package repositoryengine

import (
	"context"

	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
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
	panic("implement me")
}

func (r *paymentRepository) Store(ctx context.Context, payment *repository.Payment) error {
	panic("implement me")
}
