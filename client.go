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
}

// Client provides common methods for using SQL.
type Client struct {
	db  *sql.DB
	tag string
	l   *log
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
		db:  db,
		tag: tag,
		l:   &log{logger: &logger, level: logLevel},
	}
	ret.l.info("SQL client created successfully")
	return &ret, nil
}
