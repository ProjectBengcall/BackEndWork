package repository

import (
	"bengcall/features/user/domain"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fullname string `json:"fullname" validate:"required,alpha,min=3,max=40"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum,min=8"`
	Images   string `json:"images"`
	Role     uint   `json:"role" validate:"numeric"`
	Token    string `json:"token" validate:"multibyte"`
}

func FromDomain(du domain.UserCore) User {
	return User{
		Model:    gorm.Model{ID: du.ID},
		Fullname: du.Fullname,
		Email:    du.Email,
		Password: du.Password,
		Images:   du.Images,
		Role:     du.Role,
		Token:    du.Token,
	}
}

func ToDomain(u User) domain.UserCore {
	return domain.UserCore{
		ID:       u.ID,
		Fullname: u.Fullname,
		Email:    u.Email,
		Password: u.Password,
		Images:   u.Images,
		Role:     u.Role,
		Token:    u.Token,
	}
}
