package account

import (
	"context"
	"errors"

	"github.com/Toshik1978/go-rest-api/mock"
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/Toshik1978/go-rest-api/service/testutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type accountManagerTestSuite struct {
	suite.Suite

	accounts []repository.Account
	payments []repository.Payment
}

func (s *accountManagerTestSuite) SetupSuite() {
	account1 := testutil.RepositoryAccount()
	account2 := account1
	account2.UID = "toshik1979"
	s.accounts = []repository.Account{account1, account2}

	payment1 := testutil.RepositoryPayment()
	payment2 := payment1
	payment2.PayerAccountUID = payment1.RecipientAccountUID
	payment2.RecipientAccountUID = payment1.PayerAccountUID
	payment2.Amount = -payment1.Amount
	s.payments = []repository.Payment{payment1, payment2}
}

func (s *accountManagerTestSuite) TestAllAccountsFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repository := mock.NewMockAccountRepository(ctrl)
	repository.
		EXPECT().
		GetAll(gomock.Any()).
		Return(nil, errors.New("fail"))
	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		AccountRepository().
		Return(repository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	accountManager := NewAccountManager(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})
	accounts, err := accountManager.AllAccounts(context.Background())

	s.Error(err)
	s.Nil(accounts)
	s.Equal(0, zapRecorded.Len())
}

func (s *accountManagerTestSuite) TestAllAccountsSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repository := mock.NewMockAccountRepository(ctrl)
	repository.
		EXPECT().
		GetAll(gomock.Any()).
		Return(s.accounts, nil)
	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		AccountRepository().
		Return(repository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	accountManager := NewAccountManager(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})
	accounts, err := accountManager.AllAccounts(context.Background())

	s.NoError(err)
	s.Len(accounts, 2)
	s.Equal(0, zapRecorded.Len())
}

func (s *accountManagerTestSuite) TestAllPaymentsFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repository := mock.NewMockPaymentRepository(ctrl)
	repository.
		EXPECT().
		GetAll(gomock.Any()).
		Return(nil, errors.New("fail"))
	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		PaymentRepository().
		Return(repository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	accountManager := NewAccountManager(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})
	payments, err := accountManager.AllPayments(context.Background())

	s.Error(err)
	s.Nil(payments)
	s.Equal(0, zapRecorded.Len())
}

func (s *accountManagerTestSuite) TestAllPaymentsSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repository := mock.NewMockPaymentRepository(ctrl)
	repository.
		EXPECT().
		GetAll(gomock.Any()).
		Return(s.payments, nil)
	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		PaymentRepository().
		Return(repository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	accountManager := NewAccountManager(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})
	payments, err := accountManager.AllPayments(context.Background())

	s.NoError(err)
	s.Len(payments, 2)
	s.Equal(s.payments[0].PayerAccountUID, payments[0].UID)
	s.NotNil(payments[0].TargetUID)
	s.Equal(s.payments[0].RecipientAccountUID, *payments[0].TargetUID)
	s.Equal(outgoingPayment, payments[0].Direction)
	s.Equal(s.payments[1].PayerAccountUID, payments[1].UID)
	s.NotNil(payments[1].SourceUID)
	s.Equal(s.payments[1].RecipientAccountUID, *payments[1].SourceUID)
	s.Equal(incomingPayment, payments[1].Direction)
	s.Equal(0, zapRecorded.Len())
}

func (s *accountManagerTestSuite) TestAccountBuilderSucceeded() {
	accountManager := NewAccountManager(server.Globals{})
	builder := accountManager.AccountBuilder()

	s.NotNil(accountManager)
	s.IsType(&accountBuilder{}, builder)
}

func (s *accountManagerTestSuite) TestPaymentBuilderSucceeded() {
	accountManager := NewAccountManager(server.Globals{})
	builder := accountManager.PaymentBuilder()

	s.NotNil(accountManager)
	s.IsType(&paymentBuilder{}, builder)
}
