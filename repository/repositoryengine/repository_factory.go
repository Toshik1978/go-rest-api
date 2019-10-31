package repositoryengine

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

// repositoryFactory implements RepositoryFactory interface
type repositoryFactory struct {
	db                *sqlx.DB
	accountRepository repository.AccountRepository
	paymentRepository repository.PaymentRepository
}

// NewRepositoryFactory creates repository factory
func NewRepositoryFactory(db *sqlx.DB) repository.Factory {
	return &repositoryFactory{
		db:                db,
		accountRepository: newAccountRepository(db),
		paymentRepository: newPaymentRepository(db),
	}
}

func (f *repositoryFactory) Scope() repository.Scope {
	return newScope(f.db)
}

func (f *repositoryFactory) AccountRepository() repository.AccountRepository {
	return f.accountRepository
}

func (f *repositoryFactory) PaymentRepository() repository.PaymentRepository {
	return f.paymentRepository
}
