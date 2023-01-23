package ssql

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
)

type Options struct {
	// Struct tag to look for "sql" field names in.
	Tag string
	// If not set (0), logging will be disabled
	LogLevel int
	// If not set, the default logger will be used
	Logger *Logger
	// Main default pagination sort config. It can be overriden
	// by Find query options.
	MainSortField     string
	MainSortDirection string
}

// Client provides common methods for using SQL.
type Client struct {
	db                *sql.DB
	tag               string
	mainSortField     string
	mainSortDirection string
	l                 *log
}

// NewClient creates a new SQL client using the database connection and give options.
func NewClient(ctx context.Context, db *sql.DB, options *Options) (*Client, error) {
	tag := "sql"
	var logger Logger = &defaultLogger{}
	logLevel := 0
	if options != nil {
		if options.Tag != "" {
			tag = options.Tag
		}
		if options.Logger != nil {
			logger = *options.Logger
		}
		if options.LogLevel != 0 {
			logLevel = options.LogLevel
		}
	}

	ret := Client{
		db:                db,
		tag:               tag,
		mainSortField:     options.MainSortField,
		mainSortDirection: options.MainSortDirection,
		l: &log{
			logger: &logger,
			level:  logLevel,
		},
	}
	ret.l.info("SQL client created successfully")
	return &ret, nil
}

const (
	// Used by query options and the API.
	DirectionDesc    = "DESC"
	DirectionAsc     = "ASC"
	OPEqual          = ">"
	OPNotEqual       = "<>"
	OPGreater        = ">"
	OPLess           = "<"
	OPGreatorOrEqual = ">="
	OPLessOrEqual    = "<="
	OPLogicalIn      = "IN"
	// Only used internally.
	opLogicalAnd = "AND"
	// OPLogicalOr      = "OR"
	// OPLogicalNot     = "NOT"
	// OPLogicalLike    = "LIKE"
	// OPLogicalExists  = "EXISTS"
	// OPLogicalBetween = "BETWEEN"
)
