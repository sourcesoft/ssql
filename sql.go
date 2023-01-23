package ssql

import (
	"database/sql"
	"fmt"
	"reflect"

	"github.com/georgysavva/scany/v2/sqlscan"
)

// GetPairsByTag returns you an array of FieldsPair by tag name.
func getPairsByTag(tagName string, s interface{}) FieldValuePairs {
	var pairs FieldValuePairs
	t := reflect.TypeOf(s)
	v := reflect.ValueOf(s)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		col := field.Tag.Get(tagName)
		val := v.Field(i).Interface()
		if col == "" || val == nil || v.Field(i).IsNil() {
			continue
		}
		pairs.Fields = append(pairs.Fields, col)
		pairs.Values = append(pairs.Values, val)
	}
	return pairs
}

// MutateParamsByCursor mutates the 'pagination' and 'sort' of passed params using crsors.
// 3. If the first argument is provided => add ORDER BY id DESC LIMIT first+1 to the query.
// 4. If the last argument is provided => add ORDER BY id ASC LIMIT last+1 to the query.
//
// Note:
// - For previous steps check function 'GetSQLFieldValuePairs'.
// - For next steps 5-10 check function 'PrepareGraphQLConnection'.
func mutateParamsByCursor(arg *SQLQueryOptions) *Params {
	if arg == nil || arg.Params == nil {
		return nil
	}
	limit := 10 // Default
	dir := arg.MainSortDirection
	if arg.Params.CursorParams != nil && arg.Params.CursorParams.First != nil {
		limit = *arg.Params.CursorParams.First + 1
	} else if arg.Params.CursorParams != nil && arg.Params.CursorParams.Last != nil {
		if dir == DirectionAsc {
			dir = DirectionDesc
		} else {
			dir = DirectionAsc
		}
		limit = *arg.Params.CursorParams.Last + 1
	}

	if (arg.Params.OffsetParams == nil || arg.Params.OffsetParams.Limit == nil) && arg.Params.CursorParams != nil && arg.MainSortField != "" {
		arg.Params.OffsetParams = &OffsetParams{
			Limit: &limit,
		}
	} else if arg.Params.OffsetParams != nil {
		if arg.Params.OffsetParams.Limit == nil {
			arg.Params.OffsetParams.Limit = &limit
		} else {
			limit = *arg.Params.OffsetParams.Limit
			limit++
			arg.Params.OffsetParams.Limit = &limit
		}
	}

	field := arg.MainSortField
	arg.Params.SortParams = append(arg.Params.SortParams, &SortParams{
		Direction: &dir,
		Field:     &field,
	})
	return arg.Params
}

func getSQLStmFromPaginationAndSortParams(arg *SQLQueryOptions) (paramsOrder, paramsPagination string) {
	if arg == nil || arg.Params == nil {
		return "", ""
	}
	// ORDER BY.
	for sIndex, sParam := range arg.Params.SortParams {
		split := ""
		if sIndex > 0 && sIndex < len(arg.Params.SortParams) {
			split = ", "
		}
		if sParam != nil && sParam.getSQLStatement() != "" {
			paramsOrder = fmt.Sprintf("%s %s%s", paramsOrder, split, sParam.getSQLStatement())
		}
	}
	if paramsOrder != "" {
		paramsOrder = fmt.Sprintf("ORDER BY %s", paramsOrder)
	}
	// LIMIT/OFFSET Pagination.
	if arg.Params.OffsetParams != nil && arg.Params.OffsetParams.GetSQLStatement() != "" {
		if paramsPagination != "" {
			paramsPagination = fmt.Sprintf("%s ", paramsPagination)
		}
		paramsPagination = fmt.Sprintf("%s%s", paramsPagination, arg.Params.OffsetParams.GetSQLStatement())
	}
	return paramsOrder, paramsPagination
}

func ScanRow(dst interface{}, rows *sql.Rows) error {
	if err := sqlscan.ScanRow(dst, rows); err != nil {
		return err
	}
	return nil
}

func ScanOne(dst interface{}, rows *sql.Rows) error {
	if err := sqlscan.ScanOne(dst, rows); err != nil {
		return err
	}
	return nil
}
