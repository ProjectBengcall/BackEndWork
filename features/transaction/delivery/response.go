package delivery

import (
	"bengcall/features/transaction/domain"
	"time"
)

func SuccessResponse(msg string, data interface{}) map[string]interface{} {
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

type ResponseAll struct {
	ID       uint      `json:"id"`
	Invoice  int       `json:"invoice"`
	Fullname string    `json:"fullname"`
	Schedule time.Time `json:"date"`
	Total    int       `json:"total"`
	Status   int       `json:"status"`
}

type ResponseHistory struct {
	ID       uint      `json:"id"`
	Invoice  int       `json:"invoice"`
	Schedule time.Time `json:"date"`
	Total    int       `json:"total"`
}

type ResponseDetail struct {
	ID           uint      `json:"id"`
	Location     int       `json:"location"`
	Fullname     string    `json:"fullname"`
	Phone        string    `json:"phone"`
	Address      string    `json:"address"`
	Schedule     time.Time `json:"date"`
	Name_vehicle string    `json:"vehicle"`
	ServiceName  string    `json:"service"`
	Other        string    `json:"other"`
	Invoice      int       `json:"invoice"`
	Total        int       `json:"total"`
	Status       int       `json:"status"`
	PaymentToken string    `json:"payment_token"`
	PaymentLink  string    `json:"payment_link"`
}

type ResponsePost struct {
	ID           uint   `json:"id"`
	Invoice      int    `json:"invoice"`
	Total        int    `json:"total"`
	Status       int    `json:"status"`
	PaymentToken string `json:"payment_token"`
	PaymentLink  string `json:"payment_link"`
}

type ResponseStts struct {
	ID     uint `json:"id"`
	Status int  `json:"status"`
}

type ResponseCmmt struct {
	ID      uint   `json:"id"`
	Comment string `json:"comment"`
}

func ToResponse(core interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "stts":
		cnv := core.(domain.TransactionCore)
		res = ResponseStts{ID: cnv.ID, Status: cnv.Status}
	case "cmmt":
		cnv := core.(domain.TransactionCore)
		res = ResponseCmmt{ID: cnv.ID, Comment: cnv.Comment}
	case "detail":
		cnv := core.(domain.TransactionDetail)
		res = ResponseDetail{ID: cnv.ID, Location: cnv.Location, Fullname: cnv.Fullname, Phone: cnv.Phone, Address: cnv.Address, Schedule: cnv.Schedule, Name_vehicle: cnv.Name_vehicle, ServiceName: cnv.ServiceName, Other: cnv.Other, Invoice: cnv.Invoice, Total: cnv.Total, Status: cnv.Status, PaymentToken: cnv.PaymentToken, PaymentLink: cnv.PaymentLink}
	case "post":
		cnv := core.(domain.TransactionDetail)
		res = ResponsePost{ID: cnv.ID, Invoice: cnv.Invoice, Total: cnv.Total, Status: cnv.Status, PaymentToken: cnv.PaymentToken, PaymentLink: cnv.PaymentLink}
	case "all":
		var arr []ResponseAll
		cnv := core.([]domain.TransactionAll)
		for _, val := range cnv {
			arr = append(arr, ResponseAll{ID: val.ID, Invoice: val.Invoice, Fullname: val.Fullname, Schedule: val.Schedule, Total: val.Total, Status: val.Status})
		}
		res = arr
	case "history":
		var arr []ResponseHistory
		cnv := core.([]domain.TransactionHistory)
		for _, val := range cnv {
			arr = append(arr, ResponseHistory{ID: val.ID, Invoice: val.Invoice, Schedule: val.Schedule, Total: val.Total})
		}
		res = arr
	}
	return res
}
