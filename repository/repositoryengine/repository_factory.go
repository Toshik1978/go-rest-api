package repositoryengine

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

// repositoryFactory implements RepositoryFactory interface
type repositoryFactory struct {
	accountRepository repository.AccountRepository
	paymentRepository repository.PaymentRepository
}

// NewRepositoryFactory creates repository factory
func NewRepositoryFactory(db *sqlx.DB) repository.Factory {
	return &repositoryFactory{
		accountRepository: newAccountRepository(db),
		paymentRepository: newPaymentRepository(db),
	}
}

// GetAccountRepository returns account repository
func (f *repositoryFactory) GetAccountRepository() repository.AccountRepository {
	return f.accountRepository
}

// GetPaymentRepository returns payment repository
func (f *repositoryFactory) GetPaymentRepository() repository.PaymentRepository {
	return f.paymentRepository
}
