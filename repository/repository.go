package repository

import (
	"errors"
)

//go:generate mockgen -source repository.go -package mock -destination ../mock/repository.go

// Errors
var (
	ErrNotFound = errors.New("data not found in repository")
)

// AccountRepository declare repository for accounts
type AccountRepository interface {
}

// PaymentRepository declare repository for payments
type PaymentRepository interface {
}

// Factory define access to all repositories
type Factory interface {
	GetAccountRepository() AccountRepository
	GetPaymentRepository() PaymentRepository
}
