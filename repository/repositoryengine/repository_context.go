package repositoryengine

import (
	"database/sql"

	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/errutil"
	"github.com/jmoiron/sqlx"
)

// repositoryContext implements repository.Context interface
type repositoryContext struct {
	tx *sqlx.Tx
}

// newRepositoryContext creates new instance of repository.Context interface
func newRepositoryContext(tx *sqlx.Tx) repository.Context {
	return &repositoryContext{
		tx: tx,
	}
}

func (c *repositoryContext) Cancel() error {
	err := c.tx.Rollback()
	if err == sql.ErrTxDone {
		return nil // Ignore this error, because it's safe
	}
	return errutil.Wrap(err, "failed to rollback transaction")
}

func (c *repositoryContext) Complete() error {
	err := c.tx.Commit()
	if err == sql.ErrTxDone {
		return nil // Ignore this error, because it's safe
	}
	return errutil.Wrap(err, "failed to commit transaction")
}
