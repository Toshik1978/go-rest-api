package service

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"go.uber.org/zap"
)

// Globals define some globally initialized in main objects
type Globals struct {
	Logger            *zap.Logger
	RepositoryFactory repository.Factory

	BuildTime string
	Version   string
}
