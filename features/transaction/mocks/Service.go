// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "bengcall/features/transaction/domain"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// All provides a mock function with given fields:
func (_m *Service) All() ([]domain.TransactionAll, error) {
	ret := _m.Called()

	var r0 []domain.TransactionAll
	if rf, ok := ret.Get(0).(func() []domain.TransactionAll); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TransactionAll)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Cancel provides a mock function with given fields: ID
func (_m *Service) Cancel(ID uint) error {
	ret := _m.Called(ID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Comment provides a mock function with given fields: updateCmmt, ID
func (_m *Service) Comment(updateCmmt domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	ret := _m.Called(updateCmmt, ID)

	var r0 domain.TransactionCore
	if rf, ok := ret.Get(0).(func(domain.TransactionCore, uint) domain.TransactionCore); ok {
		r0 = rf(updateCmmt, ID)
	} else {
		r0 = ret.Get(0).(domain.TransactionCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.TransactionCore, uint) error); ok {
		r1 = rf(updateCmmt, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Detail provides a mock function with given fields: ID
func (_m *Service) Detail(ID uint) (domain.TransactionDetail, error) {
	ret := _m.Called(ID)

	var r0 domain.TransactionDetail
	if rf, ok := ret.Get(0).(func(uint) domain.TransactionDetail); ok {
		r0 = rf(ID)
	} else {
		r0 = ret.Get(0).(domain.TransactionDetail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// History provides a mock function with given fields: userID
func (_m *Service) History(userID uint) ([]domain.TransactionHistory, error) {
	ret := _m.Called(userID)

	var r0 []domain.TransactionHistory
	if rf, ok := ret.Get(0).(func(uint) []domain.TransactionHistory); ok {
		r0 = rf(userID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.TransactionHistory)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// My provides a mock function with given fields: userID
func (_m *Service) My(userID uint) (domain.TransactionHistory, error) {
	ret := _m.Called(userID)

	var r0 domain.TransactionHistory
	if rf, ok := ret.Get(0).(func(uint) domain.TransactionHistory); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(domain.TransactionHistory)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Status provides a mock function with given fields: updateStts, ID
func (_m *Service) Status(updateStts domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	ret := _m.Called(updateStts, ID)

	var r0 domain.TransactionCore
	if rf, ok := ret.Get(0).(func(domain.TransactionCore, uint) domain.TransactionCore); ok {
		r0 = rf(updateStts, ID)
	} else {
		r0 = ret.Get(0).(domain.TransactionCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.TransactionCore, uint) error); ok {
		r1 = rf(updateStts, ID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Success provides a mock function with given fields: ID
func (_m *Service) Success(ID uint) error {
	ret := _m.Called(ID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(ID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// Transaction provides a mock function with given fields: newTrx, newDtl
func (_m *Service) Transaction(newTrx domain.TransactionCore, newDtl []domain.DetailCore) (domain.TransactionDetail, error) {
	ret := _m.Called(newTrx, newDtl)

	var r0 domain.TransactionDetail
	if rf, ok := ret.Get(0).(func(domain.TransactionCore, []domain.DetailCore) domain.TransactionDetail); ok {
		r0 = rf(newTrx, newDtl)
	} else {
		r0 = ret.Get(0).(domain.TransactionDetail)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.TransactionCore, []domain.DetailCore) error); ok {
		r1 = rf(newTrx, newDtl)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewService interface {
	mock.TestingT
	Cleanup(func())
}

// NewService creates a new instance of Service. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewService(t mockConstructorTestingTNewService) *Service {
	mock := &Service{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}