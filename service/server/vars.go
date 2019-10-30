package server

import (
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

const (
	configFileName = "go-rest-api.conf"
)

// Vars declare variables for service running
type Vars struct {
	HTTPAddress string
	HTTPPort    string

	DB        string
	DBTimeout time.Duration
}

// LoadConfig load config
func LoadConfig(logger *zap.Logger) Vars {
	viper.SetConfigName(configFileName)
	viper.AddConfigPath("configs")
	viper.AddConfigPath("/etc/go-rest-api")
	if err := viper.ReadInConfig(); err != nil {
		logger.Fatal("Failed to read config", zap.Error(err))
	}

	return Vars{
		HTTPAddress: viper.GetString("http.host"),
		HTTPPort:    viper.GetString("http.port"),
		DB:          viper.GetString("db.master"),
		DBTimeout:   viper.GetDuration("db.timeout"),
	}
}
