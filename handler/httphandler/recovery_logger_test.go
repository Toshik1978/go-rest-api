package httphandler

import (
	"github.com/Toshik1978/go-rest-api/service"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

type recoveryLoggerTestSuite struct {
	suite.Suite
}

func (s *recoveryLoggerTestSuite) TestRecoveryLoggerSucceeded() {
	zapCore, zapRecorded := observer.New(zapcore.InfoLevel)

	recovery := newRecoveryLogger(service.Globals{
		Logger: zap.New(zapCore),
	})
	recovery.Println("panic")

	s.Equal(1, zapRecorded.Len())
	s.Equal("Panic happened in HTTP handler", zapRecorded.All()[0].Message)
}
