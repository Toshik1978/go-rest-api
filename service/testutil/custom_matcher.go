// +build test

package testutil

import (
	"github.com/Toshik1978/go-rest-api/repository"
	"github.com/golang/mock/gomock"
)

// equalRepositoryAccountMatcher implements custom matcher for repository accounts
type equalRepositoryAccountMatcher struct {
	account repository.Account
}

// EqualRepositoryAccount return matcher instance
func EqualRepositoryAccount(account repository.Account) gomock.Matcher {
	return &equalRepositoryAccountMatcher{
		account: account,
	}
}

func (m *equalRepositoryAccountMatcher) Matches(x interface{}) bool {
	account, ok := x.(*repository.Account)
	if !ok {
		return false
	}
	return account.UID == m.account.UID &&
		account.Balance == m.account.Balance &&
		account.Currency == m.account.Currency
}

func (m *equalRepositoryAccountMatcher) String() string {
	return "is equal account"
}

// equalRepositoryPaymentMatcher implements custom matcher for repository payments
type equalRepositoryPaymentMatcher struct {
	payment repository.Payment
}

// EqualRepositoryPayment return matcher instance
func EqualRepositoryPayment(payment repository.Payment) gomock.Matcher {
	return &equalRepositoryPaymentMatcher{
		payment: payment,
	}
}

func (m *equalRepositoryPaymentMatcher) Matches(x interface{}) bool {
	payment, ok := x.(*repository.Payment)
	if !ok {
		return false
	}
	return payment.PayerAccountUID == m.payment.PayerAccountUID &&
		payment.RecipientAccountUID == m.payment.RecipientAccountUID &&
		payment.Amount == m.payment.Amount
}

func (m *equalRepositoryPaymentMatcher) String() string {
	return "is equal payment"
}
