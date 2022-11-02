package services

import (
	"bengcall/features/user/domain"
	"bengcall/features/user/mocks"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Add User", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "joko", Email: "joko@gmail.com",
			Password: "$2a$10$wNRtkCldluk.H6QXWPSz1.5KUcNNY0ncyh/LsPtp2IVagSD.ITRiK", Images: "https://bengcallbucket.s3.ap-southeast-1.amazonaws.com/profile/Q5aWl5c2RKoHcIFIrbMi-dummy450x450.jpg", Role: 0, Token: ""}, nil).Once()

		srv := New(repo)
		input := domain.UserCore{Password: "joko", Fullname: "joko", Email: "joko@gmail.com"}
		res, err := srv.Register(input)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Register", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{}, errors.New("some problem on database")).Once()
		srv := New(repo)
		res, err := srv.Register(domain.UserCore{})
		assert.Empty(t, res)
		assert.NotNil(t, err)
		repo.AssertExpectations(t)
	})

}

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Login", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{Password: "$2a$10$wNRtkCldluk.H6QXWPSz1.5KUcNNY0ncyh/LsPtp2IVagSD.ITRiK "}, nil).Once()
		srv := New(repo)
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.NotEmpty(t, res)
		assert.Nil(t, err)
		repo.AssertExpectations(t)
	})

	t.Run("Wrong Login", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{Password: "asgfasg"}, errors.New("no data")).Once()
		srv := New(repo)
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.EqualError(t, err, "password not match")
		repo.AssertExpectations(t)
	})

}

func TestMyProfile(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("success", func(t *testing.T) {
		repo.On("GetMyUser", mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "joko", Email: "joko@gmail.com", Images: "https://bengcallbucket.s3.ap-southeast-1.amazonaws.com/profile/Q5aWl5c2RKoHcIFIrbMi-dummy450x450.jpg", Role: 0}, nil).Once()
		srv := New(repo)
		res, err := srv.MyProfile(uint(1))
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Failed Get User", func(t *testing.T) {
		repo.On("GetMyUser", mock.Anything).Return(domain.UserCore{}, errors.New("no data")).Once()
		srv := New(repo)
		res, _ := srv.MyProfile(uint(1))
		//assert.NotNil(t, err)
		assert.Empty(t, res)
		repo.AssertExpectations(t)
	})

}

func TestUpdateProfile(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Update User", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "joko soleh", Password: "joko",
			Email: "joko@gmail.com"}, nil).Once()
		srv := New(repo)
		input := domain.UserCore{ID: uint(1), Fullname: "joko soleh", Password: "joko",
			Email: "joko@gmail.com"}
		res, err := srv.UpdateProfile(input, 1)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Update Profile", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("error update user")).Once()
		srv := New(repo)
		var input domain.UserCore
		res, err := srv.UpdateProfile(input, 1)
		assert.Empty(t, res)
		assert.NotNil(t, err)
		repo.AssertExpectations(t)
	})
}

func TestDeactivate(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Delete User", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(nil).Once()
		srv := New(repo)
		err := srv.Deactivate(1)
		assert.Nil(t, err)
		//assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

}
