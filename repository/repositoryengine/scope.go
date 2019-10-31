package repositoryengine

import (
	"context"
	"database/sql"

	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/errutil"
	"github.com/jmoiron/sqlx"
)

// scope implements repository.Scope interface
type scope struct {
	db *sqlx.DB
}

// newScope creates new instance of repository.Scope interface
func newScope(db *sqlx.DB) repository.Scope {
	return &scope{
		db: db,
	}
}

func (s *scope) WithContext(ctx context.Context) (context.Context, error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, errutil.Wrap(err, "failed to start transaction")
	}
	return contextWithTransaction(ctx, tx), nil
}

func (s *scope) Complete(ctx context.Context) error {
	tx := transactionFromContext(ctx)
	if tx == nil {
		return nil
	}

	err := tx.Commit()
	if err == sql.ErrTxDone {
		return nil // Ignore this error, because it's safe
	}
	return errutil.Wrap(err, "failed to commit transaction")
}

func (s *scope) Cancel(ctx context.Context) error {
	tx := transactionFromContext(ctx)
	if tx == nil {
		return nil
	}

	err := tx.Rollback()
	if err == sql.ErrTxDone {
		return nil // Ignore this error, because it's safe
	}
	return errutil.Wrap(err, "failed to rollback transaction")
}
