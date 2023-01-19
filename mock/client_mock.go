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
	FindOne(ctx context.Context, table string, idKey, idValue string) (*sql.Rows, error)
	Find(ctx context.Context, options *ssql.SQLQueryOptions) (*sql.Rows, *int, error)
	DeleteOne(ctx context.Context, table string, idKey, idValue string) (*sql.Result, error)
	Delete(ctx context.Context, table string, conds []*ssql.ConditionPair) (*sql.Result, error)
	UpdateOne(ctx context.Context, table string, idKey, idValue string, payload interface{}) (*sql.Result, error)
	Update(ctx context.Context, table string, conds []*ssql.ConditionPair, payload interface{}) (*sql.Result, error)
}

type MockClient struct {
	OnRaw       func(query string, values []interface{}) (*sql.Result, error)
	OnInsert    func(table string, payload interface{}) (*sql.Result, error)
	OnFindOne   func(table string, idKey, idValue string) (*sql.Rows, error)
	OnFind      func(options *ssql.SQLQueryOptions) (*sql.Rows, *int, error)
	OnDeleteOne func(table string, idKey, idValue string) (*sql.Result, error)
	OnDelete    func(table string, conds []*ssql.ConditionPair) (*sql.Result, error)
	OnUpdateOne func(table string, idKey, idValue string, payload interface{}) (*sql.Result, error)
	OnUpdate    func(table string, conds []*ssql.ConditionPair, payload interface{}) (*sql.Result, error)
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

func (m MockClient) FindOne(_ context.Context, table string, idKey, idValue string) (*sql.Rows, error) {
	return m.OnFindOne(table, idKey, idValue)
}

func (m MockClient) Find(_ context.Context, options *ssql.SQLQueryOptions) (*sql.Rows, *int, error) {
	return m.OnFind(options)
}

func (m MockClient) DeleteOne(_ context.Context, table, idKey, idValue string) (*sql.Result, error) {
	return m.OnDeleteOne(table, idKey, idValue)
}

func (m MockClient) Delete(_ context.Context, table string, conds []*ssql.ConditionPair) (*sql.Result, error) {
	return m.OnDelete(table, conds)
}

func (m MockClient) UpdateOne(_ context.Context, table string, idKey, idValue string, payload interface{}) (*sql.Result, error) {
	return m.OnUpdateOne(table, idKey, idValue, payload)
}

func (m MockClient) Update(_ context.Context, table string, conds []*ssql.ConditionPair, payload interface{}) (*sql.Result, error) {
	return m.OnUpdate(table, conds, payload)
}
