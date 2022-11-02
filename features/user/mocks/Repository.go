// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "bengcall/features/user/domain"

	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// AddUser provides a mock function with given fields: newUser
func (_m *Repository) AddUser(newUser domain.UserCore) (domain.UserCore, error) {
	ret := _m.Called(newUser)

	var r0 domain.UserCore
	if rf, ok := ret.Get(0).(func(domain.UserCore) domain.UserCore); ok {
		r0 = rf(newUser)
	} else {
		r0 = ret.Get(0).(domain.UserCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.UserCore) error); ok {
		r1 = rf(newUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Delete provides a mock function with given fields: userID
func (_m *Repository) Delete(userID uint) error {
	ret := _m.Called(userID)

	var r0 error
	if rf, ok := ret.Get(0).(func(uint) error); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetMyUser provides a mock function with given fields: userID
func (_m *Repository) GetMyUser(userID uint) (domain.UserCore, error) {
	ret := _m.Called(userID)

	var r0 domain.UserCore
	if rf, ok := ret.Get(0).(func(uint) domain.UserCore); ok {
		r0 = rf(userID)
	} else {
		r0 = ret.Get(0).(domain.UserCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetUser provides a mock function with given fields: existUser
func (_m *Repository) GetUser(existUser domain.UserCore) (domain.UserCore, error) {
	ret := _m.Called(existUser)

	var r0 domain.UserCore
	if rf, ok := ret.Get(0).(func(domain.UserCore) domain.UserCore); ok {
		r0 = rf(existUser)
	} else {
		r0 = ret.Get(0).(domain.UserCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.UserCore) error); ok {
		r1 = rf(existUser)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// Update provides a mock function with given fields: updatedUser, userID
func (_m *Repository) Update(updatedUser domain.UserCore, userID uint) (domain.UserCore, error) {
	ret := _m.Called(updatedUser, userID)

	var r0 domain.UserCore
	if rf, ok := ret.Get(0).(func(domain.UserCore, uint) domain.UserCore); ok {
		r0 = rf(updatedUser, userID)
	} else {
		r0 = ret.Get(0).(domain.UserCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.UserCore, uint) error); ok {
		r1 = rf(updatedUser, userID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

type mockConstructorTestingTNewRepository interface {
	mock.TestingT
	Cleanup(func())
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
func NewRepository(t mockConstructorTestingTNewRepository) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
