package repositoryengine

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/jmoiron/sqlx"
)

// paymentRepository implements PaymentRepository interface
type paymentRepository struct {
	db *sqlx.DB
}

// newPaymentRepository creates new payment repository
func newPaymentRepository(db *sqlx.DB) repository.PaymentRepository {
	return &paymentRepository{
		db: db,
	}
}
