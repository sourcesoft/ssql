package ssql

import (
	"context"
	"database/sql"

	"github.com/sourcesoft/ssql"
)

// MClient provides common methods interface for using SQL.
type MClient interface {
	Raw(ctx context.Context, query string, values []interface{}) (*sql.Result, error)
	Insert(ctx context.Context, table string, payload interface{}) (*sql.Result, error)
	SelectByID(ctx context.Context, table string, idKey, idValue string) (*sql.Rows, error)
	Find(ctx context.Context, options *ssql.SQLQueryOptions) (*sql.Rows, *int, error)
	DeleteByID(ctx context.Context, table string, idKey, idValue string) (*sql.Result, error)
	UpdateByID(ctx context.Context, table string, idKey, idValue string, payload interface{}) (*sql.Result, error)
}

type MockClient struct {
	OnRaw        func(query string, values []interface{}) (*sql.Result, error)
	OnInsert     func(table string, payload interface{}) (*sql.Result, error)
	OnSelectByID func(table string, idKey, idValue string) (*sql.Rows, error)
	OnFind       func(options *ssql.SQLQueryOptions) (*sql.Rows, *int, error)
	OnDeleteByID func(table string, idKey, idValue string) (*sql.Result, error)
	OnUpdateByID func(table string, idKey, idValue string, payload interface{}) (*sql.Result, error)
}

func NewMockClient(mock MockClient) MClient {
	return mock
}

func (m MockClient) Raw(_ context.Context, query string, values []interface{}) (*sql.Result, error) {
	return m.OnRaw(query, values)
}

func (m MockClient) Insert(_ context.Context, table string, payload interface{}) (*sql.Result, error) {
	return m.OnInsert(table, payload)
}

func (m MockClient) SelectByID(_ context.Context, table string, idKey, idValue string) (*sql.Rows, error) {
	return m.OnSelectByID(table, idKey, idValue)
}

func (m MockClient) Find(_ context.Context, options *ssql.SQLQueryOptions) (*sql.Rows, *int, error) {
	return m.OnFind(options)
}

func (m MockClient) DeleteByID(_ context.Context, table string, idKey, idValue string) (*sql.Result, error) {
	return m.OnDeleteByID(table, idKey, idValue)
}

func (m MockClient) UpdateByID(_ context.Context, table string, idKey, idValue string, payload interface{}) (*sql.Result, error) {
	return m.OnUpdateByID(table, idKey, idValue, payload)
}
