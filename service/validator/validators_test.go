package validator

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

func TestValidators(t *testing.T) {
	suite.Run(t, new(validatorTestSuite))
}
