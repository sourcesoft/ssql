package ssql

import (
	"fmt"
	"strconv"
	"strings"
)

// GetSQLFieldValuePairs uses cursor params to get 'conditions' using the following rules:
// 1. If the after argument is provided => add id > parsed_cursor to the WHERE clause.
// 2. If the before argument is provided => add id < parsed_cursor to the WHERE clause.
//
// Note:
// - For next steps 3 and 4 check function 'MutateParamsByCursor'.
// - For next steps 5-10 check function 'PrepareGraphQLConnection'.
func (param CursorParams) getSQLFieldValuePairs() *ConditionPair {
	if param.After == nil && param.Before == nil {
		return nil
	}
	comparisonOp := ">"
	var iterator int = 0
	if param.After != nil && CursorToInt(*param.After) > 0 {
		comparisonOp = ">"
		iterator = CursorToInt(*param.After)
	} else if param.Before != nil && CursorToInt(*param.Before) > 0 {
		comparisonOp = "<"
		iterator = CursorToInt(*param.Before)
	}

	fvp := &ConditionPair{
		Field: cursorField,
		Value: iterator,
		Op:    comparisonOp,
	}
	return fvp
}

// CursorParams is the cursor pagination used by GraphQL usually.
// 'after': Returns the elements in the list that come after the specified cursor.
// 'before': Returns the elements in the list that come before the specified cursor.
// 'first': Returns the first n elements from the list.
// 'last': Returns the last n elements from the list.
type CursorParams struct {
	After  *string `json:"after" query:"after"`
	Before *string `json:"before" query:"before"`
	First  *int    `json:"first" query:"first"`
	Last   *int    `json:"last" query:"last"`
	Used   bool    `json:"-" query:"-"`
}

const cursorField = "created_at"

func IntToCursor(cursor int) string {
	return base64Encode(fmt.Sprintf("cursor:%d", cursor))
}

func CursorToInt(cursor string) int {
	origin, err := base64Decode(cursor)
	if err != nil {
		return 0
	}
	parts := strings.Split(origin, ":")
	if len(parts) < 2 {
		return 0
	}
	c, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}
	return c
}

func StrToCursor(cursor string) string {
	return base64Encode(fmt.Sprintf("cursor:%s", cursor))
}
