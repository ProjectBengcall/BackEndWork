package services

import (
	"bengcall/config"
	"bengcall/features/user/domain"
	"bengcall/features/user/mocks"
	rep "bengcall/features/user/repository"
	"errors"
	"mime/multipart"
	"testing"

	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func TestLogin(t *testing.T) {
	repo := mocks.NewRepository(t)
	t.Run("Sukses Login", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{
			ID:       uint(1),
			Fullname: "teamfour",
			Email:    "joko@gmail.com",
			Password: "$2a$10$wNRtkCldluk.H6QXWPSz1.5KUcNNY0ncyh/LsPtp2IVagSD.ITRiK"}, nil).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Nil(t, err)
		assert.NotNil(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Wrong Login", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{}, errors.New("no data")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "email not found")
		repo.AssertExpectations(t)
	})

	t.Run("Fail Database Error", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{}, errors.New("table not exists")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed. Email or Password not found.")
		repo.AssertExpectations(t)
	})

	t.Run("Fail No Data", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{}, errors.New("data not found")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: "joko"}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Failed. Email or Password not found.")
		repo.AssertExpectations(t)
	})

	t.Run("Wrong Password", func(t *testing.T) {
		repo.On("GetUser", mock.Anything).Return(domain.UserCore{}, errors.New("wrong password")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Email: "joko@gmail.com", Password: ""}
		res, err := srv.Login(input)
		assert.Empty(t, res)
		assert.NotNil(t, err)
		assert.Equal(t, input.Password, res.Password)
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

func TestDeactivate(t *testing.T) {
	repo := mocks.NewRepository(t)

	t.Run("Sukses Delete User", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(domain.UserCore{ID: uint(1)}, nil).Once()
		srv := New(repo, validator.New())
		res, err := srv.Deactivate(1)
		assert.NotNil(t, res)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		assert.Equal(t, domain.UserCore{ID: uint(1)}, res)
		repo.AssertExpectations(t)
	})

	t.Run("Fail Database Error", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(domain.UserCore{}, errors.New(config.DATABASE_ERROR)).Once()
		srv := New(repo, validator.New())
		_, err := srv.Deactivate(1)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.EqualError(t, err, config.DATABASE_ERROR)
		repo.AssertExpectations(t)
	})

	t.Run("Fail No Data", func(t *testing.T) {
		repo.On("Delete", mock.Anything).Return(domain.UserCore{}, gorm.ErrRecordNotFound).Once()
		srv := New(repo, validator.New())
		_, err := srv.Deactivate(1)
		assert.NotNil(t, err)
		assert.Error(t, err)
		assert.Equal(t, err, gorm.ErrRecordNotFound)
		repo.AssertExpectations(t)
	})
}

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
		repo.AssertExpectations(t)
	})

	t.Run("problem database", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{}, errors.New("some problem on database")).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{Fullname: "Khafid", Email: "khafid@gmail.com", Password: "khafid123456789"}
		res, err := srv.Register(input)
		assert.NotNil(t, err)
		assert.EqualError(t, err, "some problem on database", "pesan error tidak sesuai")
		assert.Empty(t, res, "karena object gagal dibuat")
		assert.Equal(t, uint(0), res.ID, "id harusnya 0 karena tidak ada data")
		repo.AssertExpectations(t)

	})

	t.Run("invalid email", func(t *testing.T) {
		var cnv rep.User
		cnv.Fullname = "joko"
		cnv.Email = "jokogmail.com"
		cnv.Password = "jokojoko1"
		srv := New(repo, validator.New())
		res, err := srv.Register(domain.UserCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: cnv.Password})
		err = validator.New().Struct(cnv)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'User.Email' Error:Field validation for 'Email' failed on the 'email' tag")
		repo.AssertExpectations(t)
	})

	t.Run("invalid fullname or password", func(t *testing.T) {
		var cnv rep.User
		cnv.Fullname = "jo"
		cnv.Email = "joko@gmail.com"
		cnv.Password = "jokojoko"
		srv := New(repo, validator.New())
		res, err := srv.Register(domain.UserCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: cnv.Password})
		err = validator.New().Struct(cnv)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'User.Fullname' Error:Field validation for 'Fullname' failed on the 'min' tag\nKey: 'User.Password' Error:Field validation for 'Password' failed on the 'containsany' tag")
		repo.AssertExpectations(t)
	})

	t.Run("encrypt password", func(t *testing.T) {
		repo.On("AddUser", mock.Anything).Return(domain.UserCore{}, errors.New("cannot encrypt password")).Once()
		var cnv rep.User
		cnv.Fullname = "joko"
		cnv.Email = "joko@gmail.com"
		cnv.Password = "jokojoko"
		srv := New(repo, validator.New())
		generate, err := bcrypt.GenerateFromPassword([]byte(cnv.Password), bcrypt.DefaultCost)
		res, err := srv.Register(domain.UserCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: string(generate)})
		assert.NotEqual(t, "$2a$10$2nyoTrlQGgzcC9mAeexI9edkq4CeiRvSf4sgGhl9DDKIiirqVytpC", string(generate))
		assert.Equal(t, uint(0), res.ID, "id harusnya 0 karena tidak ada data")
		assert.NotNil(t, err)
		assert.Error(t, err, "cannot encrypt password")
		repo.AssertExpectations(t)
	})

	t.Run("Duplicate data", func(t *testing.T) {
		repo.On("AddUser", mock.Anything, mock.Anything).Return(domain.UserCore{ID: 0x0, Fullname: "Khafid", Email: "khafid@gmail.com", Password: "$2a$10$hslY9/LyHsjfIFtSQk3mv.lWcoCinpxOoVM6GIS.6/4QmiKuVxrqG", Images: "https://bengcallbucket.s3.ap-southeast-1.amazonaws.com/profile/Q5aWl5c2RKoHcIFIrbMi-dummy450x450.jpg", Role: 0x0, Token: ""}, errors.New("Duplicate")).Once()
		var cnv rep.User
		cnv.Fullname = "Khafid"
		cnv.Email = "khafid@gmail.com"
		cnv.Password = "khafid123456789"
		srv := New(repo, validator.New())
		input := domain.UserCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: cnv.Password}
		res, err := srv.Register(input)
		assert.NotNil(t, err)
		//assert.Empty(t, res, "karena object gagal dibuat")
		assert.Equal(t, uint(0), res.ID, "id harusnya 0 karena tidak ada data")
		assert.EqualError(t, err, "already exist")
		repo.AssertExpectations(t)
	})
}

