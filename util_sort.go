package ssql

import (
	"fmt"
)

type SortParams struct {
	Direction *string `json:"direction" query:"direction"` // ASC/DESC (sortby).
	Field     *string `json:"field" query:"field"`         // Field name (order by).
}

func (param SortParams) getSQLStatement() string {
	stm := ""
	if param.Field != nil {
		stm = *param.Field
	}
	if param.Direction != nil {
		stm = fmt.Sprintf("%s %s", stm, *param.Direction)
	}
	return stm
}
