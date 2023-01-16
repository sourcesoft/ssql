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
	Logger Logger
}

// Client provides common methods for using SQL.
type Client struct {
	db  *sql.DB
	tag string
	l   *log
}

// NewClient creates a new SQL client using the database connection and give options.
func NewClient(ctx context.Context, db *sql.DB, options Options) (*Client, error) {
	tag := "sql"
	if options.Tag != "" {
		tag = options.Tag
	}
	var logger Logger = defaultLogger{}
	if options.Logger != nil {
		logger = options.Logger
	}
	ret := Client{
		db:  db,
		tag: tag,
		l:   &log{logger: &logger, level: options.LogLevel},
	}
	ret.l.info("SQL client created successfully")
	return &ret, nil
}
