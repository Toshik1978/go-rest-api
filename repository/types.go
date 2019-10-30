package repository

import "time"

// Account define account entity
type Account struct {
	ID        int64     `db:"id"`
	UID       string    `db:"uid"`
	Currency  string    `db:"currency"`
	Balance   int64     `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
}

// Payment define payment entity
type Payment struct {
	ID                  int64     `db:"id"`
	Amount              int64     `db:"amount"`
	PayerAccountUID     string    `db:"payer_account_uid"`
	RecipientAccountUID string    `db:"recipient_account_uid"`
	CreatedAt           time.Time `db:"created_at"`
}
