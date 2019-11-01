package account

import (
	"context"
	"strings"
	"time"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/Toshik1978/go-rest-api/service/validator"
	"go.uber.org/zap"
)

// accountBuilder implements AccountBuilder interface
type accountBuilder struct {
	logger            *zap.Logger
	repositoryFactory repository.Factory

	account repository.Account
	v       *validator.Validator
}

// newAccountBuilder creates new AccountBuilder implementation
func newAccountBuilder(globals server.Globals) handler.AccountBuilder {
	return &accountBuilder{
		logger:            globals.Logger,
		repositoryFactory: globals.RepositoryFactory,
		account:           repository.Account{CreatedAt: time.Now()},
		v:                 validator.NewValidator(),
	}
}

func (b *accountBuilder) SetUID(uid string) handler.AccountBuilder {
	b.account.UID = uid
	return b
}

func (b *accountBuilder) SetBalance(balance float64) handler.AccountBuilder {
	b.account.Balance = int64(balance * 100)
	return b
}

func (b *accountBuilder) SetCurrency(currency string) handler.AccountBuilder {
	b.account.Currency = strings.ToUpper(currency)
	return b
}

func (b *accountBuilder) Build(ctx context.Context) (*handler.Account, error) {
	b.v.
		ValidateUID("uid", b.account.UID).
		ValidateBalance(float64(b.account.Balance) / 100).
		ValidateCurrency(b.account.Currency)
	if err := b.v.Error(); err != nil {
		return nil, handler.WrapError(err, "failed to validate account", handler.ClientError)
	}

	if err := b.repositoryFactory.AccountRepository().Store(ctx, &b.account); err != nil {
		return nil, handler.WrapError(err, "failed to create account", handler.ServerError)
	}
	return mapRepositoryAccount(b.account), nil
}
