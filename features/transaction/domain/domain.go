package domain

import "time"

type TransactionCore struct {
	ID           uint
	Location     int
	Schedule     time.Time
	Phone        string
	Address      string
	Invoice      int
	Total        int
	PaymentToken string
	PaymentLink  string
	Other        string
	Additional   int
	Status       int
	Comment      string
	UserID       uint
}

type DetailCore struct {
	ID            uint
	VehicleID     uint
	ServiceID     uint
	TransactionID uint
	SubTotal      int
}

type TransactionAll struct {
	ID       uint
	Schedule time.Time
	Invoice  int
	Total    int
	Status   int
	Fullname string
}

type TransactionHistory struct {
	ID       uint
	Schedule time.Time
	Invoice  int
	Total    int
}

type TransactionDetail struct {
	ID           uint
	Location     int
	Schedule     time.Time
	Phone        string
	Address      string
	Invoice      int
	Total        int
	PaymentToken string
	PaymentLink  string
	Other        string
	Status       int
	Fullname     string
	Name_vehicle string
	ServiceName  string
}

type Repository interface {
	GetAll() ([]TransactionAll, error)
	//GetMy(userID uint) (TransactionDetail, error)
	GetHistory(userID uint) ([]TransactionHistory, error)
	GetDetail(ID uint) (TransactionDetail, error)
	PutStts(updateStts TransactionCore, ID uint) (TransactionCore, error)
	PutCmmt(updateCmmt TransactionCore, ID uint) (TransactionCore, error)
	Delete(ID uint) error
}

type Service interface {
	All() ([]TransactionAll, error)
	//My(userID uint) (TransactionDetail, error)
	History(userID uint) ([]TransactionHistory, error)
	Detail(ID uint) (TransactionDetail, error)
	Status(updateStts TransactionCore, ID uint) (TransactionCore, error)
	Comment(updateCmmt TransactionCore, ID uint) (TransactionCore, error)
	Cancel(ID uint) error
}
