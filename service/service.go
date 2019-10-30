package service

import (
	"github.com/jmoiron/sqlx"
)

//go:generate mockgen -source service.go -package mock -destination ../mock/service.go

// PostgresClient declare interface for PostgreSQL connections
type PostgresClient interface {
	// GetConnection retrieve database connection to be used in repositories
	GetConnection() *sqlx.DB
	// Stop finish background tasks of this client (db reconnection)
	Stop()
}
