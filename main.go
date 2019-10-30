package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/handler/account"
	"github.com/Toshik1978/go-rest-api/handler/httphandler"
	"github.com/Toshik1978/go-rest-api/repository/repositoryengine"
	"github.com/Toshik1978/go-rest-api/service"
	"github.com/Toshik1978/go-rest-api/service/postgres"
	"github.com/Toshik1978/go-rest-api/service/server"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	httpShutdownTimeout = 5 * time.Second
)

var (
	BuildTime  = "undefined"
	GitVersion = "undefined"
)

func main() {
	logger := initializeLogger()
	defer func() {
		if recErr := recover(); recErr != nil {
			// Log error
			logger.Error("Panic in main", zap.Any("panic", recErr))
		}
	}()

	// Initialize signals
	interruptCh := make(chan os.Signal, 1)
	signal.Notify(interruptCh, syscall.SIGINT, syscall.SIGTERM)

	logger.Info("Start service", zap.String("git_version", GitVersion))
	vars := server.LoadConfig(logger)

	dbClient := initializeDB(logger, vars)
	globals := initializeGlobals(logger, dbClient)
	accountManager := account.NewAccountManager(globals)
	server := initializeHTTP(vars, globals, accountManager)

	waitShutdown(interruptCh, logger, dbClient, server)
}

// initializeLogger initialized logger
func initializeLogger() *zap.Logger {
	config := zap.NewProductionConfig()
	config.DisableCaller = true
	config.DisableStacktrace = true
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		log.Fatal("Initialize logger failed", err)
	}
	return logger
}

// initializeDB initializes DB
func initializeDB(logger *zap.Logger, vars server.Vars) service.PostgresClient {
	client, err := postgres.NewPostgresClient(logger, vars)
	if err != nil {
		logger.Fatal("DB initialization failed", zap.Error(err))
		return nil
	}
	logger.Info("DB initialized")
	return client
}

// initializeGlobals initialize globals
func initializeGlobals(logger *zap.Logger, dbClient service.PostgresClient) server.Globals {
	db := dbClient.GetConnection()
	return server.Globals{
		Logger:            logger,
		RepositoryFactory: repositoryengine.NewRepositoryFactory(db),
		BuildTime:         BuildTime,
		Version:           GitVersion,
	}
}

// initializeHTTP initializes HTTP server
func initializeHTTP(vars server.Vars, globals server.Globals, accountManager handler.AccountManager) *http.Server {
	server := &http.Server{
		Addr:    vars.HTTPAddress + ":" + vars.HTTPPort,
		Handler: httphandler.NewHTTPHandler(globals, accountManager),
	}

	go func() {
		globals.Logger.Info("HTTP server initializing",
			zap.String("http_addr", vars.HTTPAddress),
			zap.String("http_port", vars.HTTPPort))
		if err := server.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				globals.Logger.Fatal("HTTP server failed", zap.Error(err))
			} else {
				globals.Logger.Info("HTTP server shutdown")
			}
		}
	}()

	return server
}

// waitShutdown waits for shutdown signal
func waitShutdown(interruptCh <-chan os.Signal,
	logger *zap.Logger, dbClient service.PostgresClient, server *http.Server) {

	// Wait for interrupt
	<-interruptCh

	dbClient.Stop()
	dbClient = nil

	// Graceful shutdown HTTP server
	ctx, cancel := context.WithTimeout(context.Background(), httpShutdownTimeout)
	defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		logger.Info("Failed to graceful shutdown server", zap.Error(err))
	}
	server = nil

	logger.Info("Stop service", zap.String("git_version", GitVersion))
}
