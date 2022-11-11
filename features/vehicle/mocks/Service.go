// Code generated by mockery v2.14.0. DO NOT EDIT.

package mocks

import (
	domain "bengcall/features/vehicle/domain"

	mock "github.com/stretchr/testify/mock"
)

// Service is an autogenerated mock type for the Service type
type Service struct {
	mock.Mock
}

// AddVehicle provides a mock function with given fields: newItem
func (_m *Service) AddVehicle(newItem domain.VehicleCore) (domain.VehicleCore, error) {
	ret := _m.Called(newItem)

	var r0 domain.VehicleCore
	if rf, ok := ret.Get(0).(func(domain.VehicleCore) domain.VehicleCore); ok {
		r0 = rf(newItem)
	} else {
		r0 = ret.Get(0).(domain.VehicleCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(domain.VehicleCore) error); ok {
		r1 = rf(newItem)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteVehicle provides a mock function with given fields: vehicleID
func (_m *Service) DeleteVehicle(vehicleID uint) (domain.VehicleCore, error) {
	ret := _m.Called(vehicleID)

	var r0 domain.VehicleCore
	if rf, ok := ret.Get(0).(func(uint) domain.VehicleCore); ok {
		r0 = rf(vehicleID)
	} else {
		r0 = ret.Get(0).(domain.VehicleCore)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(uint) error); ok {
		r1 = rf(vehicleID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetService provides a mock function with given fields:
func (_m *Service) GetService() ([]domain.VehicleCore, []domain.ServiceVehicle, error) {
	ret := _m.Called()

	var r0 []domain.VehicleCore
	if rf, ok := ret.Get(0).(func() []domain.VehicleCore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.VehicleCore)
		}
	}

	var r1 []domain.ServiceVehicle
	if rf, ok := ret.Get(1).(func() []domain.ServiceVehicle); ok {
		r1 = rf()
	} else {
		if ret.Get(1) != nil {
			r1 = ret.Get(1).([]domain.ServiceVehicle)
		}
	}

	var r2 error
	if rf, ok := ret.Get(2).(func() error); ok {
		r2 = rf()
	} else {
		r2 = ret.Error(2)
	}

	return r0, r1, r2
}

// GetVehicle provides a mock function with given fields:
func (_m *Service) GetVehicle() ([]domain.VehicleCore, error) {
	ret := _m.Called()

	var r0 []domain.VehicleCore
	if rf, ok := ret.Get(0).(func() []domain.VehicleCore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]domain.VehicleCore)
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
