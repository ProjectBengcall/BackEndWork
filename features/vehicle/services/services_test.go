package services

import (
	"bengcall/config"
	"bengcall/features/vehicle/domain"
	"bengcall/features/vehicle/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func TestAddVehicle(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Success Add ", func(t *testing.T) {
		repo.On("Add", mock.Anything).Return(domain.VehicleCore{ID: uint(1), Name_vehicle: "Supra 125 125cc"}, nil).Once()
		srv := New(repo)
		input := domain.VehicleCore{Name_vehicle: "Supra 125 125cc"}
		res, err := srv.AddVehicle(input)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.NotEmpty(t, res.ID, "harusnya ada id yang terbuat")
		assert.Equal(t, input.Name_vehicle, res.Name_vehicle, "seharusnya nama sama")
		repo.AssertExpectations(t)
	})

	t.Run("Duplicate data", func(t *testing.T) {
		repo.On("Add", mock.Anything).Return(domain.VehicleCore{}, errors.New("there's duplicate data")).Once()
		srv := New(repo)
		input := domain.VehicleCore{Name_vehicle: "Supra 125 125cc"}
		res, err := srv.AddVehicle(input)
		assert.NotNil(t, err)
		assert.Empty(t, res, "karena object gagal dibuat")
		assert.Equal(t, uint(0), res.ID, "id harusnya 0 karena tidak ada data")
		assert.EqualError(t, err, "rejected from database")
		repo.AssertExpectations(t)
	})

	t.Run("Problem", func(t *testing.T) {
		repo.On("Add", mock.Anything).Return(domain.VehicleCore{}, errors.New("cannot connect")).Once()
		srv := New(repo)
		input := domain.VehicleCore{Name_vehicle: "Supra 125 125cc"}
		res, err := srv.AddVehicle(input)
		assert.NotNil(t, err)
		assert.Empty(t, res, "karena object gagal dibuat")
		assert.Equal(t, uint(0), res.ID, "id harusnya 0 karena tidak ada data")
		assert.EqualError(t, err, "some problem on database")
		repo.AssertExpectations(t)
	})
}

func TestGetVehicle(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Succses Get Vehicle", func(t *testing.T) {
		repo.On("GetAll", mock.Anything).Return([]domain.VehicleCore{{ID: uint(1), Name_vehicle: "Supra 125 125cc"}}, nil).Once()
		srv := New(repo)
		res, err := srv.GetVehicle()
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.Greater(t, res[0].ID, uint(0)) //lebih besar
		assert.GreaterOrEqual(t, len(res), 1) //lebih besar atau sama
		repo.AssertExpectations(t)
	})

	t.Run("Cant Retrive on database", func(t *testing.T) {
		repo.On("GetAll", mock.Anything).Return(nil, errors.New(config.DATABASE_ERROR)).Once()
		srv := New(repo)
		res, err := srv.GetVehicle()
		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, config.DATABASE_ERROR, "pesan error tidak sesuai")
		assert.Equal(t, len(res), 0, "len harusnya 0 karena tidak ada data")
		repo.AssertExpectations(t)
	})

	t.Run("Data not found", func(t *testing.T) {
		repo.On("GetAll", mock.Anything).Return(nil, gorm.ErrRecordNotFound).Once()
		srv := New(repo)
		res, err := srv.GetVehicle()
		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error(), "pesan error tidak sesuai")
		assert.Equal(t, len(res), 0, "len harusnya 0 karena tidak ada data")
		repo.AssertExpectations(t)
	})
}

func TestGetService(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Success Get Vehicle", func(t *testing.T) {
		repo.On("Get", mock.Anything).Return([]domain.VehicleCore{{ID: uint(1), Name_vehicle: "Supra"}}, []domain.ServiceVehicle{{ID: uint(1), ServiceName: "Full Service", Price: 100000}}, nil).Once()
		srv := New(repo)
		res, ser, err := srv.GetService()
		assert.Nil(t, err)
		assert.NotNil(t, res)
		assert.NotNil(t, ser)
		assert.Greater(t, res[0].ID, uint(0)) //lebih besar
		assert.GreaterOrEqual(t, len(res), 1) //lebih besar atau sama
		assert.Greater(t, ser[0].ID, uint(0)) //lebih besar
		assert.GreaterOrEqual(t, len(ser), 1) //lebih besar atau sama
		repo.AssertExpectations(t)
	})

	t.Run("Cant Retrive on database", func(t *testing.T) {
		repo.On("Get", mock.Anything).Return(nil, nil, errors.New(config.DATABASE_ERROR)).Once()
		srv := New(repo)
		res, ser, err := srv.GetService()
		assert.NotNil(t, err)
		assert.Nil(t, res)
		assert.Nil(t, ser)
		assert.EqualError(t, err, "get service error", "pesan error tidak sesuai")
		assert.Equal(t, len(res), 0, "len harusnya 0 karena tidak ada data")
		repo.AssertExpectations(t)
	})
}

func TestDeleteVehicle(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sucses delete profil", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(nil).Once()
		srv := New(repo)
		var id = uint(1)
		err := srv.DeleteVehicle(id)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Data not found", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(gorm.ErrRecordNotFound).Once()
		srv := New(repo)
		var id = uint(1)
		err := srv.DeleteVehicle(id)
		assert.NotNil(t, err)
		assert.EqualError(t, err, gorm.ErrRecordNotFound.Error(), "pesan error tidak sesuai")
		repo.AssertExpectations(t)
	})

	t.Run("error data on database", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New(config.DATABASE_ERROR)).Once()
		srv := New(repo)
		var id = uint(1)
		err := srv.DeleteVehicle(id)
		assert.NotNil(t, err)
		assert.EqualError(t, err, config.DATABASE_ERROR, "pesan error tidak sesuai")
		repo.AssertExpectations(t)
	})

}
