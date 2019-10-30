package repositoryengine

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/errutil"
	"github.com/jmoiron/sqlx"
)

// repositoryFactory implements RepositoryFactory interface
type repositoryFactory struct {
	db *sqlx.DB
}

// NewRepositoryFactory creates repository factory
func NewRepositoryFactory(db *sqlx.DB) repository.Factory {
	return &repositoryFactory{
		db: db,
	}
}

func (f *repositoryFactory) Context() (repository.Context, error) {
	tx, err := f.db.Beginx()
	if err != nil {
		return nil, errutil.Wrap(err, "failed to create repository context")
	}
	return newRepositoryContext(tx), nil
}

func (f *repositoryFactory) AccountRepository(ctx repository.Context) repository.AccountRepository {
	if txCtx, ok := ctx.(*repositoryContext); ok {
		return newAccountRepository(txCtx.tx)
	}
	return newAccountRepository(f.db)
}

func (f *repositoryFactory) PaymentRepository(ctx repository.Context) repository.PaymentRepository {
	if txCtx, ok := ctx.(*repositoryContext); ok {
		return newPaymentRepository(txCtx.tx)
	}
	return newPaymentRepository(f.db)
}
