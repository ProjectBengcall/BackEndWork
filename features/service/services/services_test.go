package services

import (
	"bengcall/features/service/domain"
	"bengcall/features/service/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetSpesific(t *testing.T) {
	repo := new(mocks.Repository)
	service := []domain.Core{{ID: uint(1), ServiceName: "Full Service", Price: 150000, VehicleID: 1}, {ID: uint(2), ServiceName: "Half Service", Price: 75000, VehicleID: 1}, {ID: uint(3), ServiceName: "Bit Service", Price: 37500, VehicleID: 1}}

	t.Run("Success", func(t *testing.T) {
		repo.On("Get", 1).Return(service, nil).Once()
		srv := New(repo)
		_, err := srv.GetSpesific(1)
		assert.NoError(t, err)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Get", 2).Return(nil, errors.New("table not exists")).Once()
		srv := New(repo)
		_, err := srv.GetSpesific(2)
		assert.Error(t, err)
		assert.EqualError(t, err, "Database Error")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Get", 2).Return(nil, errors.New("data not found")).Once()
		srv := New(repo)
		_, err := srv.GetSpesific(2)
		assert.Error(t, err)
		assert.EqualError(t, err, "No Data")
		repo.AssertExpectations(t)
	})
}

func TestAddService(t *testing.T) {
	repo := new(mocks.Repository)
	service := domain.Core{ID: uint(1), ServiceName: "Full Service", Price: 150000, VehicleID: 1}

	t.Run("Success", func(t *testing.T) {
		repo.On("Add", mock.Anything).Return(service, nil).Once()

		srv := New(repo)
		input := domain.Core{ServiceName: "Full Service", Price: 150000, VehicleID: 1}
		res, err := srv.AddService(input)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, service.ID, res.ID)
		assert.Equal(t, service.ServiceName, res.ServiceName)
		assert.Equal(t, service.Price, res.Price)
		assert.Equal(t, service.VehicleID, res.VehicleID)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Add", mock.Anything).Return(domain.Core{}, errors.New("there's duplicate data")).Once()
		srv := New(repo)
		res, err := srv.AddService(domain.Core{})
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.ServiceName)
		assert.Error(t, err)
		assert.EqualError(t, err, "rejected from database")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Add", mock.Anything).Return(domain.Core{}, errors.New("there is some error")).Once()
		srv := New(repo)
		res, err := srv.AddService(domain.Core{})
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.ServiceName)
		assert.Error(t, err)
		assert.EqualError(t, err, "some problem on database")
		repo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	repo := new(mocks.Repository)

	t.Run("Success", func(t *testing.T) {
		repo.On("Del", mock.Anything).Return(nil).Once()

		srv := New(repo)

		err := srv.Delete(1)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Del", mock.Anything).Return(errors.New("data not found")).Once()

		srv := New(repo)

		err := srv.Delete(2)
		assert.Error(t, err)
		assert.EqualError(t, err, "no data")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Del", mock.Anything).Return(errors.New("table not exists")).Once()

		srv := New(repo)

		err := srv.Delete(3)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		repo.AssertExpectations(t)
	})
}
