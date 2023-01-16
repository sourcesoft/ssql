package ssql

import (
	"context"
	"database/sql"
	"fmt"
)

func (c Client) DeleteByID(ctx context.Context, table string, idKey, idValue string) (*sql.Result, error) {
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
