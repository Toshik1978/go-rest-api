package repository

import "time"

// Account define account entity
type Account struct {
	ID        int64     `db:"id"`
	UUID      string    `db:"uuid"`
	Balance   int64     `db:"balance"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
}

// Payment define payment entity
type Payment struct {
	ID        int64     `db:"id"`
	Amount    int64     `db:"amount"`
	AccountID int64     `db:"account_id"`
	Direction int8      `db:"direction"`
	CreatedAt time.Time `db:"created_at"`
}
