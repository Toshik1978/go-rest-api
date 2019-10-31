package account

import (
	"context"
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
		v:                 validator.NewValidator(),
	}
}

func (b *accountBuilder) SetUID(uid string) handler.AccountBuilder {
	b.account.UID = uid
	b.v.ValidateUID("uid", uid)
	return b
}

func (b *accountBuilder) SetBalance(balance float64) handler.AccountBuilder {
	b.account.Balance = int64(balance * 100)
	b.v.ValidateBalance(balance)
	return b
}

func (b *accountBuilder) SetCurrency(currency string) handler.AccountBuilder {
	b.account.Currency = currency
	b.v.ValidateCurrency(currency)
	return b
}

func (b *accountBuilder) Build(ctx context.Context) (*handler.Account, error) {
	if err := b.v.Error(); err != nil {
		return nil, handler.WrapError(err, "failed to validate account", handler.ClientError)
	}
	b.account.CreatedAt = time.Now()
	if err := b.repositoryFactory.AccountRepository().Store(ctx, &b.account); err != nil {
		return nil, handler.WrapError(err, "failed to create account", handler.ServerError)
	}
	return mapRepositoryAccount(b.account), nil
}
