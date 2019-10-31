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

func (s *validatorTestSuite) TestValidateCurrencyFailed() {
	v := NewValidator()
	s.Error(v.ValidateCurrency("rub").Error())
}

func (s *validatorTestSuite) TestValidateCurrencySucceeded() {
	v := NewValidator()
	s.NoError(v.ValidateCurrency("usd").Error())
}

func (s *validatorTestSuite) TestValidateBalanceFailed() {
	v := NewValidator()
	s.Error(v.ValidateBalance(-100).Error())
}

func (s *validatorTestSuite) TestValidateBalanceSucceeded() {
	v := NewValidator()
	s.NoError(v.ValidateBalance(100).Error())
}

func (s *validatorTestSuite) TestValidateAmountFailed() {
	v := NewValidator()
	s.Error(v.ValidateAmount(-100).Error())
}

func (s *validatorTestSuite) TestValidateAmountSucceeded() {
	v := NewValidator()
	s.NoError(v.ValidateAmount(100).Error())
}
