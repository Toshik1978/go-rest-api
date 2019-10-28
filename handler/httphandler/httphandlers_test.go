package httphandler

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestHTTPHandler(t *testing.T) {
	suite.Run(t, new(httpHandlerTestSuite))
	suite.Run(t, new(logFormatterTestSuite))
	suite.Run(t, new(recoveryLoggerTestSuite))
	suite.Run(t, new(apiHandlerTestSuite))
}
