package httphandler

import (
	"net/http"
	"net/http/httptest"
	"time"

	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/service"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// failResponseWriter declare ResponseWriter, which always fails
type failResponseWriter struct {
	headers http.Header
}

func (r *failResponseWriter) Header() http.Header {
	return r.headers
}

func (r *failResponseWriter) Write(body []byte) (int, error) {
	return 0, errors.New("fail")
}

func (r *failResponseWriter) WriteHeader(status int) {
}

func newFailResponseWriter() http.ResponseWriter {
	return &failResponseWriter{
		headers: make(http.Header),
	}
}

// apiHandlerTestSuite test suite
type apiHandlerTestSuite struct {
	suite.Suite
}

func (s *apiHandlerTestSuite) TestWriteResponseFailed1() {
	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)

	apiHandler := newAPIHandler(service.Globals{
		Logger: zap.New(zapCore),
	})
	apiHandler.writeResponse(newFailResponseWriter(), make(chan struct{}))

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to marshal HTTP response", zapRecorded.All()[0].Message)
}

func (s *apiHandlerTestSuite) TestWriteResponseFailed2() {
	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)

	apiHandler := newAPIHandler(service.Globals{
		Logger: zap.New(zapCore),
	})
	apiHandler.writeResponse(newFailResponseWriter(), handler.ServerStatusResponse{})

	s.Equal(1, zapRecorded.Len())
	s.Equal("Failed to write HTTP response", zapRecorded.All()[0].Message)
}

func (s *apiHandlerTestSuite) TestServerStatusHandlerSucceeded() {
	ctrl := gomock.NewController(s.T())
	defer ctrl.Finish()

	req, err := http.NewRequest("GET", "/v1/server/status", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	apiHandler := newAPIHandler(service.Globals{
		Logger:    zap.New(zapCore),
		BuildTime: time.Now().String(),
		Version:   "test",
	}).ServerStatusHandler()

	r := httptest.NewRecorder()
	apiHandler.ServeHTTP(r, req)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusOK, r.Code)
}
