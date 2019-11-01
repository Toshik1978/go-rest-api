// +build test

package testutil

import (
	"time"

	"github.com/AlekSi/pointer"
	"github.com/Toshik1978/go-rest-api/handler"
	"github.com/Toshik1978/go-rest-api/repository"
)

// Package contains helpers for testing purpose

func RepositoryAccount() repository.Account {
	return repository.Account{
		ID:        1234,
		UID:       "toshik1978",
		Currency:  "USD",
		Balance:   10000,
		CreatedAt: time.Now().Round(time.Millisecond),
	}
}

func RepositoryPayment() repository.Payment {
	return repository.Payment{
		ID:                  1234,
		Amount:              10000,
		PayerAccountUID:     "toshik1978",
		RecipientAccountUID: "toshik1979",
		CreatedAt:           time.Now().Round(time.Millisecond),
	}
}

func AccountRequest() handler.AccountRequest {
	return handler.AccountRequest{
		UID:      "toshik1978",
		Currency: "USD",
		Balance:  100,
	}
}

func AccountResponse() handler.Account {
	return handler.Account{
		UID:       "toshik1978",
		Currency:  "USD",
		Balance:   100,
		CreatedAt: time.Now().Round(time.Millisecond),
	}
}

func PaymentRequest() handler.PaymentRequest {
	return handler.PaymentRequest{
		RecipientUID: "toshik1979",
		Amount:       100,
	}
}

func PaymentResponse() handler.Payment {
	return handler.Payment{
		UID:       "toshik1978",
		TargetUID: pointer.ToString("toshik1979"),
		Direction: "outgoing",
		Amount:    100,
		CreatedAt: time.Now().Round(time.Millisecond),
	}
}

func EqualStrings(s1 *string, s2 *string) bool {
	if s1 == s2 {
		return true
	}
	if s1 != nil && s2 != nil && *s1 == *s2 {
		return true
	}
	return false
}
