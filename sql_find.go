package ssql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

func (c Client) FindOne(ctx context.Context, table string, idKey, idValue string) (*sql.Rows, error) {
	c.l.debugf(fmt.Sprintf("SELECT * FROM \"%s\" WHERE %s=$1", table, idKey))
	rows, err := c.db.Query(fmt.Sprintf("SELECT * FROM \"%s\" WHERE %s=$1", table, idKey), idValue)
	return rows, err
}

type FindResult struct {
	Rows        *sql.Rows
	TotalCount  *int
	optionsUsed *SQLQueryOptions
}

func (c Client) Find(ctx context.Context, options *SQLQueryOptions) (*FindResult, error) {
	if options == nil {
		options = &SQLQueryOptions{}
	}
	if options.Tag == "" {
		options.Tag = c.tag
	}
	if (options.MainSortDirection == "" && c.mainSortDirection == "") || (options.MainSortField == "" && c.mainSortField == "") {
		return nil, errors.New("'Find' method requires you to set the 'MainSortDirection' & 'MainSortField' in client config or query options")
	}
	if options.MainSortDirection == "" {
		options.MainSortDirection = c.mainSortDirection
	}
	if options.MainSortField == "" {
		options.MainSortField = c.mainSortField
	}
	options.MainSortDirection = strings.ToUpper(options.MainSortDirection)
	if options.MainSortDirection != DirectionAsc && options.MainSortDirection != DirectionDesc {
		return nil, errors.New("Sort direction can only be ASC or DESC")
	}
	// Calculate totalCount.
	var totalCount int
	if options.WithTotalCount {
		// Or use: SELECT reltuples AS estimate FROM pg_class where relname = 'user';
		// But note that the above query requires: ANALYZE VERBOSE "tablename".
		qTotalCount := fmt.Sprintf("SELECT COUNT(*) FROM \"%s\"", options.Table)
		whereStmPure, argsPure := getSQLWhereClauseFromConditions(options.Conditions, 0)
		if whereStmPure != "" {
			qTotalCount = fmt.Sprintf("%s WHERE %s", qTotalCount, whereStmPure)
		}
		err := c.db.QueryRow(qTotalCount, argsPure...).Scan(&totalCount)
		if err != nil {
			c.l.error(err, "Error preparing postgres count")
			return nil, err
		}
	}

	// Build args and prepare params.
	cursorConditions := []*ConditionPair{}
	order := ""
	pagination := ""
	if options.Params == nil {
		options.Params = &Params{}
	}
	mutateParamsByCursor(options)
	order, pagination = getSQLStmFromPaginationAndSortParams(options)
	cursorConditions = sqlParamsToConditionPairs(options)
	whereStm, args := getSQLWhereClauseFromConditions(append(options.Conditions, cursorConditions...), 0)

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
	return &FindResult{
		Rows:        rows,
		TotalCount:  &totalCount,
		optionsUsed: options,
	}, err
}

func SuperScan[T any](result *[]T, f *FindResult) (*PageInfo, error) {
	for f.Rows.Next() {
		var user T
		if err := ScanRow(&user, f.Rows); err != nil {
			return nil, err
		}
		*result = append(*result, user)
	}
	// Build the pageInfo connection.
	_, pageInfo := PrepareGraphQLConnection(result, f.optionsUsed)
	// Set StartCursor and EndCursor.
	if f.optionsUsed != nil && f.optionsUsed.MainSortField != "" && result != nil && len(*result) > 0 {
		userMappingsByTags, _ := ExtractStructMappings([]string{f.optionsUsed.Tag}, (*result)[0])
		cursorStructFieldName := userMappingsByTags[f.optionsUsed.Tag][f.optionsUsed.MainSortField]
		if pageInfo.HasPreviousPage != nil && *pageInfo.HasPreviousPage {
			field := reflect.Indirect(reflect.ValueOf((*result)[0])).FieldByName(cursorStructFieldName)
			if field.IsValid() && !field.IsNil() {
				startCursor := StrToCursor(field.Interface())
				pageInfo.StartCursor = &startCursor
			}
		}
		if pageInfo.HasNextPage != nil && *pageInfo.HasNextPage {
			field := reflect.Indirect(reflect.ValueOf((*result)[len(*result)-1])).FieldByName(cursorStructFieldName)
			if field.IsValid() && !field.IsNil() {
				endCursor := StrToCursor(field.Interface())
				pageInfo.EndCursor = &endCursor
			}
		}
	}
	if f.TotalCount != nil {
		pageInfo.TotalCount = f.TotalCount
	}
	return pageInfo, nil
}
