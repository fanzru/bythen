package sqlwrap

import (
	"context"
	"database/sql"

	sq "github.com/Masterminds/squirrel"
)

type DB struct {
	DB     *sql.DB
	before []BeforeFunc
	after  []AfterFunc
}

func NewDB(db *sql.DB) *DB {
	return &DB{
		DB:     db,
		before: make([]BeforeFunc, 0),
		after:  make([]AfterFunc, 0),
	}
}

var _ Database = (*DB)(nil)

func (db *DB) AddBeforeFunc(f BeforeFunc) {
	db.before = append(db.before, f)
}

func (db *DB) AddAfterFunc(f AfterFunc) {
	db.after = append(db.after, f)
}

var _ sq.StdSqlCtx = (*DB)(nil)

func (db DB) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	ctx := db.doBefore(context.Background(), query, args...)
	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.Exec(query, args...)
	return
}

func (db DB) Query(query string, args ...interface{}) (res *sql.Rows, err error) {
	ctx := db.doBefore(context.Background(), query, args...)
	defer db.doAfter(ctx, err, query, args...)

	res, err = db.DB.Query(query, args...)
	return
}

func (db DB) QueryRow(query string, args ...interface{}) *sql.Row {
	ctx := db.doBefore(context.Background(), query, args...)

	res := db.DB.QueryRow(query, args...)
	err := res.Err()

	db.doAfter(ctx, err, query, args...)
	return res
}

func (db DB) QueryContext(ctx context.Context, query string, args ...interface{}) (res *sql.Rows, err error) {
	ctx = db.doBefore(ctx, query, args...)
	defer db.doAfter(ctx, err, query, args...)

	res, err = TxFromContext(ctx, db.DB).QueryContext(ctx, query, args...)
	return
}

func (db DB) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	ctx = db.doBefore(ctx, query, args...)

	res := TxFromContext(ctx, db.DB).QueryRowContext(ctx, query, args...)
	err := res.Err()

	db.doAfter(ctx, err, query, args...)
	return res
}

func (db DB) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	ctx = db.doBefore(ctx, query, args...)
	defer db.doAfter(ctx, err, query, args...)

	res, err = TxFromContext(ctx, db.DB).ExecContext(ctx, query, args...)
	return
}

func (db *DB) BeginTx(ctx context.Context, options *sql.TxOptions) (Transaction, error) {
	tx, err := db.DB.BeginTx(ctx, options)
	if err != nil {
		return nil, err
	}

	return &Tx{
		tx:     tx,
		before: db.before,
		after:  db.after,
	}, nil
}

func (db *DB) Ping() (err error) {
	return db.DB.Ping()
}

func (db *DB) PingContext(ctx context.Context) error {
	return db.DB.PingContext(ctx)
}

func (db *DB) Close() error {
	return db.DB.Close()
}

func (db DB) doBefore(ctx context.Context, query string, args ...interface{}) context.Context {
	for _, f := range db.before {
		ctx = f(ctx, query, args...)
	}
	return ctx
}

func (db DB) doAfter(ctx context.Context, err error, query string, args ...interface{}) {
	for _, f := range db.after {
		f(ctx, err, query, args...)
	}
}

type Tx struct {
	tx     *sql.Tx
	before []BeforeFunc
	after  []AfterFunc
}

var _ Transaction = (*Tx)(nil)

func (t Tx) Exec(query string, args ...interface{}) (res sql.Result, err error) {
	t.doBefore(context.Background(), query, args...)
	defer t.doAfter(context.Background(), err, query, args...)

	res, err = t.tx.Exec(query, args...)
	return
}

func (t Tx) Query(query string, args ...interface{}) (res *sql.Rows, err error) {
	t.doBefore(context.Background(), query, args...)
	defer t.doAfter(context.Background(), err, query, args...)

	res, err = t.tx.Query(query, args...)
	return
}

func (t Tx) QueryRow(query string, args ...interface{}) *sql.Row {

	t.doBefore(context.Background(), query, args...)

	res := t.tx.QueryRow(query, args...)
	err := res.Err()

	t.doAfter(context.Background(), err, query, args...)
	return res
}

func (t Tx) QueryContext(ctx context.Context, query string, args ...interface{}) (res *sql.Rows, err error) {
	res, err = t.tx.QueryContext(ctx, query, args...)
	return
}

func (t Tx) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	res := t.tx.QueryRowContext(ctx, query, args...)
	return res
}

func (t Tx) ExecContext(ctx context.Context, query string, args ...interface{}) (res sql.Result, err error) {
	res, err = t.tx.ExecContext(ctx, query, args...)
	return
}

func (t Tx) Rollback() error {
	return t.tx.Rollback()
}

func (t Tx) Commit() error {
	return t.tx.Commit()
}

func (t Tx) doBefore(ctx context.Context, query string, args ...interface{}) {
	for _, f := range t.before {
		f(ctx, query, args...)
	}
}

func (t Tx) doAfter(ctx context.Context, err error, query string, args ...interface{}) {
	for _, f := range t.after {
		f(ctx, err, query, args...)
	}
}