func TestUpdateProfile(t *testing.T) {
	var file multipart.File
	var FileHeader *multipart.FileHeader
	repo := mocks.NewRepository(t)

	t.Run("Sukses Update User", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{ID: uint(1), Fullname: "joko soleh", Password: "jokojoko1", Email: "joko@gmail.com"}, nil).Once()
		srv := New(repo, validator.New())
		input := domain.UserCore{ID: uint(1), Fullname: "joko soleh", Password: "jokojoko1", Email: "joko@gmail.com"}
		res, err := srv.UpdateProfile(input, file, FileHeader, 1)
		assert.Nil(t, err)
		assert.NotEmpty(t, res)
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Update Profile", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("error on bcrypt password updated user")).Once()
		srv := New(repo, validator.New())
		var input domain.UserCore
		res, err := srv.UpdateProfile(input, file, FileHeader, 1)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "some problem on database")
		repo.AssertExpectations(t)
	})

	t.Run("Gagal Update Profile", func(t *testing.T) {
		repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("column doesnt exists")).Once()
		srv := New(repo, validator.New())
		var input domain.UserCore
		res, err := srv.UpdateProfile(input, file, FileHeader, 1)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "rejected from database")
		repo.AssertExpectations(t)
	})

	t.Run("invalid email", func(t *testing.T) {
		//repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("Key: 'Email.Email' Error:Field validation for 'Email' failed on the 'email' tag")).Once()
		var cnv rep.Email
		cnv.Email = "jokogmail.com"
		srv := New(repo, validator.New())
		res, err := srv.UpdateProfile(domain.UserCore{Email: cnv.Email}, file, FileHeader, 1)
		err = validator.New().Struct(cnv)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Email.Email' Error:Field validation for 'Email' failed on the 'email' tag")
		repo.AssertExpectations(t)
	})

	t.Run("invalid fullname", func(t *testing.T) {
		//repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("Key: 'Email.Email' Error:Field validation for 'Email' failed on the 'email' tag")).Once()
		var cnv rep.Fullname
		cnv.Fullname = "jo"
		srv := New(repo, validator.New())
		res, err := srv.UpdateProfile(domain.UserCore{Fullname: cnv.Fullname}, file, FileHeader, 1)
		err = validator.New().Struct(cnv)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Fullname.Fullname' Error:Field validation for 'Fullname' failed on the 'min' tag")
		repo.AssertExpectations(t)
	})

	t.Run("invalid password", func(t *testing.T) {
		//repo.On("Update", mock.Anything, mock.Anything).Return(domain.UserCore{}, errors.New("Key: 'Email.Email' Error:Field validation for 'Email' failed on the 'email' tag")).Once()
		var cnv rep.Password
		cnv.Password = "jokojoko"
		srv := New(repo, validator.New())
		res, err := srv.UpdateProfile(domain.UserCore{Password: cnv.Password}, file, FileHeader, 1)
		err = validator.New().Struct(cnv)
		assert.Empty(t, res)
		assert.Error(t, err)
		assert.EqualError(t, err, "Key: 'Password.Password' Error:Field validation for 'Password' failed on the 'containsany' tag")
		repo.AssertExpectations(t)
	})

}
