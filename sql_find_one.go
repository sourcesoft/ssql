package ssql

import (
	"context"
	"database/sql"
	"fmt"
)

func (c Client) SelectByID(ctx context.Context, table string, idKey, idValue string) (*sql.Rows, error) {
	c.l.debugf(fmt.Sprintf("SELECT * FROM \"%s\" WHERE %s=$1", table, idKey))
	rows, err := c.db.Query(fmt.Sprintf("SELECT * FROM \"%s\" WHERE %s=$1", table, idKey), idValue)
	return rows, err
}
