package delivery

import "bengcall/features/transaction/domain"

// type TransactionFormat struct {
// 	Fullname string `json:"fullname" form:"fullname" validate:"required,min=3,max=20" `
// 	Email    string `json:"email" form:"email" validate:"required,email"`
// 	Password string `json:"password" form:"password" validate:"required,min=8,containsany=1234567890"`
// }

type StatusFormat struct {
	Other      string `json:"other" form:"other"`
	Additional int    `json:"additional" form:"additional"`
	Status     int    `json:"status" form:"status"`
}
type CommentFormat struct {
	Comment string `json:"comment" form:"comment"`
}

func ToDomain(i interface{}) domain.TransactionCore {
	switch i.(type) {
	// case TransactionFormat:
	// 	cnv := i.(TransactionFormat)
	// 	return domain.TransactionCore{Fullname: cnv.Fullname, Email: cnv.Email, Password: cnv.Password}
	case StatusFormat:
		cnv := i.(StatusFormat)
		return domain.TransactionCore{Other: cnv.Other, Additional: cnv.Additional, Status: cnv.Status}
	case CommentFormat:
		cnv := i.(CommentFormat)
		return domain.TransactionCore{Comment: cnv.Comment}
	}
	return domain.TransactionCore{}
}
