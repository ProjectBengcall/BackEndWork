package repository

import (
	"bengcall/features/user/domain"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Fullname string `json:"fullname" validate:"required,min=3,max=20"`
	Email    string `gorm:"unique" json:"email"  validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,containsany=1234567890" `
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
