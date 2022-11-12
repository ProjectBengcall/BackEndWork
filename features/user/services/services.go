package services

import (
	"bengcall/config"
	"bengcall/features/user/domain"
	rep "bengcall/features/user/repository"
	"bengcall/utils/helper"
	"errors"
	lo "log"
	"mime/multipart"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/gommon/log"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
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
func (us *userService) Deactivate(userID uint) (domain.UserCore, error) {
	res, err := us.qry.Delete(userID)
	if err == gorm.ErrRecordNotFound {
		log.Error(err.Error())
		return domain.UserCore{}, gorm.ErrRecordNotFound
	} else if err != nil {
		log.Error(err.Error())
		return domain.UserCore{}, errors.New(config.DATABASE_ERROR)
	}
	return res, nil
}

// Login implements domain.Service
func (us *userService) Login(existUser domain.UserCore) (domain.UserCore, error) {

	res, err := us.qry.GetUser(existUser)
	lo.Println("hasil res", res)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "found") {
			return domain.UserCore{}, errors.New("Failed. Email or Password not found.")
		} else if strings.Contains(err.Error(), "table") {
			return domain.UserCore{}, errors.New("Failed. Email or Password not found.")
		}
		return domain.UserCore{}, errors.New("email not found")
	} else {
		if res.ID == 0 {
			return domain.UserCore{}, errors.New("email not found")
		}
		err := bcrypt.CompareHashAndPassword([]byte(res.Password), []byte(existUser.Password))
		if err != nil {
			log.Error(err, " wrong password")
			return domain.UserCore{}, errors.New("wrong password")
		}
		return res, nil
	}

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
		//return domain.UserCore{}, err
		if strings.Contains(err.Error(), "email") {
			log.Error("error on validation", err.Error())
			return domain.UserCore{}, errors.New("invalid email")
		} else {
			return domain.UserCore{}, errors.New("invalid password or fullname")
		}
	} else {

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
			if strings.Contains(err.Error(), "Duplicate") {
				return domain.UserCore{}, errors.New("already exist")
			}
			return domain.UserCore{}, errors.New("some problem on database")
		}
		res.Password = orgPass
		return res, nil
	}

}

// UpdateProfile implements domain.Service
func (us *userService) UpdateProfile(updatedUser domain.UserCore, file multipart.File, fileheader *multipart.FileHeader, userID uint) (domain.UserCore, error) {
	if updatedUser.Password != "" {
		var cnv rep.Password
		cnv.Password = updatedUser.Password
		err := us.validate.Struct(cnv)
		if err != nil {
			log.Error("error on validation", err.Error())
			return domain.UserCore{}, errors.New("invalid password")
		} else {
			generate, err := bcrypt.GenerateFromPassword([]byte(updatedUser.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Error("error on bcrypt password updated user", err.Error())
				return domain.UserCore{}, errors.New("cannot encrypt password")
			} else {

				updatedUser.Password = string(generate)
			}
		}
	}

	if updatedUser.Email != "" {
		var cnv rep.Email
		cnv.Email = updatedUser.Email
		err := us.validate.Struct(cnv)
		if err != nil {
			log.Error("error on validation", err.Error())
			return domain.UserCore{}, errors.New("invalid email")
		}
	}

	if updatedUser.Fullname != "" {
		var cnv rep.Fullname
		cnv.Fullname = updatedUser.Fullname
		//Fullname := updatedUser.Fullname
		err := us.validate.Struct(cnv)
		if err != nil {
			log.Error("error on validation", err.Error())
			return domain.UserCore{}, errors.New("invalid fullname")
		}

	}

	if fileheader != nil {
		res, err := helper.UploadProfile(file, fileheader)
		if err != nil {
			return domain.UserCore{}, err
		}
		updatedUser.Images = res
	}

	res, err := us.qry.Update(updatedUser, userID)
	if err != nil {
		if strings.Contains(err.Error(), "column") {
			return domain.UserCore{}, errors.New("rejected from database")
		}
		return domain.UserCore{}, errors.New("some problem on database")
	}

	return res, nil
}
