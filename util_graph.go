package ssql

// PrepareGraphQLConnection puts the final touches on the result.
// 5. If the last argument is provided => reverse the order of the results.
// 6. If no 'less than first+1' results are returned => set hasNextPage: true.
// 7. If the 'last' argument is provided =>  set hasNextPage: false.
// 8. If no 'less last+1' results are returned => set hasPreviousPage: true
// 9. If 'after' is set AND server can efficiently determine that elements exist prior to after => set hasNextPage: false.
// 10. Finally remove the last result cause we get one extra at the end.
//
// Note:
// - For previous steps 1 and 2 check function 'GetSQLFieldValuePairs'.
// - For previous steps 3 and 4 check function 'MutateParamsByCursor'.
func PrepareGraphQLConnection[T any](result []T, params *Params) (*[]T, *PageInfo) {
	falseVal := false
	trueVal := true
	pageInfo := PageInfo{
		HasPreviousPage: &falseVal,
		HasNextPage:     &falseVal,
		StartCursor:     nil,
		EndCursor:       nil,
	}
	// 5. If the last argument is provided => reverse the order of the results.
	r := reverse(&result)
	if params != nil {
		// 6. If no 'less than first+1' results are returned => set hasNextPage: true.
		firstValue := 10 // Default limit.
		if params.CursorParams.First != nil {
			firstValue = *params.CursorParams.First
		}
		if firstValue < len(*r) {
			pageInfo.HasNextPage = &trueVal
		}
		// 7. If the 'last' argument is provided => set hasNextPage: false.
		if params.CursorParams.Last != nil {
			pageInfo.HasNextPage = &falseVal
		}
		// 8. If no 'less last+1' results are returned => set hasPreviousPage: true
		if params.CursorParams.Last != nil && *params.CursorParams.Last < len(*r) {
			pageInfo.HasPreviousPage = &trueVal
		}
		// 9. If 'after' is set AND server can efficiently determine that elements exist prior to after => set hasNextPage: false.
		if params.CursorParams.After != nil {
			pageInfo.HasPreviousPage = &trueVal
		}
	}
	// 10. Finally remove the last result cause we get one extra at the end.
	limit := 10
	if params.OffsetParams.Used && params.OffsetParams.Limit != nil {
		limit = *params.OffsetParams.Limit
	}
	if params.CursorParams.Used && params.CursorParams.First != nil {
		limit = *params.CursorParams.First
	}
	if len(*r) > limit {
		r = popArray(r)
	}
	return r, &pageInfo
}
