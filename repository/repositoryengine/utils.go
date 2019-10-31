package repositoryengine

import (
	"context"

	"github.com/jmoiron/sqlx"
)

// sqlxExt retrieve current active sqlx.Ext instance. ext points to default value
func sqlxExt(ctx context.Context, ext sqlx.Ext) sqlx.Ext {
	if tx := transactionFromContext(ctx); tx != nil {
		return tx
	}
	return ext
}
