package delivery

import "bengcall/features/user/domain"

type LoginResponse struct {
	Fullname string `json:"fullname"`
	Images   string `json:"images"`
	Role     uint   `json:"role"`
	Token    string `json:"token"`
}

type RegistResponses struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
}

type EditResponse struct {
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Images   string `json:"images"`
}

type UserResponse struct {
	ID       uint   `json:"id"`
	Fullname string `json:"fullname"`
	Email    string `json:"email"`
	Images   string `json:"images"`
}

func ToResponse(core interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "login":
		cnv := core.(domain.UserCore)
		res = LoginResponse{Fullname: cnv.Fullname, Images: cnv.Images, Role: cnv.Role, Token: cnv.Token}
	case "reg":
		cnv := core.(domain.UserCore)
		res = RegistResponses{Fullname: cnv.Fullname, Email: cnv.Email}
	case "edit":
		cnv := core.(domain.UserCore)
		res = EditResponse{Fullname: cnv.Fullname, Email: cnv.Email, Images: cnv.Images}
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

func SuccessDeleteResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
	}
}
