package services

import (
	"bengcall/features/user/domain"
	"bengcall/features/user/mocks"
	"errors"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Add User", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "Khafid", Email: "khafid@gmail.com",
			Password: "$2a$10$wNRtkCldluk.H6QXWPSz1.5KUcNNY0ncyh/LsPtp2IVagSD.ITRiK", Images: "https://bengcallbucket.s3.ap-southeast-1.amazonaws.com/profile/Q5aWl5c2RKoHcIFIrbMi-dummy450x450.jpg", Role: 0, Token: ""}, nil).Once()

		srv := New(repo, validator.New())
		input := domain.UserCore{Fullname: "Khafid", Email: "khafid@gmail.com", Password: "khafid123456789"}
		res, err := srv.Register(input)
		assert.Nil(t, err)
		assert.Equal(t, "Khafid", res.Fullname)
		assert.Equal(t, "khafid@gmail.com", res.Email)
		assert.Equal(t, "$2a$10$wNRtkCldluk.H6QXWPSz1.5KUcNNY0ncyh/LsPtp2IVagSD.ITRiK", res.Password, "Password tidak sesuai")
		//assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Register", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{}, errors.New("email can't be use")).Once()
		srv := New(repo, validator.New())
		res, err := srv.Register(domain.UserCore{})
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid email")

	})

	t.Run("Validator error", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{}, errors.New("validation error")).Once()

		input := domain.UserCore{Password: "jokoagung", Fullname: "Joko", Email: "joko@gmail.com"}
		srv := New(repo, validator.New())
		res, err := srv.Register(input)
		assert.NotNil(t, err, "validation error") // Apakah errornya nil
		//assert.Greater(t, res.ID, 0)              // Apakah ID nya lebih besar dari 0
		assert.Equal(t, "", res.Fullname) // Apakah nama yang di insertkan sama
		assert.Equal(t, "", res.Email)
		assert.Equal(t, "", res.Password, "Password tidak sesuai")
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Register", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{}, errors.New("duplicate data")).Once()
		srv := New(repo, validator.New())
		res, err := srv.Register(domain.UserCore{})
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "rejected from database")
		repo.AssertExpectations(t)
	})
}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Login", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{Password: "$2a$10$wNRtkCldluk.H6QXWPSz1.5KUcNNY0ncyh/LsPtp2IVagSD.ITRiK "}, nil).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.NotEmpty(t, res)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Wrong Login", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{Password: "asgfasg"}, errors.New("no data")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.EqualError(t, err, "email not found")
		repo.AssertExpectations(t)
	})

	t.Run("Fail Database Error", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{Password: "asgfasg"}, errors.New("table not exists")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed. Email or Password not found.")
		repo.AssertExpectations(t)
	})

	t.Run("Fail No Data", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{Password: "asgfasg"}, errors.New("data not found")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed. Email or Password not found.")
		repo.AssertExpectations(t)
	})
}

func TestMyProfile(t *testing.T) {
	repo := mocks.NewRepository(t)

	t.Run("success", func(t *testing.T) {
		repo.On("GetMyUser", mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "joko", Email: "joko@gmail.com", Images: "https://bengcallbucket.s3.ap-southeast-1.amazonaws.com/profile/Q5aWl5c2RKoHcIFIrbMi-dummy450x450.jpg", Role: 0}, nil).Once()
		srv := New(repo, validator.New())
		res, err := srv.MyProfile(uint(1))
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Get User", func(t *testing.T) {
		repo.On("GetMyUser", mock.Anything).Return(domain.UserCore{}, errors.New("data not found")).Once()
		srv := New(repo, validator.New())
		res, err := srv.MyProfile(uint(2))
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "no data")
		repo.AssertExpectations(t)
	})

	t.Run("Failed Get User", func(t *testing.T) {
		repo.On("GetMyUser", mock.Anything).Return(domain.UserCore{}, errors.New("table not exists")).Once()
		srv := New(repo, validator.New())
		res, err := srv.MyProfile(uint(3))
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		repo.AssertExpectations(t)
	})
}

func TestUpdateProfile(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Update User", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "joko soleh", Password: "joko",
			Email: "joko@gmail.com"}, nil).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{ID: uint(1), Fullname: "joko soleh", Password: "joko",
			Email: "joko@gmail.com"}
		res, err := srv.UpdateProfile(input, 1)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Update Profile", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("error on bcrypt password updated user")).Once()
		srv := New(repo, validator.New())
		var input domain.UserCore
		res, err := srv.UpdateProfile(input, 1)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "some problem on database")
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Update Profile", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("column doesnt exists")).Once()
		srv := New(repo, validator.New())
		var input domain.UserCore
		res, err := srv.UpdateProfile(input, 1)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "rejected from database")
		repo.AssertExpectations(t)
	})
}

func TestDeactivate(t *testing.T) {
	repo := mocks.NewRepository(t)

	t.Run("Sukses Delete User", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(nil).Once()
		srv := New(repo, validator.New())
		err := srv.Deactivate(1)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Fail Database Error", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New("table not exists")).Once()
		srv := New(repo, validator.New())
		err := srv.Deactivate(2)
		assert.Error(t, err)
		assert.EqualError(t, err, "database error")
		repo.AssertExpectations(t)
	})

	t.Run("Fail No Data", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(errors.New("data not found")).Once()
		srv := New(repo, validator.New())
		err := srv.Deactivate(3)
		assert.Error(t, err)
		assert.EqualError(t, err, "no data")
		repo.AssertExpectations(t)
	})
}
