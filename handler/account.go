package handler

import "time"

// AccountRequest define request to create new account
type AccountRequest struct {
	UID      string  `json:"uid"`
	Currency string  `json:"currency"`
	Balance  float64 `json:"balance"`
}

// Account define account description
type Account struct {
	UID       string    `json:"uid"`
	Currency  string    `json:"currency"`
	Balance   float64   `json:"balance"`
	CreatedAt time.Time `json:"created_at"`
}

// PaymentRequest define request to create new payment
type PaymentRequest struct {
	RecipientUID string  `json:"recipient"`
	Amount       float64 `json:"amount"`
}

// Payment define payment description
type Payment struct {
	UID       string    `json:"account"`
	SourceUID *string   `json:"from_account,omitempty"`
	TargetUID *string   `json:"to_account,omitempty"`
	Direction string    `json:"direction"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at"`
}
