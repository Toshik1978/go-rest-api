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

	repositoryContext, err := b.repositoryFactory.Context()
	if err != nil {
		return nil, errutil.Wrap(err, "failed to create repository context")
	}
	// Here we can defer Cancel operation, because it's safe
	defer func() { _ = repositoryContext.Cancel() }()

	// We should create new payment and update balance for accounts
	if err := b.repositoryFactory.PaymentRepository(repositoryContext).Store(ctx, &b.payment); err != nil {
		return nil, handler.WrapError(err, "failed to create payment", handler.ServerError)
	}
	if err := b.updateBalance(ctx, repositoryContext); err != nil {
		return nil, handler.WrapError(err, "failed to update balance", handler.ServerError)
	}

	if err := repositoryContext.Complete(); err != nil {
		return nil, errutil.Wrap(err, "failed to complete repository context")
	}
	return mapRepositoryPayment(b.payment, true), nil
}

// updateBalance updates balance for both accounts - payer and recipient
func (b *paymentBuilder) updateBalance(ctx context.Context, repositoryContext repository.Context) error {
	repo := b.repositoryFactory.AccountRepository(repositoryContext)
	if err := repo.IncrementBalance(ctx, b.payment.PayerAccountUID, -b.payment.Amount); err != nil {
		return err
	}
	if err := repo.IncrementBalance(ctx, b.payment.RecipientAccountUID, b.payment.Amount); err != nil {
		return err
	}
	return nil
}
