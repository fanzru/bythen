package sqlwrap

import (
	"context"

	sq "github.com/Masterminds/squirrel"
)

type ctxKeyType int

const txKey ctxKeyType = 0

// TxFromContext returns the transaction object from the context. Return db if not exists
func TxFromContext(ctx context.Context, db sq.StdSqlCtx) Queryer {
	tx, ok := ctx.Value(txKey).(Transaction)
	if !ok {
		return db
	}
	return tx
}

// TransactionFromContext returns the transaction object from the context. Return db if not exists
func TransactionFromContext(ctx context.Context) Transaction {
	tx, ok := ctx.Value(txKey).(Transaction)
	if !ok {
		return nil
	}
	return tx
}

// ContextWithTx add database transaction to context
func ContextWithTx(parentContext context.Context, tx Transaction) context.Context {
	return context.WithValue(parentContext, txKey, tx)
}
