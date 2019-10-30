package httphandler

import (
	"github.com/Toshik1978/go-rest-api/service/server"
	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

// recoveryLogger implements logger for panic recovery middleware
type recoveryLogger struct {
	logger *zap.Logger
}

// newRecoveryLogger creates new RecoveryHandlerLogger instance
func newRecoveryLogger(globals server.Globals) handlers.RecoveryHandlerLogger {
	return &recoveryLogger{
		logger: globals.Logger,
	}
}

// Println log panic
func (l *recoveryLogger) Println(v ...interface{}) {
	l.logger.With(
		zap.String("mode", "panic_log"),
		zap.Any("panic", v),
		zap.Stack("stack"),
	).Error("Panic happened in HTTP handler")
}
