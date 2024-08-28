package sqlwrap

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type BeforeFunc func(ctx context.Context, query string, args ...interface{}) context.Context
type AfterFunc func(ctx context.Context, err error, query string, args ...interface{})

// Queryer interface
type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

// Database interface
type Database interface {
	sq.StdSqlCtx
	BeginTx(context.Context, *sql.TxOptions) (Transaction, error)
	Ping() error
	PingContext(ctx context.Context) error
	Close() error
}

// Transaction interface
type Transaction interface {
	sq.StdSqlCtx
	Rollback() error
	Commit() error
}
