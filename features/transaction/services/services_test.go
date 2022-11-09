package services

import (
	"bengcall/features/transaction/domain"
	"bengcall/features/transaction/mocks"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestTransaction(t *testing.T) {
	repo := new(mocks.Repository)
	service := domain.TransactionDetail{ID: uint(1), Invoice: 237863, Total: 150000, PaymentToken: "9843j2jnasja12830sd9", PaymentLink: "https://sanbox.midtrans.com/kcldskchuhv2093840vdvsec", Status: 1}

	t.Run("Success", func(t *testing.T) {
		repo.On("Post", mock.Anything, mock.Anything).Return(service, nil).Once()

		srv := New(repo)
		input := domain.TransactionCore{Location: 1, Schedule: time.Now(), Phone: "081234567890", Address: "Jl. Pahlawan No. 32, Surabaya"}
		inputs := []domain.DetailCore{{VehicleID: 1, ServiceID: 1, SubTotal: 50000}, {VehicleID: 1, ServiceID: 2, SubTotal: 75000}, {VehicleID: 1, ServiceID: 3, SubTotal: 25000}}
		res, err := srv.Transaction(input, inputs)
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Post", mock.Anything, mock.Anything).Return(domain.TransactionDetail{}, errors.New("there's duplicate data")).Once()
		srv := New(repo)
		res, err := srv.Transaction(domain.TransactionCore{}, []domain.DetailCore{})
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, 0, res.Invoice)
		assert.Error(t, err)
		assert.EqualError(t, err, "rejected from database")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Post", mock.Anything, mock.Anything).Return(domain.TransactionDetail{}, errors.New("there is some error")).Once()
		srv := New(repo)
		res, err := srv.Transaction(domain.TransactionCore{}, []domain.DetailCore{})
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, 0, res.Invoice)
		assert.Error(t, err)
		assert.EqualError(t, err, "some problem on database")
		repo.AssertExpectations(t)
	})
}

func TestSuccess(t *testing.T) {
	repo := new(mocks.Repository)

	t.Run("Success", func(t *testing.T) {
		repo.On("PutScss", uint(1)).Return(nil).Once()
		srv := New(repo)
		err := srv.Success(uint(1))
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("PutScss", uint(2)).Return(errors.New("there's duplicate data")).Once()
		srv := New(repo)
		err := srv.Success(uint(2))
		assert.Error(t, err)
		assert.EqualError(t, err, "Rejected from Database")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("PutScss", uint(3)).Return(errors.New("there is some error")).Once()
		srv := New(repo)
		err := srv.Success(uint(3))
		assert.Error(t, err)
		assert.EqualError(t, err, "Some Problem on Database")
		repo.AssertExpectations(t)
	})
}

func TestStatus(t *testing.T) {
	repo := new(mocks.Repository)
	service := domain.TransactionCore{ID: uint(1), Status: 3}

	t.Run("Success", func(t *testing.T) {
		repo.On("PutStts", mock.Anything, uint(1)).Return(service, nil).Once()

		srv := New(repo)
		input := domain.TransactionCore{Status: 3}
		res, err := srv.Status(input, uint(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, service.ID, res.ID)
		assert.Equal(t, service.Status, res.Status)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("PutStts", mock.Anything, uint(2)).Return(domain.TransactionCore{}, errors.New("there's duplicate data")).Once()
		srv := New(repo)
		res, err := srv.Status(domain.TransactionCore{}, uint(2))
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.Comment)
		assert.Error(t, err)
		assert.EqualError(t, err, "Rejected from Database")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("PutStts", mock.Anything, uint(3)).Return(domain.TransactionCore{}, errors.New("there is some error")).Once()
		srv := New(repo)
		res, err := srv.Status(domain.TransactionCore{}, uint(3))
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.Comment)
		assert.Error(t, err)
		assert.EqualError(t, err, "Some Problem on Database")
		repo.AssertExpectations(t)
	})
}

func TestComment(t *testing.T) {
	repo := new(mocks.Repository)
	service := domain.TransactionCore{ID: uint(1), Comment: "Great and nice."}

	t.Run("Success", func(t *testing.T) {
		repo.On("PutCmmt", mock.Anything, uint(1)).Return(service, nil).Once()

		srv := New(repo)
		input := domain.TransactionCore{Comment: "Great and nice."}
		res, err := srv.Comment(input, uint(1))
		assert.NoError(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, service.ID, res.ID)
		assert.Equal(t, service.Comment, res.Comment)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("PutCmmt", mock.Anything, uint(2)).Return(domain.TransactionCore{}, errors.New("there's duplicate data")).Once()
		srv := New(repo)
		res, err := srv.Comment(domain.TransactionCore{}, uint(2))
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.Comment)
		assert.Error(t, err)
		assert.EqualError(t, err, "Rejected from Database")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("PutCmmt", mock.Anything, uint(3)).Return(domain.TransactionCore{}, errors.New("there is some error")).Once()
		srv := New(repo)
		res, err := srv.Comment(domain.TransactionCore{}, uint(3))
		assert.Equal(t, uint(0), res.ID)
		assert.Equal(t, "", res.Comment)
		assert.Error(t, err)
		assert.EqualError(t, err, "Some Problem on Database")
		repo.AssertExpectations(t)
	})
}

