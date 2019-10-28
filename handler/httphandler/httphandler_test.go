package httphandler

import (
	"net/http"
	"net/http/httptest"

	"github.com/Toshik1978/go-rest-api/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

// panicResponseWriter declare ResponseWriter, which always panic
type panicResponseWriter struct {
	headers http.Header
	status  int
}

func (r *panicResponseWriter) Header() http.Header {
	return r.headers
}

func (r *panicResponseWriter) Write(body []byte) (int, error) {
	panic("panic")
}

func (r *panicResponseWriter) WriteHeader(status int) {
	r.status = status
}

func newPanicResponseWriter() *panicResponseWriter {
	return &panicResponseWriter{
		headers: make(http.Header),
	}
}

// HTTPHandler test suite
type httpHandlerTestSuite struct {
	suite.Suite
}

func (s *httpHandlerTestSuite) TestHTTPHandlerSucceeded1() {
	req, err := http.NewRequest("GET", "/v1/server/status", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	handler := NewHTTPHandler(service.Globals{
		Logger: zap.New(zapCore),
	})

	r := httptest.NewRecorder()
	handler.ServeHTTP(r, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Handled HTTP request", zapRecorded.All()[0].Message)
	s.Equal(http.StatusOK, r.Code)
}

func (s *httpHandlerTestSuite) TestHTTPHandlerSucceeded2() {
	req, err := http.NewRequest("GET", "/debug/pprof/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	handler := NewHTTPHandler(service.Globals{
		Logger: zap.New(zapCore),
	})

	r := httptest.NewRecorder()
	handler.ServeHTTP(r, req)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusOK, r.Code)
}

func (s *httpHandlerTestSuite) TestHTTPHandlerSucceeded3() {
	req, err := http.NewRequest("GET", "/test", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	handler := NewHTTPHandler(service.Globals{
		Logger: zap.New(zapCore),
	})

	r := httptest.NewRecorder()
	handler.ServeHTTP(r, req)

	s.Equal(0, zapRecorded.Len())
	s.Equal(http.StatusNotFound, r.Code)
}

func (s *httpHandlerTestSuite) TestHTTPHandlerSucceeded4() {
	req, err := http.NewRequest("GET", "/debug/pprof/", nil)
	if err != nil {
		s.T().Fatal(err)
	}

	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)
	handler := NewHTTPHandler(service.Globals{
		Logger: zap.New(zapCore),
	})

	w := newPanicResponseWriter()
	handler.ServeHTTP(w, req)

	s.Equal(1, zapRecorded.Len())
	s.Equal("Panic happened in HTTP handler", zapRecorded.All()[0].Message)
	s.Equal(http.StatusInternalServerError, w.status)
}
