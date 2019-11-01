package repositoryengine

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestRepositoryEngine(t *testing.T) {
	suite.Run(t, new(contextTestSuite))
	suite.Run(t, new(repositoryFactoryTestSuite))
	suite.Run(t, new(paymentRepositoryTestSuite))
	suite.Run(t, new(accountRepositoryTestSuite))
	suite.Run(t, new(utilsTestSuite))
	suite.Run(t, new(scopeTestSuite))
}
