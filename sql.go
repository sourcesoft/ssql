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
func mutateParamsByCursor(arg *Params) *Params {
	if arg == nil {
		return nil
	}
	if arg.CursorParams != nil && arg.CursorParams.Used {
		field := cursorField
		// Note: offset/limit pagination is out of question when cursor is used.
		// For cursor pagination sometimes we need limit in addition to ordering.
		dir := "DESC"
		limit := 0
		if arg.CursorParams.First != nil && *arg.CursorParams.First != 0 {
			limit = (*arg.CursorParams.First) + 1
		} else if arg.CursorParams.Last != nil && *arg.CursorParams.Last != 0 {
			dir = "ASC"
			limit = (*arg.CursorParams.Last) + 1
		}
		arg.SortParams = append(arg.SortParams, &SortParams{
			Direction: &dir,
			Field:     &field,
		})
		arg.OffsetParams = &OffsetParams{
			Limit: &limit,
		}
	}
	return arg
}

func getSQLStmFromPaginationAndSortParams(arg *Params) (paramsOrder, paramsPagination string) {
	if arg == nil {
		return "", ""
	}
	// ORDER BY.
	for sIndex, sParam := range arg.SortParams {
		split := ""
		if sIndex > 0 && sIndex < len(arg.SortParams) {
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
	if arg.OffsetParams != nil && arg.OffsetParams.GetSQLStatement() != "" {
		if paramsPagination != "" {
			paramsPagination = fmt.Sprintf("%s ", paramsPagination)
		}
		paramsPagination = fmt.Sprintf("%s%s", paramsPagination, arg.OffsetParams.GetSQLStatement())
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
