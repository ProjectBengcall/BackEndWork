package repository

import (
	"bengcall/features/transaction/domain"
	"time"

	"gorm.io/gorm"
)

type Transaction struct {
	gorm.Model
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

type Detail struct {
	gorm.Model
	VehicleID     uint
	ServiceID     uint
	TransactionID uint
	SubTotal      int
}

type TransactionComplete struct {
	gorm.Model
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
	Fullname     string
	Name_vehicle string
	ServiceName  string
}

func FromDomStts(du domain.TransactionCore) Transaction {
	return Transaction{
		Model:      gorm.Model{ID: du.ID},
		Other:      du.Other,
		Additional: du.Additional,
		Status:     du.Status,
	}
}

func FromDomCmmt(du domain.TransactionCore) Transaction {
	return Transaction{
		Model:   gorm.Model{ID: du.ID},
		Comment: du.Comment,
	}
}

func ToDomStts(t Transaction) domain.TransactionCore {
	return domain.TransactionCore{
		ID:     t.ID,
		Status: t.Status,
	}
}

func ToDomCmmt(t Transaction) domain.TransactionCore {
	return domain.TransactionCore{
		ID:      t.ID,
		Comment: t.Comment,
	}
}

func ToDomDetail(t TransactionComplete) domain.TransactionDetail {
	return domain.TransactionDetail{
		ID:           t.ID,
		Location:     t.Location,
		Schedule:     t.Schedule,
		Phone:        t.Phone,
		Address:      t.Address,
		Invoice:      t.Invoice,
		Total:        t.Total,
		PaymentToken: t.PaymentToken,
		PaymentLink:  t.PaymentLink,
		Other:        t.Other,
		Status:       t.Status,
		Fullname:     t.Fullname,
		Name_vehicle: t.Name_vehicle,
		ServiceName:  t.ServiceName,
	}
}

func ToDomainAll(ta []TransactionComplete) []domain.TransactionAll {
	var res []domain.TransactionAll
	for _, val := range ta {
		res = append(res, domain.TransactionAll{ID: val.ID, Schedule: val.Schedule, Invoice: val.Invoice, Total: val.Total, Status: val.Status, Fullname: val.Fullname})
	}
	return res
}

func ToDomainHistory(ta []TransactionComplete) []domain.TransactionHistory {
	var res []domain.TransactionHistory
	for _, val := range ta {
		res = append(res, domain.TransactionHistory{ID: val.ID, Schedule: val.Schedule, Invoice: val.Invoice, Total: val.Total})
	}
	return res
}
