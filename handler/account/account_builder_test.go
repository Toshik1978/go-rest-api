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

type accountBuilderTestSuite struct {
	suite.Suite

	account repository.Account
}

func (s *accountBuilderTestSuite) SetupSuite() {
	s.account = testutil.RepositoryAccount()
}

func (s *accountBuilderTestSuite) TestAccountBuilderValidationFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newAccountBuilder(server.Globals{
		Logger: zap.New(zapCore),
	})

	account, err := builder.
		SetCurrency(s.account.Currency).
		SetBalance(float64(s.account.Balance) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(account)
	s.Equal(0, zapRecorded.Len())
}

func (s *accountBuilderTestSuite) TestAccountBuilderStoreFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repository := mock.NewMockAccountRepository(ctrl)
	repository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryAccount(s.account)).
		Return(errors.New("fail"))
	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		AccountRepository().
		Return(repository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newAccountBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	account, err := builder.
		SetUID(s.account.UID).
		SetCurrency(s.account.Currency).
		SetBalance(float64(s.account.Balance) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(account)
	s.Equal(0, zapRecorded.Len())
}

func (s *accountBuilderTestSuite) TestAccountBuilderSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	repository := mock.NewMockAccountRepository(ctrl)
	repository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryAccount(s.account)).
		Return(nil)
	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		AccountRepository().
		Return(repository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newAccountBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	account, err := builder.
		SetUID(s.account.UID).
		SetCurrency(s.account.Currency).
		SetBalance(float64(s.account.Balance) / 100).
		Build(context.Background())

	s.NoError(err)
	s.NotNil(account)
	s.Equal(s.account.UID, account.UID)
	s.Equal(s.account.Currency, account.Currency)
	s.Equal(0, zapRecorded.Len())
}
