package repository

import (
	"context"
	"errors"
)

//go:generate mockgen -source repository.go -package mock -destination ../mock/repository.go

// Errors
var (
	ErrNotFound = errors.New("data not found in repository")
)

// AccountRepository declare repository for accounts
type AccountRepository interface {
	// GetAll return all accounts in storage
	GetAll(ctx context.Context) ([]Account, error)
	// Store save new account in storage
	Store(ctx context.Context, account *Account) error
	// Update balance for given account by incrementing on given value
	IncrementBalance(ctx context.Context, uid string, incr int64) error
}

// PaymentRepository declare repository for payments
type PaymentRepository interface {
	// GetAll return all payments in storage
	GetAll(ctx context.Context) ([]Payment, error)
	// Store save new payment in storage
	Store(ctx context.Context, payment *Payment) error
}

// Repository pattern and transactions are not very good combination, so here we are declare some context.
// It has semantic of unit of work, calling code should not know about nature of context,
// but code can cancel or complete unit of work.

// Context define some operation context for repository operations (unit of work)
// It's safe to call Cancel/Complete multiple times, only the first one will be actually done
type Context interface {
	// Cancel cancel current context and all unit of work
	Cancel() error
	// Complete finish current unit of work
	Complete() error
}

// Factory define access to all repositories
type Factory interface {
	// Context creates new unit of work context
	Context() (Context, error)

	// AccountRepository return account repository instance
	// ctx can be nil to create repository w/o context
	AccountRepository(ctx Context) AccountRepository
	// PaymentRepository return payment repository instance
	// ctx can be nil to create repository w/o context
	PaymentRepository(ctx Context) PaymentRepository
}
