package account

import (
	"context"

	"github.com/Toshik1978/go-rest-api/service/errutil"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/server"
	"go.uber.org/zap"
)

// accountManager implements AccountManager interface
type accountManager struct {
	logger            *zap.Logger
	repositoryFactory repository.Factory
}

// NewAccountManager creates new implementation of AccountManager interface
func NewAccountManager(globals server.Globals) handler.AccountManager {
	return &accountManager{
		logger:            globals.Logger,
		repositoryFactory: globals.RepositoryFactory,
	}
}

func (m *accountManager) AllAccounts(ctx context.Context) ([]handler.Account, error) {
	accounts, err := m.repositoryFactory.AccountRepository(nil).GetAll(ctx)
	if err != nil {
		return nil, errutil.Wrap(err, "failed to get accounts")
	}
	return mapRepositoryAccounts(accounts), nil
}

func (m *accountManager) AllPayments(ctx context.Context) ([]handler.Payment, error) {
	payments, err := m.repositoryFactory.PaymentRepository(nil).GetAll(ctx)
	if err != nil {
		return nil, errutil.Wrap(err, "failed to get payments")
	}
	return mapRepositoryPayments(payments), nil
}

func (m *accountManager) AccountBuilder() handler.AccountBuilder {
	return newAccountBuilder(server.Globals{
		Logger:            m.logger,
		RepositoryFactory: m.repositoryFactory,
	})
}

func (m *accountManager) PaymentBuilder() handler.PaymentBuilder {
	return newPaymentBuilder(server.Globals{
		Logger:            m.logger,
		RepositoryFactory: m.repositoryFactory,
	})
}
