// Code generated by MockGen. DO NOT EDIT.
// Source: repository.go

// Package mock is a generated GoMock package.
package mock

import (
	repository "github.com/Toshik1978/go-rest-api/repository"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockAccountRepository is a mock of AccountRepository interface
type MockAccountRepository struct {
	ctrl     *gomock.Controller
	recorder *MockAccountRepositoryMockRecorder
}

// MockAccountRepositoryMockRecorder is the mock recorder for MockAccountRepository
type MockAccountRepositoryMockRecorder struct {
	mock *MockAccountRepository
}

// NewMockAccountRepository creates a new mock instance
func NewMockAccountRepository(ctrl *gomock.Controller) *MockAccountRepository {
	mock := &MockAccountRepository{ctrl: ctrl}
	mock.recorder = &MockAccountRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAccountRepository) EXPECT() *MockAccountRepositoryMockRecorder {
	return m.recorder
}

// MockPaymentRepository is a mock of PaymentRepository interface
type MockPaymentRepository struct {
	ctrl     *gomock.Controller
	recorder *MockPaymentRepositoryMockRecorder
}

// MockPaymentRepositoryMockRecorder is the mock recorder for MockPaymentRepository
type MockPaymentRepositoryMockRecorder struct {
	mock *MockPaymentRepository
}

// NewMockPaymentRepository creates a new mock instance
func NewMockPaymentRepository(ctrl *gomock.Controller) *MockPaymentRepository {
	mock := &MockPaymentRepository{ctrl: ctrl}
	mock.recorder = &MockPaymentRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockPaymentRepository) EXPECT() *MockPaymentRepositoryMockRecorder {
	return m.recorder
}

// MockFactory is a mock of Factory interface
type MockFactory struct {
	ctrl     *gomock.Controller
	recorder *MockFactoryMockRecorder
}

// MockFactoryMockRecorder is the mock recorder for MockFactory
type MockFactoryMockRecorder struct {
	mock *MockFactory
}

// NewMockFactory creates a new mock instance
func NewMockFactory(ctrl *gomock.Controller) *MockFactory {
	mock := &MockFactory{ctrl: ctrl}
	mock.recorder = &MockFactoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockFactory) EXPECT() *MockFactoryMockRecorder {
	return m.recorder
}

// GetAccountRepository mocks base method
func (m *MockFactory) GetAccountRepository() repository.AccountRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAccountRepository")
	ret0, _ := ret[0].(repository.AccountRepository)
	return ret0
}

// GetAccountRepository indicates an expected call of GetAccountRepository
func (mr *MockFactoryMockRecorder) GetAccountRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAccountRepository", reflect.TypeOf((*MockFactory)(nil).GetAccountRepository))
}

// GetPaymentRepository mocks base method
func (m *MockFactory) GetPaymentRepository() repository.PaymentRepository {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPaymentRepository")
	ret0, _ := ret[0].(repository.PaymentRepository)
	return ret0
}

// GetPaymentRepository indicates an expected call of GetPaymentRepository
func (mr *MockFactoryMockRecorder) GetPaymentRepository() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPaymentRepository", reflect.TypeOf((*MockFactory)(nil).GetPaymentRepository))
}
