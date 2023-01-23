package ssql

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

// GetSQLFieldValuePairs uses cursor params to get 'conditions' using the following rules:
// 1. If the after argument is provided => add cursor-field > after to the WHERE clause.
// 2. If the before argument is provided => add cursor-field < before to the WHERE clause.
//
// Note:
// - For next steps 3 and 4 check function 'MutateParamsByCursor'.
// - For next steps 5-10 check function 'PrepareGraphQLConnection'.
func (param CursorParams) getSQLFieldValuePairs(arg *SQLQueryOptions) *ConditionPair {
	if (param.After == nil && param.Before == nil) || arg.MainSortField == "" {
		return nil
	}
	comparisonOp := OPGreater
	var iterator int = 0
	if param.After != nil && CursorToInt(*param.After) > 0 {
		comparisonOp = OPGreater
		iterator = CursorToInt(*param.After)
	} else if param.Before != nil && CursorToInt(*param.Before) > 0 {
		comparisonOp = OPLess
		iterator = CursorToInt(*param.Before)
	}

	if arg.MainSortDirection == DirectionDesc {
		if comparisonOp == OPGreater {
			comparisonOp = OPLess
		} else {
			comparisonOp = OPGreater
		}
	}

	fvp := &ConditionPair{
		Field: arg.MainSortField,
		Op:    comparisonOp,
		Value: iterator,
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

func StrToCursor(cursor interface{}) string {
	v := reflect.ValueOf(cursor)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int32, reflect.Int64:
		return base64Encode(fmt.Sprintf("c:%d", v.Int()))
	case reflect.Uint, reflect.Uint8, reflect.Uint32, reflect.Uint64:
		return base64Encode(fmt.Sprintf("c:%d", v.Uint()))
	case reflect.String:
		return base64Encode(fmt.Sprintf("c:%s", v.String()))
	case reflect.Float32, reflect.Float64:
		return base64Encode(fmt.Sprintf("c:%b", v.Float()))
	}
	return base64Encode(fmt.Sprintf("c:%+v", cursor))
}
