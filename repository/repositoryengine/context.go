package repositoryengine

import (
	"context"

	"github.com/jmoiron/sqlx"
)

type contextTxKey string

const (
	txKey contextTxKey = "go-rest-api.tx"
)

// contextWithTransaction creates context with transaction
func contextWithTransaction(ctx context.Context, tx *sqlx.Tx) context.Context {
	return context.WithValue(ctx, txKey, tx)
}

// transactionFromContext retrieve transaction from context
func transactionFromContext(ctx context.Context) *sqlx.Tx {
	if value := ctx.Value(txKey); value != nil {
		if tx, ok := value.(*sqlx.Tx); ok {
			return tx
		}
	}
	return nil
}
