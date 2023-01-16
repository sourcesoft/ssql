package ssql

import (
	"context"
	"database/sql"
)

func (c Client) Raw(ctx context.Context, query string, values []interface{}) (*sql.Result, error) {
	stmt, err := c.db.Prepare(query)
	if err != nil {
		c.l.error(err, "Error preparing postgres raw query")
		return nil, err
	}
	defer stmt.Close() // Prepared statements take up server resources and should be closed after use.

	result, err := stmt.Exec(values...)
	if err != nil {
		c.l.error(err, "Error inserting to postgres")
	}

	return &result, err
}
