package repositoryengine

import (
	"context"

	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

const (
	getAllAccountsSQL = `
		SELECT id, uid, currency, balance, created_at
		FROM accounts`

	storeAccountSQL = `
		INSERT INTO accounts
			(uid, currency, balance, created_at)
		VALUES
			(:uid, :currency, :balance, :created_at)`
	updateBalanceSQL = `
		UPDATE accounts
		SET balance = balance + $2
		WHERE uid = $1`
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
	var accounts []repository.Account
	if err := sqlx.Select(sqlxExt(ctx, r.ext), &accounts, getAllAccountsSQL); err != nil {
		return nil, err
	}
	return accounts, nil
}

func (r *accountRepository) Store(ctx context.Context, account *repository.Account) error {
	res, err := sqlx.NamedExec(sqlxExt(ctx, r.ext), storeAccountSQL, account)
	if err != nil {
		return err
	}
	account.ID, _ = res.LastInsertId()
	return nil
}

func (r *accountRepository) UpdateBalance(ctx context.Context, uid string, incr int64) error {
	_, err := sqlxExt(ctx, r.ext).Exec(updateBalanceSQL, uid, incr)
	if err != nil {
		return err
	}
	return nil
}