func TestAll(t *testing.T) {
	repo := new(mocks.Repository)
	trx := []domain.TransactionAll{{ID: uint(1), Schedule: time.Now(), Invoice: 983672, Total: 150000, Status: 1, Fullname: "Gerdo Trada Wijaya"}, {ID: uint(2), Schedule: time.Now(), Invoice: 487231, Total: 50000, Status: 1, Fullname: "Lukmanul"}, {ID: uint(3), Schedule: time.Now(), Invoice: 873241, Total: 75000, Status: 1, Fullname: "Jerry Young"}}

	t.Run("Success", func(t *testing.T) {
		repo.On("GetAll").Return(trx, nil).Once()
		srv := New(repo)
		_, err := srv.All()
		assert.NoError(t, err)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("GetAll").Return(nil, errors.New("table not exists")).Once()
		srv := New(repo)
		_, err := srv.All()
		assert.Error(t, err)
		assert.EqualError(t, err, "Database Error")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("GetAll").Return(nil, errors.New("data not found")).Once()
		srv := New(repo)
		_, err := srv.All()
		assert.Error(t, err)
		assert.EqualError(t, err, "No Data")
		repo.AssertExpectations(t)
	})
}

func TestHistory(t *testing.T) {
	repo := new(mocks.Repository)
	trx := []domain.TransactionHistory{{ID: uint(1), Schedule: time.Now(), Invoice: 983672, Total: 150000}, {ID: uint(2), Schedule: time.Now(), Invoice: 487231, Total: 50000}, {ID: uint(3), Schedule: time.Now(), Invoice: 873241, Total: 75000}}

	t.Run("Success", func(t *testing.T) {
		repo.On("GetHistory", uint(1)).Return(trx, nil).Once()
		srv := New(repo)
		_, err := srv.History(uint(1))
		assert.NoError(t, err)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("GetHistory", uint(2)).Return(nil, errors.New("table not exists")).Once()
		srv := New(repo)
		_, err := srv.History(uint(2))
		assert.Error(t, err)
		assert.EqualError(t, err, "Database Error")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("GetHistory", uint(3)).Return(nil, errors.New("data not found")).Once()
		srv := New(repo)
		_, err := srv.History(uint(3))
		assert.Error(t, err)
		assert.EqualError(t, err, "No Data")
		repo.AssertExpectations(t)
	})
}

func TestDetail(t *testing.T) {
	repo := new(mocks.Repository)
	trx := domain.TransactionDetail{ID: uint(1), Schedule: time.Now(), Invoice: 983672, Total: 150000}

	t.Run("Success", func(t *testing.T) {
		repo.On("GetDetail", uint(1)).Return(trx, nil).Once()
		srv := New(repo)
		_, err := srv.Detail(uint(1))
		assert.NoError(t, err)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("GetDetail", uint(2)).Return(domain.TransactionDetail{}, errors.New("table not exists")).Once()
		srv := New(repo)
		_, err := srv.Detail(uint(2))
		assert.Error(t, err)
		assert.EqualError(t, err, "Database Error")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("GetDetail", uint(3)).Return(domain.TransactionDetail{}, errors.New("data not found")).Once()
		srv := New(repo)
		_, err := srv.Detail(uint(3))
		assert.Error(t, err)
		assert.EqualError(t, err, "No Data")
		repo.AssertExpectations(t)
	})
}

func TestCancel(t *testing.T) {
	repo := new(mocks.Repository)

	t.Run("Success", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(nil).Once()

		srv := New(repo)

		err := srv.Cancel(1)
		assert.NoError(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New("data not found")).Once()

		srv := New(repo)

		err := srv.Cancel(2)
		assert.Error(t, err)
		assert.EqualError(t, err, "no data")
		repo.AssertExpectations(t)
	})

	t.Run("Failed", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New("table not exists")).Once()

		srv := New(repo)

		err := srv.Cancel(3)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		repo.AssertExpectations(t)
	})
}
