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

type paymentBuilderTestSuite struct {
	suite.Suite

	payments []repository.Payment
}

func (s *paymentBuilderTestSuite) SetupSuite() {
	payment1 := testutil.RepositoryPayment()
	payment2 := payment1
	payment2.PayerAccountUID = payment1.RecipientAccountUID
	payment2.RecipientAccountUID = payment1.PayerAccountUID
	payment2.Amount = -payment1.Amount
	s.payments = []repository.Payment{payment1, payment2}
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderValidationFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger: zap.New(zapCore),
	})

	payment, err := builder.
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderWithContextFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(nil, errors.New("fail"))

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentStoreFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(nil)

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(errors.New("fail"))

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderStoreReverseFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(nil)

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(nil)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[1])).
		Return(errors.New("fail"))

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository).
		Times(2)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderUpdatePayerBalanceFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(nil)

	accountRepository := mock.NewMockAccountRepository(ctrl)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[0].PayerAccountUID), gomock.Eq(-s.payments[0].Amount)).
		Return(errors.New("fail"))

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(nil)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[1])).
		Return(nil)

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		AccountRepository().
		Return(accountRepository)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository).
		Times(2)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderUpdateRecepientBalanceFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(nil)

	accountRepository := mock.NewMockAccountRepository(ctrl)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[0].PayerAccountUID), gomock.Eq(-s.payments[0].Amount)).
		Return(nil)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[1].PayerAccountUID), gomock.Eq(-s.payments[1].Amount)).
		Return(errors.New("fail"))

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(nil)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[1])).
		Return(nil)

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		AccountRepository().
		Return(accountRepository).
		Times(2)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository).
		Times(2)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderCompleteFailed() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Complete(gomock.Any()).
		Return(errors.New("fail"))
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(nil)

	accountRepository := mock.NewMockAccountRepository(ctrl)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[0].PayerAccountUID), gomock.Eq(-s.payments[0].Amount)).
		Return(nil)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[1].PayerAccountUID), gomock.Eq(-s.payments[1].Amount)).
		Return(nil)

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(nil)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[1])).
		Return(nil)

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		AccountRepository().
		Return(accountRepository).
		Times(2)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository).
		Times(2)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.Error(err)
	s.Nil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderCancelFailedButSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Complete(gomock.Any()).
		Return(nil)
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(errors.New("fail"))

	accountRepository := mock.NewMockAccountRepository(ctrl)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[0].PayerAccountUID), gomock.Eq(-s.payments[0].Amount)).
		Return(nil)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[1].PayerAccountUID), gomock.Eq(-s.payments[1].Amount)).
		Return(nil)

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(nil)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[1])).
		Return(nil)

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		AccountRepository().
		Return(accountRepository).
		Times(2)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository).
		Times(2)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.NoError(err)
	s.NotNil(payment)
	s.Equal(0, zapRecorded.Len())
}

func (s *paymentBuilderTestSuite) TestPaymentBuilderSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	scope := mock.NewMockScope(ctrl)
	scope.
		EXPECT().
		WithContext(gomock.Any()).
		Return(context.Background(), nil)
	scope.
		EXPECT().
		Complete(gomock.Any()).
		Return(nil)
	scope.
		EXPECT().
		Cancel(gomock.Any()).
		Return(nil)

	accountRepository := mock.NewMockAccountRepository(ctrl)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[0].PayerAccountUID), gomock.Eq(-s.payments[0].Amount)).
		Return(nil)
	accountRepository.
		EXPECT().
		UpdateBalance(gomock.Any(), gomock.Eq(s.payments[1].PayerAccountUID), gomock.Eq(-s.payments[1].Amount)).
		Return(nil)

	paymentRepository := mock.NewMockPaymentRepository(ctrl)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[0])).
		Return(nil)
	paymentRepository.
		EXPECT().
		Store(gomock.Any(), testutil.EqualRepositoryPayment(s.payments[1])).
		Return(nil)

	factory := mock.NewMockFactory(ctrl)
	factory.
		EXPECT().
		Scope().
		Return(scope)
	factory.
		EXPECT().
		AccountRepository().
		Return(accountRepository).
		Times(2)
	factory.
		EXPECT().
		PaymentRepository().
		Return(paymentRepository).
		Times(2)

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	builder := newPaymentBuilder(server.Globals{
		Logger:            zap.New(zapCore),
		RepositoryFactory: factory,
	})

	payment, err := builder.
		SetPayer(s.payments[0].PayerAccountUID).
		SetRecipient(s.payments[0].RecipientAccountUID).
		SetAmount(float64(s.payments[0].Amount) / 100).
		Build(context.Background())

	s.NoError(err)
	s.NotNil(payment)
	s.Equal(0, zapRecorded.Len())
}
