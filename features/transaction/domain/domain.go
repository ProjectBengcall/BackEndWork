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

type TransactionCores struct {
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
	Detail       []DetailCores
}

type DetailCores struct {
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
	GetMy(userID uint) (TransactionHistory, error)
	GetHistory(userID uint) ([]TransactionHistory, error)
	GetDetail(ID uint) (TransactionDetail, error)
	Post(newTrx TransactionCore, newDtl []DetailCore) (TransactionDetail, error)
	PutScss(ID uint) error
	PutStts(updateStts TransactionCore, ID uint) (TransactionCore, error)
	PutCmmt(updateCmmt TransactionCore, ID uint) (TransactionCore, error)
	Delete(ID uint) error
}

type Service interface {
	All() ([]TransactionAll, error)
	My(userID uint) (TransactionHistory, error)
	History(userID uint) ([]TransactionHistory, error)
	Detail(ID uint) (TransactionDetail, error)
	Transaction(newTrx TransactionCore, newDtl []DetailCore) (TransactionDetail, error)
	Success(ID uint) error
	Status(updateStts TransactionCore, ID uint) (TransactionCore, error)
	Comment(updateCmmt TransactionCore, ID uint) (TransactionCore, error)
	Cancel(ID uint) error
}
