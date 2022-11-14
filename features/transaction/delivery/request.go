package delivery

import (
	"bengcall/features/transaction/domain"
)

type TransactionFormat struct {
	Location int            `json:"location" form:"location"`
	Phone    string         `json:"phone" form:"phone"`
	Address  string         `json:"address" form:"address"`
	Schedule string         `json:"date" form:"date"`
	Other    string         `json:"other" form:"other"`
	Detail   []DetailFormat `json:"detail" form:"detail"`
}

type SuccessFormat struct {
	Status string `json:"transaction_status" form:"transaction_status"`
	Order  string `json:"order_id" form:"order_id"`
}

type DetailFormat struct {
	VehicleID uint `json:"vehicle_id" form:"vehicle_id"`
	ServiceID uint `json:"service_id" form:"service_id"`
	SubTotal  int  `json:"sub_total" form:"sub_total"`
}

type StatusFormat struct {
	Other      string `json:"other" form:"other"`
	Additional int    `json:"additional" form:"additional"`
	Status     int    `json:"status" form:"status"`
}
type CommentFormat struct {
	Comment string `json:"comment" form:"comment"`
}

func ToDom(i interface{}) []domain.DetailCore {
	switch i.(type) {
	case TransactionFormat:
		cnv := i.(TransactionFormat)
		var detail []domain.DetailCore
		for x := 0; x < len(cnv.Detail); x++ {
			detail = append(detail, domain.DetailCore{VehicleID: cnv.Detail[x].VehicleID, ServiceID: cnv.Detail[x].ServiceID, SubTotal: cnv.Detail[x].SubTotal})
		}
		return detail
	}
	return []domain.DetailCore{}
}

func ToDomain(i interface{}) domain.TransactionCore {
	switch i.(type) {
	case TransactionFormat:
		cnv := i.(TransactionFormat)
		return domain.TransactionCore{Location: cnv.Location, Phone: cnv.Phone, Address: cnv.Address, Schedule: cnv.Schedule}
	case StatusFormat:
		cnv := i.(StatusFormat)
		return domain.TransactionCore{Other: cnv.Other, Additional: cnv.Additional, Status: cnv.Status}
	case CommentFormat:
		cnv := i.(CommentFormat)
		return domain.TransactionCore{Comment: cnv.Comment}
	}
	return domain.TransactionCore{}
}

func ToSucc(i interface{}) domain.TransactionSuccess {
	switch i.(type) {
	case SuccessFormat:
		cnv := i.(SuccessFormat)
		return domain.TransactionSuccess{Status: cnv.Status, Order: cnv.Order}
	}
	return domain.TransactionSuccess{}
}
