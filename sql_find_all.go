package ssql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

func (c Client) Find(ctx context.Context, options *SQLQueryOptions) (*sql.Rows, *int, error) {
	// Calculate totalCount.
	var totalCount int
	if options.WithTotalCount {
		// Or use: SELECT reltuples AS estimate FROM pg_class where relname = 'user';
		// But note that the above query requires: ANALYZE VERBOSE "tablename".
		qTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", options.Table)
		whereStmPure, argsPure := getSQLWhereClauseFromConditions(options.Conditions)
		if whereStmPure != "" {
			qTotalCount = fmt.Sprintf("%s WHERE %s", qTotalCount, whereStmPure)
		}
		c.l.debugf("count statement %s with args %+v", qTotalCount, argsPure)
		err := c.db.QueryRow(qTotalCount, argsPure...).Scan(&totalCount)
		if err != nil {
			c.l.error(err, "Error preparing postgres count")
			return nil, &totalCount, err
		}
	}

	// Build args and prepare params.
	cursorConditions := []*ConditionPair{}
	order := ""
	pagination := ""
	if options.Params != nil {
		mutateParamsByCursor(options.Params)
		order, pagination = getSQLStmFromPaginationAndSortParams(options.Params)
		cursorConditions = sqlParamsToConditionPairs(options.Params)
	}
	whereStm, args := getSQLWhereClauseFromConditions(append(options.Conditions, cursorConditions...))

	// Build query statement.
	fieldsSeparated := "*"
	if len(options.Fields) > 0 {
		fields := []string{}
		for key := range options.Fields {
			fields = append(fields, key)
		}
		fieldsSeparated = strings.Join(fields, ", ")
	}
	q := fmt.Sprintf("SELECT %s FROM \"%s\"", fieldsSeparated, options.Table)
	if whereStm != "" {
		q = fmt.Sprintf("%s WHERE %s", q, whereStm)
	}
	if order != "" {
		q = fmt.Sprintf("%s %s", q, order)
	}
	if pagination != "" {
		q = fmt.Sprintf("%s %s", q, pagination)
	}
	c.l.debugf("find statement %s with args %+v", q, args)

	// Run the query.
	rows, err := c.db.Query(q, args...)
	return rows, &totalCount, err
}
