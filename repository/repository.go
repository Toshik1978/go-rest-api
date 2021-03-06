package repository

import (
	"context"
)

//go:generate mockgen -source repository.go -package mock -destination ../mock/repository.go

// AccountRepository declare repository for accounts
type AccountRepository interface {
	// GetAll return all accounts in storage
	GetAll(ctx context.Context) ([]Account, error)
	// Store save new account in storage
	Store(ctx context.Context, account *Account) error
	// Update balance for given account by incrementing on given value
	UpdateBalance(ctx context.Context, uid string, incr int64) error
}

// PaymentRepository declare repository for payments
type PaymentRepository interface {
	// GetAll return all payments in storage
	GetAll(ctx context.Context) ([]Payment, error)
	// Store save new payment in storage
	Store(ctx context.Context, payment *Payment) error
}

// Repository pattern and transactions are not very good combination, so here we are declare some scope.
// It has semantic of unit of work, calling code should not know about nature of scope,
// but code can cancel or complete it.

// Scope define some operation context for repository operations (unit of work)
// It's safe to call Cancel/Complete multiple times, only the first one will be actually done
type Scope interface {
	// WithContext initializes scope with context and return new context to use in repository's operations
	WithContext(ctx context.Context) (context.Context, error)

	// Complete finish current scope
	Complete(ctx context.Context) error
	// Cancel cancel current scope
	Cancel(ctx context.Context) error
}

// Factory define accessors to all repositories
type Factory interface {
	// Context creates new scope for repository activities
	Scope() Scope

	// AccountRepository return account repository instance
	AccountRepository() AccountRepository
	// PaymentRepository return payment repository instance
	PaymentRepository() PaymentRepository
}
