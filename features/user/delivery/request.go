package delivery

import "bengcall/features/user/domain"

type UserFormat struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Images   string `json:"images" form:"images"`
	Role     uint   `json:"role" form:"role"`
}

type LoginFormat struct {
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
}
type EditFormat struct {
	Fullname string `json:"fullname" form:"fullname"`
	Email    string `json:"email" form:"email"`
	Password string `json:"password" form:"password"`
	Images   string `json:"images" form:"images"`
}

func ToDomain(i interface{}) domain.UserCore {
	switch i.(type) {
	case UserFormat:
		cnv := i.(UserFormat)
		return domain.UserCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: cnv.Password}
	case LoginFormat:
		cnv := i.(LoginFormat)
		return domain.UserCore{Email: cnv.Email, Password: cnv.Password}
	case EditFormat:
		cnv := i.(EditFormat)
		return domain.UserCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: cnv.Password, Images: cnv.Images}
	}
	return domain.UserCore{}
}
