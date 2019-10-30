package validator

import (
	"errors"
	"fmt"
	"strings"
)

const (
	errorMessageFmt = "field %v should be %v, %v detected"

	usd = "USD"
)

// Validator defines validator object for input parameters (don't trust to anybody!)
type Validator struct {
	errors []string
}

// NewValidator creates new validator
func NewValidator() *Validator {
	return &Validator{
		errors: make([]string, 0),
	}
}

// AddField add error for a field with expected value
func (v *Validator) AddField(field string, actual string, expected string) {
	v.errors = append(v.errors, fmt.Sprintf(errorMessageFmt, field, expected, actual))
}

// Error returns current error status
func (v *Validator) Error() error {
	if len(v.errors) == 0 {
		return nil
	}
	return errors.New(strings.Join(v.errors, "\n"))
}

// ValidateUID validates user's UID
func (v *Validator) ValidateUID(field string, uid string) *Validator {
	if len(uid) == 0 {
		v.AddField(field, "nil", "string")
	}
	return v
}

// ValidateCurrency validates user's currency
func (v *Validator) ValidateCurrency(currency string) *Validator {
	if strings.ToUpper(currency) != usd {
		v.AddField("currency", currency, usd)
	}
	return v
}

// ValidateBalance validates user's balance
func (v *Validator) ValidateBalance(balance float64) *Validator {
	if balance < 0 {
		v.AddField("balance", fmt.Sprintf("%.2f", balance), "> 0")
	}
	return v
}

// ValidateBalance validates payment's amount
func (v *Validator) ValidateAmount(amount float64) *Validator {
	if amount < 0 {
		v.AddField("amount", fmt.Sprintf("%.2f", amount), "> 0")
	}
	return v
}
