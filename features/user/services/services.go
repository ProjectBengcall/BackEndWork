package services

import (
	"bengcall/features/user/domain"
	rep "bengcall/features/user/repository"
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	qry      domain.Repository
	validate *validator.Validate
}

func New(repo domain.Repository, val *validator.Validate) domain.Service {
	return &userService{
		qry:      repo,
		validate: val,
	}
}

// Deactivate implements domain.Service
func (us *userService) Deactivate(userID uint) error {
	err := us.qry.Delete(userID)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return errors.New("no data")
		}
	}
	return nil
}

// Login implements domain.Service
func (us *userService) Login(existUser domain.UserCore) (domain.UserCore, error) {
	var cnv = rep.FromDomain(existUser)
	if strings.TrimSpace(cnv.Email) == "" || strings.TrimSpace(cnv.Password) == "" {
		return domain.UserCore{}, errors.New("Email or password empty")
	}

	res, err := us.qry.GetUser(existUser)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return domain.UserCore{}, errors.New("Failed. Error database.")
		} else if strings.Contains(err.Error(), "found") {
			return domain.UserCore{}, errors.New("Failed. Email or Password not found.")
		} else {
			return domain.UserCore{}, errors.New("Failed. Process error. Please contact Admin")
		}
	}

	//pass := domain.UserCore{Password: res.Password}
	err = bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(existUser.Password))
	if err != nil {
		log.Error(err, " wrong password")
		if strings.Contains(err.Error(), "found") {
			return domain.UserCore{}, errors.New("Failed. Email or Password not found.")
		} else {
			return domain.UserCore{}, errors.New("wrong password")
		}
	}
	return res, nil

}

// MyProfile implements domain.Service
func (us *userService) MyProfile(userID uint) (domain.UserCore, error) {
	res, err := us.qry.GetMyUser(userID)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return domain.UserCore{}, errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.UserCore{}, errors.New("no data")
		}
	}
	return res, nil
}

// Register implements domain.Service
func (us *userService) Register(newUser domain.UserCore) (domain.UserCore, error) {
	var cnv = rep.FromDomain(newUser)
	err := us.validate.Struct(cnv)
	if err != nil {
		if strings.Contains(err.Error(), "email") {
			return domain.UserCore{}, errors.New("invalid email")
		} else {
			return domain.UserCore{}, errors.New("invalid password")
		}
	}

	generate, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error("error on bcrypt", err.Error())
		return domain.UserCore{}, errors.New("cannot encrypt password")
	}

	newUser.Password = string(generate)
	newUser.Images = "https://bengcallbucket.s3.ap-southeast-1.amazonaws.com/profile/Q5aWl5c2RKoHcIFIrbMi-dummy450x450.jpg"
	newUser.Role = 0
	orgPass := newUser.Password
	res, err := us.qry.AddUser(newUser)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("email alredy exist")
	}
	res.Password = orgPass
	return res, nil

}

// UpdateProfile implements domain.Service
func (us *userService) UpdateProfile(updatedUser domain.UserCore, userID uint) (domain.UserCore, error) {
	if updatedUser.Password != "" {
		generate, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Error("error on bcrypt password updated user", err.Error())
			return domain.UserCore{}, errors.New("cannot encrypt password")
		}
		updatedUser.Password = string(generate)
	}
	// var cnv = rep.FromDomain(updatedUser)
	// err := us.validate.Struct(cnv)
	// if err != nil {
	// 	log.Error("Validation errror : ", err.Error())
	// 	return domain.UserCore{}, err
	// }

	res, err := us.qry.Update(updatedUser, userID)
	if err != nil {
		if strings.Contains(err.Error(), "column") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("some problem on database")
	}

	return res, nil
}
