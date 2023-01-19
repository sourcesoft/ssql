package ssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (c Client) UpdateOne(ctx context.Context, table string, idKey, idValue string, payload interface{}) (*sql.Result, error) {
	// First get struct pairs by tag.
	pairs := getPairsByTag(c.tag, payload)
	c.l.debugf("UPDATE: DB pairs %+v", pairs)

	// Prepare statement.
	for index, value := range pairs.Fields {
		pairs.Fields[index] = fmt.Sprintf("%s = $%d", value, index+1)
	}
	fieldsStr := strings.Join(pairs.Fields, ", ")
	c.l.debugf("UPDATE: DB fieldsStr %+v", fieldsStr)
	stmt, err := c.db.Prepare(fmt.Sprintf("UPDATE \"%s\" SET %s WHERE %s=$%d", table, fieldsStr, idKey, len(pairs.Fields)+1))
	c.l.debug("UPDATE \"%s\" SET %s WHERE %s=$%d", table, fieldsStr, idKey, len(pairs.Fields)+1)
	if err != nil {
		c.l.error(err, "Error preparing postgres update")
		return nil, err
	}
	defer stmt.Close()

	// Execute with values.
	pairs.Values = append(pairs.Values, idValue) // No need because id is $1 again.
	c.l.debugf("UPDATE: DB pairs Values %+v", pairs.Values)
	result, err := stmt.Exec(pairs.Values...)
	if err != nil {
		c.l.error(err, "Error executing query in postgres")
	}

	return &result, err
}

func (c Client) Update(ctx context.Context, table string, conditions []*ConditionPair, payload interface{}) (*sql.Result, error) {
	// First get struct pairs by tag.
	pairs := getPairsByTag(c.tag, payload)
	c.l.debugf("UPDATE: DB pairs %+v", pairs)

	// Prepare statement.
	for index, value := range pairs.Fields {
		pairs.Fields[index] = fmt.Sprintf("%s = $%d", value, index+1)
	}
	whereStm, args := getSQLWhereClauseFromConditions(conditions, len(pairs.Fields))
	fieldsStr := strings.Join(pairs.Fields, ", ")
	c.l.debugf("UPDATE: DB fieldsStr %+v", fieldsStr)
	stmt, err := c.db.Prepare(fmt.Sprintf("UPDATE \"%s\" SET %s WHERE %s", table, fieldsStr, whereStm))
	c.l.debugf("UPDATE \"%s\" SET %s WHERE %s", table, fieldsStr, whereStm)
	if err != nil {
		c.l.error(err, "Error preparing postgres update")
		return nil, err
	}
	defer stmt.Close()

	// Execute with values.
	pairs.Values = append(pairs.Values, args...)
	c.l.debugf("UPDATE: DB pairs Values %+v", pairs.Values)
	result, err := stmt.Exec(pairs.Values...)
	if err != nil {
		c.l.error(err, "Error executing query in postgres")
	}

	return &result, err
}
