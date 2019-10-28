package postgres

import (
	"time"

	"github.com/Toshik1978/go-rest-api/service"
	_ "github.com/jackc/pgx/v4/stdlib" // PostgreSQL driver
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const (
	dbSuffix           = ""
	connectLifetime    = 45 * time.Second
	maxIdleConnections = 20
	maxOpenConnections = 20
)

// postgresClient implements PostgresClient
type postgresClient struct {
	db         *sqlx.DB
	timeout    time.Duration
	logger     *zap.Logger
	stopNotify chan struct{}
}

// NewPostgresClient creates new PostgreSQL client
func NewPostgresClient(logger *zap.Logger, vars service.Vars) (service.PostgresClient, error) {
	db, err := sqlx.Connect("pgx", vars.DB+dbSuffix)
	if err != nil {
		return nil, errors.Wrap(err, "db initialization failed")
	}
	initializeConnection(db)

	client := &postgresClient{
		db:         db,
		timeout:    vars.DBTimeout,
		logger:     logger,
		stopNotify: make(chan struct{}),
	}
	client.backgroundReconnect()
	return client, nil
}

// initializeConnection initializes connection parameters
func initializeConnection(db *sqlx.DB) {
	db.SetConnMaxLifetime(connectLifetime)
	db.SetMaxOpenConns(maxOpenConnections)
	db.SetMaxIdleConns(maxIdleConnections)
}

// Stop stops background job
func (c *postgresClient) Stop() {
	if c.stopNotify != nil {
		select {
		case <-c.stopNotify:
		default:
			close(c.stopNotify)
		}
	}

	if c.db != nil {
		c.db.Close()
		c.db = nil
	}
}

// GetConnections retrieve DB connections
func (c *postgresClient) GetConnection() *sqlx.DB {
	return c.db
}

// backgroundReconnect process background reconnect
func (c *postgresClient) backgroundReconnect() {
	c.logger.Info("PostgreSQL start background reconnect")

	ticker := time.NewTicker(c.timeout)
	go func() {
		for {
			if err := c.db.Ping(); err != nil {
				c.logger.Error("Background reconnect for database failed", zap.Error(err))
			}

			select {
			case <-c.stopNotify:
				c.logger.Info("PostgreSQL stop background reconnect")
				return
			case <-ticker.C:
			}
		}
	}()
}
