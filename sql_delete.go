package ssql

import (
	"context"
	"database/sql"
	"fmt"
)

func (c Client) DeleteOne(ctx context.Context, table string, idKey, idValue string) (*sql.Result, error) {
	stmt, err := c.db.Prepare(fmt.Sprintf("DELETE FROM \"%s\" WHERE %s=$1", table, idKey))
	if err != nil {
		c.l.error(err, "Error preparing postgres deleting")
		return nil, err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	result, err := stmt.Exec(idValue)
	if err != nil {
		c.l.error(err, "Error deleting to postgres")
	}

	return &result, err
}

func (c Client) Delete(ctx context.Context, table string, conditions []*ConditionPair) (*sql.Result, error) {
	whereStm, args := getSQLWhereClauseFromConditions(conditions, 0)
	stmt, err := c.db.Prepare(fmt.Sprintf("DELETE FROM \"%s\" WHERE %s", table, whereStm))
	if err != nil {
		c.l.error(err, "Error preparing postgres deleting")
		return nil, err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	result, err := stmt.Exec(args...)
	if err != nil {
		c.l.error(err, "Error deleting to postgres")
	}

	return &result, err
}
