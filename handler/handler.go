package handler

import "context"

//go:generate mockgen -source handler.go -package mock -destination ../mock/handler.go

// AccountBuilder declare interface to build new account
type AccountBuilder interface {
	// SetUID initializes UID for the new account
	SetUID(uid string) AccountBuilder
	// SetBalance initializes initial balance for the new account
	SetBalance(balance float64) AccountBuilder
	// SetCurrency initializes currency for the new account
	SetCurrency(currency string) AccountBuilder

	// Build actually creates new account
	Build(ctx context.Context) (*Account, error)
}

// PaymentBuilder declare interface to build new payment
type PaymentBuilder interface {
	// SetAmount initializes amount for the new payment
	SetAmount(amount float64) PaymentBuilder
	// SetPayer initializes payer for the new payment
	SetPayer(uid string) PaymentBuilder
	// SetRecipient initializes recipient for the new payment
	SetRecipient(uid string) PaymentBuilder

	// Build actually creates new payment
	Build(ctx context.Context) (*Payment, error)
}

// AccountFactory declare interface to access accounts and payments information (kind of facade to simplify interface)
type AccountManager interface {
	// AllAccounts return all available accounts in the system
	AllAccounts(ctx context.Context) ([]Account, error)
	// AllPayments return all available payments in the system
	AllPayments(ctx context.Context) ([]Payment, error)

	// AccountBuilder instantiate new account builder
	AccountBuilder() AccountBuilder
	// PaymentBuilder instantiate new payment builder
	PaymentBuilder() PaymentBuilder
}
