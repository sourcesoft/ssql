package ssql

import (
	"fmt"
)

func (param OffsetParams) GetSQLStatement() string {
	stm := ""
	if param.Limit != nil {
		stm = fmt.Sprintf("LIMIT %d", *param.Limit)
	}
	if param.Offset != nil {
		stm = fmt.Sprintf("%s OFFSET %d", stm, *param.Offset)
	}
	return stm
}

type OffsetParams struct {
	Limit  *int `json:"limit,omitempty"`
	Offset *int `json:"offset,omitempty"`
	Used   bool `json:"-"`
}
