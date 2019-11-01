package repositoryengine

import (
	"github.com/stretchr/testify/suite"
)

type repositoryFactoryTestSuite struct {
	suite.Suite
}

func (s *repositoryFactoryTestSuite) TestCreateScopeSucceeded() {
	factory := NewRepositoryFactory(nil)
	scope := factory.Scope()

	s.NotNil(scope)
}

func (s *repositoryFactoryTestSuite) TestGetAccountRepositorySucceeded() {
	factory := NewRepositoryFactory(nil)
	repository := factory.AccountRepository()

	s.NotNil(repository)
	s.Equal(factory.(*repositoryFactory).accountRepository, repository)
}

func (s *repositoryFactoryTestSuite) TestGetPaymentRepositorySucceeded() {
	factory := NewRepositoryFactory(nil)
	repository := factory.PaymentRepository()

	s.NotNil(repository)
	s.Equal(factory.(*repositoryFactory).paymentRepository, repository)
}
