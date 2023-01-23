package ssql

type ConditionPair struct {
	Field string
	Value interface{}
	Op    string
}

type SQLQueryOptions struct {
	Table             string
	Tag               string // Override client tag config explicitly per query.
	MainSortField     string
	MainSortDirection string
	Fields            map[string]bool // omit to query all fields.
	WithTotalCount    bool
	Params            *Params
	Conditions        []*ConditionPair
}

type Params struct {
	OffsetParams *OffsetParams
	SortParams   []*SortParams
	CursorParams *CursorParams
}

// FieldValuePairs describes a piece of data that is stored in a database table column.
type FieldValuePairs struct {
	Fields []string
	Values []interface{}
}

// Used for in response for Find method, returning pagination info for both offset and cursor pagination.
type PageInfo struct {
	HasPreviousPage *bool   `json:"hasPreviousPage,omitempty"`
	HasNextPage     *bool   `json:"hasNextPage,omitempty"`
	StartCursor     *string `json:"startCursor,omitempty"`
	EndCursor       *string `json:"endCursor,omitempty"`
	TotalCount      *int    `json:"-"`
}
