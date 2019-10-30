package validator

import (
	"github.com/stretchr/testify/suite"
)

type validatorTestSuite struct {
	suite.Suite
}

func (s *validatorTestSuite) TestValidateUIDFailed() {
	v := NewValidator()
	s.Error(v.ValidateUID("", "").Error())
}

func (s *validatorTestSuite) TestValidateUIDSucceeded() {
	v := NewValidator()
	s.NoError(v.ValidateUID("", "uid").Error())
}
