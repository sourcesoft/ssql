package ssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (c Client) Insert(ctx context.Context, table string, payload interface{}) (*sql.Result, error) {
	// First get struct pairs by tag.
	pairs := getPairsByTag(c.tag, payload)
	c.l.debugf("DB pairs %+v", pairs)

	// Prepare statement.
	fieldsStr := strings.Join(pairs.Fields, ", ")
	indexStr := ""
	valuesArray := []string{}
	for index := range pairs.Fields {
		valuesArray = append(valuesArray, fmt.Sprintf("$%d", index+1))
	}
	indexStr = strings.Join(valuesArray, ", ")
	stmt, err := c.db.Prepare(fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s)", table, fieldsStr, indexStr))
	c.l.debugf(fmt.Sprintf("INSERT INTO \"%s\" (%s) VALUES (%s)", table, fieldsStr, indexStr))
	if err != nil {
		c.l.error(err, "Error preparing postgres insert")
		return nil, err
	}
	defer stmt.Close()

	// Execute with values.
	c.l.debug(pairs.Values...)
	result, err := stmt.Exec(pairs.Values...)
	if err != nil {
		c.l.errorf(err, "Error executing query in postgres")
	}

	return &result, err
}
