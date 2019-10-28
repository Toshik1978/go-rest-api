package service

import (
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source service.go -package mock -destination ../mock/service.go

// PostgresClient declare interface for PostgreSQL connections
type PostgresClient interface {
	Stop()

	GetConnection() *sqlx.DB
}
