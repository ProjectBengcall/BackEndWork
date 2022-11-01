package delivery

import "bengcall/features/user/domain"

type LoginResponse struct {
	Fullname string `json:"fullname"`
	Token    string `json:"token"`
}

type RegistResponses struct {
	Fullname string `json:"username"`
	Email    string `json:"email"`
	Images   string `json:"images"`
	Role     uint   `json:"role"`
}

type EditResponse struct {
	Fullname string `json:"username"`
	Email    string `json:"email"`
	Images   string `json:"images"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Fullname string `json:"username"`
	Email    string `json:"email"`
	Images   string `json:"images"`
}

func ToResponse(core interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "login":
		cnv := core.(domain.UserCore)
		res = LoginResponse{Fullname: cnv.Fullname, Token: cnv.Token}
	case "reg":
		cnv := core.(domain.UserCore)
		res = RegistResponses{Fullname: cnv.Fullname, Email: cnv.Email, Images: cnv.Images, Role: cnv.Role}
	case "edit":
		cnv := core.(domain.UserCore)
		res = RegistResponses{Fullname: cnv.Fullname, Email: cnv.Email, Images: cnv.Images}
	case "user":
		cnv := core.(domain.UserCore)
		res = UserResponse{ID: cnv.ID, Fullname: cnv.Fullname, Email: cnv.Email, Images: cnv.Images}
	}
	return res
}

func SuccessResponse(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		"data":    data,
	}
}

func SuccessLogin(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		"data":    data,
	}
}

func FailResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
	}
}
