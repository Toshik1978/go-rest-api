package repositoryengine

import (
	"context"

	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

const (
	getAllPaymentsSQL = `
		SELECT id, amount, payer_account_uid, recipient_account_uid, created_at
		FROM payments`

	storePaymentSQL = `
		INSERT INTO payments
			(amount, payer_account_uid, recipient_account_uid, created_at)
		VALUES
			(:amount, :payer_account_uid, :recipient_account_uid, :created_at)`
)

// paymentRepository implements PaymentRepository interface
type paymentRepository struct {
	ext sqlx.Ext
}

// newPaymentRepository creates new payment repository
func newPaymentRepository(ext sqlx.Ext) repository.PaymentRepository {
	return &paymentRepository{
		ext: ext,
	}
}

func (r *paymentRepository) GetAll(ctx context.Context) ([]repository.Payment, error) {
	var payments []repository.Payment
	if err := sqlx.Select(sqlxExt(ctx, r.ext), &payments, getAllPaymentsSQL); err != nil {
		return nil, err
	}
	return payments, nil
}

func (r *paymentRepository) Store(ctx context.Context, payment *repository.Payment) error {
	res, err := sqlx.NamedExec(sqlxExt(ctx, r.ext), storePaymentSQL, payment)
	if err != nil {
		return err
	}
	payment.ID, _ = res.LastInsertId()
	return nil
}
