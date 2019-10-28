package repositoryengine

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

// accountRepository implements AccountRepository interface
type accountRepository struct {
	db *sqlx.DB
}

// newAccountRepository creates new account repository
func newAccountRepository(db *sqlx.DB) repository.AccountRepository {
	return &accountRepository{
		db: db,
	}
}
