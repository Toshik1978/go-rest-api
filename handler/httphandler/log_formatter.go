package httphandler

import (
	"io"
	"time"

	"github.com/Toshik1978/go-rest-api/service"
	"github.com/gorilla/handlers"
	"go.uber.org/zap"
)

// logFormatter implements formatter for access logger middleware
type logFormatter struct {
	logger *zap.Logger
}

// newLogFormatter creates new access logger instance and return formatter function
func newLogFormatter(globals service.Globals) handlers.LogFormatter {
	return logFormatter{
		logger: globals.Logger,
	}.format
}

// format actually write access log (no need writer due to zap using)
func (f logFormatter) format(_ io.Writer, params handlers.LogFormatterParams) {
	f.logger.With(
		zap.Duration("duration", time.Since(params.TimeStamp)),
		zap.String("url", params.URL.String()),
		zap.String("method", params.Request.Method),
		zap.String("remote_addr", params.Request.RemoteAddr),
		zap.String("user_agent", params.Request.UserAgent()),
		zap.String("mode", "access_log"),
	).Info("Handled HTTP request")
}
