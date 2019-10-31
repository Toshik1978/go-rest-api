package account

import (
	"context"
	"time"

	"github.com/Toshik1978/go-rest-api/service/errutil"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/Toshik1978/go-rest-api/service/validator"
	"go.uber.org/zap"
)

// paymentBuilder implements PaymentBuilder interface
type paymentBuilder struct {
	logger            *zap.Logger
	repositoryFactory repository.Factory

	payment repository.Payment
	v       *validator.Validator
}

// newPaymentBuilder creates new PaymentBuilder implementation
func newPaymentBuilder(globals server.Globals) handler.PaymentBuilder {
	return &paymentBuilder{
		logger:            globals.Logger,
		repositoryFactory: globals.RepositoryFactory,
	}
}

func (b *paymentBuilder) SetAmount(amount float64) handler.PaymentBuilder {
	b.payment.Amount = int64(amount * 100)
	b.v.ValidateAmount(amount)
	return b
}

func (b *paymentBuilder) SetPayer(uid string) handler.PaymentBuilder {
	b.payment.PayerAccountUID = uid
	b.v.ValidateUID("payer_uid", uid)
	return b
}

func (b *paymentBuilder) SetRecipient(uid string) handler.PaymentBuilder {
	b.payment.RecipientAccountUID = uid
	b.v.ValidateUID("recipient_uid", uid)
	return b
}

func (b *paymentBuilder) Build(ctx context.Context) (*handler.Payment, error) {
	if err := b.v.Error(); err != nil {
		return nil, handler.WrapError(err, "failed to validate payment", handler.ClientError)
	}
	b.payment.CreatedAt = time.Now()

	scope := b.repositoryFactory.Scope()
	ctx, err := scope.WithContext(ctx)
	if err != nil {
		return nil, errutil.Wrap(err, "failed to start repository scope")
	}
	// Here we can defer Cancel operation, because it's safe
	defer func() { _ = scope.Cancel(ctx) }()

	// We should create new payment and update balance for accounts
	if err := b.storePayment(ctx); err != nil {
		return nil, handler.WrapError(err, "failed to create payment", handler.ServerError)
	}
	if err := b.updateBalance(ctx); err != nil {
		return nil, handler.WrapError(err, "failed to update balance", handler.ServerError)
	}

	// Complete scope
	if err := scope.Complete(ctx); err != nil {
		return nil, errutil.Wrap(err, "failed to complete repository scope")
	}
	return mapRepositoryPayment(b.payment), nil
}

// storePayment store payments in storage
// For finance purpose is not bad idea to store always 2 transactions per payment:
// 1. Incoming - payment to recipient with positive amount
// 2. Outgoing - recipient to payment with negative amount
// It cause more simple balance calculation operations
func (b *paymentBuilder) storePayment(ctx context.Context) error {
	if err := b.repositoryFactory.PaymentRepository().Store(ctx, &b.payment); err != nil {
		return err
	}
	if err := b.repositoryFactory.PaymentRepository().Store(ctx, b.reversePayment()); err != nil {
		return err
	}
	return nil
}

// reversePayment create reverse payment from current payment
func (b *paymentBuilder) reversePayment() *repository.Payment {
	payment := b.payment
	payer := payment.PayerAccountUID
	payment.PayerAccountUID = payment.RecipientAccountUID
	payment.RecipientAccountUID = payer
	payment.Amount = -payment.Amount
	return &payment
}

// updateBalance updates balances in storage
func (b *paymentBuilder) updateBalance(ctx context.Context) error {
	if err := b.repositoryFactory.AccountRepository().
		UpdateBalance(ctx, b.payment.PayerAccountUID, -b.payment.Amount); err != nil {

		return err
	}
	if err := b.repositoryFactory.AccountRepository().
		UpdateBalance(ctx, b.payment.RecipientAccountUID, b.payment.Amount); err != nil {

		return err
	}
	return nil
}
