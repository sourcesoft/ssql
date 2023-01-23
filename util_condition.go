package ssql

import (
	"fmt"
	"reflect"
	"strings"
)

func sqlParamsToConditionPairs(arg *SQLQueryOptions) []*ConditionPair {
	cursorPairs := []*ConditionPair{}
	// Cursors also may enforce conditions.
	if arg.Params.CursorParams != nil {
		pair := arg.Params.CursorParams.getSQLFieldValuePairs(arg)
		if pair != nil {
			cursorPairs = append(cursorPairs, pair)
		}
	}
	return cursorPairs
}

func getSQLWhereClauseFromConditions(conditions []*ConditionPair, startIndex int) (whereStm string, args []interface{}) {
	if len(conditions) > 0 {
		// Prepare statement.
		pairs := []string{}
		for index, value := range conditions {
			if strings.ToUpper(value.Op) == OPLogicalIn && reflect.TypeOf(value.Value).Kind() == reflect.Slice {
				// Trying to print "fieldName IN ($1, $2, $3)"
				vals := []interface{}{}
				rv := reflect.ValueOf(value.Value)
				iter := ""
				for i := 0; i < rv.Len(); i++ {
					vals = append(vals, rv.Index(i).Interface())
					iter = fmt.Sprintf("%s, $%d", iter, startIndex+i+1)
				}
				iter = strings.TrimPrefix(iter, ", ")
				pairs = append(pairs, fmt.Sprintf("%s %s (%s)", value.Field, value.Op, iter))
				args = append(args, vals...)
			} else {
				pairs = append(pairs, fmt.Sprintf("%s %s $%d", value.Field, value.Op, startIndex+index+1))
				args = append(args, value.Value)
			}
		}
		whereStm = strings.Join(pairs, fmt.Sprintf(" %s ", opLogicalAnd))
	}
	return whereStm, args
}
