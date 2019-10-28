package httphandler

import (
	"net/http"
	"net/url"
	"time"

	"github.com/Toshik1978/go-rest-api/service"
	"github.com/gorilla/handlers"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type logFormatterTestSuite struct {
	suite.Suite
}

func (s *logFormatterTestSuite) TestLogFormatterSucceeded() {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)

	lf := newLogFormatter(service.Globals{
		Logger: zap.New(zapCore),
	})

	lf(nil, handlers.LogFormatterParams{
		Request:    req,
		URL:        url.URL{},
		TimeStamp:  time.Now(),
		StatusCode: http.StatusOK,
		Size:       1024,
	})

	s.Equal(1, zapRecorded.Len())
	s.Equal("Handled HTTP request", zapRecorded.All()[0].Message)
}
