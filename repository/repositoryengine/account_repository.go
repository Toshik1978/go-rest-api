package repositoryengine

import (
	"context"

	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

// accountRepository implements AccountRepository interface
type accountRepository struct {
	ext sqlx.Ext
}

// newAccountRepository creates new account repository
func newAccountRepository(ext sqlx.Ext) repository.AccountRepository {
	return &accountRepository{
		ext: ext,
	}
}

func (r *accountRepository) GetAll(ctx context.Context) ([]repository.Account, error) {
	panic("implement me")
}

func (r *accountRepository) Store(ctx context.Context, account *repository.Account) error {
	panic("implement me")
}

func (r *accountRepository) IncrementBalance(ctx context.Context, uid string, incr int64) error {
	panic("implement me")
}
