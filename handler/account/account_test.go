package account

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestAccount(t *testing.T) {
	suite.Run(t, new(accountManagerTestSuite))
	suite.Run(t, new(accountBuilderTestSuite))
	suite.Run(t, new(paymentBuilderTestSuite))
}
