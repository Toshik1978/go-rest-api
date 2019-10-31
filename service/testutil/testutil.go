package testutil

import (
	"time"

	"github.com/Toshik1978/go-rest-api/repository"
)

func RepositoryPayment() repository.Payment {
	return repository.Payment{
		ID:                  1234,
		Amount:              100,
		PayerAccountUID:     "toshik1978",
		RecipientAccountUID: "toshik1979",
		CreatedAt:           time.Now(),
	}
}
