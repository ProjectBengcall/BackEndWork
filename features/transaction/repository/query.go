package repository

import (
	"bengcall/features/transaction/domain"
	"strconv"
	"time"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/snap"
	"gorm.io/gorm"
)

type repoQuery struct {
	db *gorm.DB
}

func New(db *gorm.DB) domain.Repository {
	return &repoQuery{
		db: db,
	}
}

var s snap.Client

func (rq *repoQuery) Post(newTrx domain.TransactionCore, newDtl []domain.DetailCore) (domain.TransactionDetail, error) {
	var resQry TransactionComplete
	var total int

	for i := 0; i < len(newDtl); i++ {
		if err := rq.db.Exec("INSERT INTO details (id, created_at, updated_at, deleted_at, vehicle_id, service_id, transaction_id, sub_total) values (?,?,?,?,?,?,?,?)",
			nil, time.Now(), time.Now(), nil, newDtl[i].VehicleID, newDtl[i].ServiceID, newTrx.Invoice, newDtl[i].SubTotal).Error; err != nil {
			return domain.TransactionDetail{}, err
		}
		total += newDtl[i].SubTotal
	}

	v := strconv.Itoa(newTrx.Invoice)

	// 1. Initiate Snap client
	s.New("SB-Mid-server-eKSCGMJJG-IEL_LscFEV9-nP", midtrans.Sandbox)

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  v,
			GrossAmt: int64(total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	newTrx.PaymentToken = snapResp.Token
	newTrx.PaymentLink = snapResp.RedirectURL

	if err := rq.db.Exec("INSERT INTO transactions (id, created_at, updated_at, deleted_at, location, phone, address, invoice, total, payment_token, payment_link, schedule, status, user_id) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		nil, time.Now(), time.Now(), nil, newTrx.Location, newTrx.Phone, newTrx.Address, newTrx.Invoice, total, newTrx.PaymentToken, newTrx.PaymentLink, newTrx.Schedule, 1, newTrx.UserID).Error; err != nil {
		return domain.TransactionDetail{}, err
	}

	if er := rq.db.Table("transactions").Select("id", "invoice", "total", "status", "payment_token", "payment_link").Where("invoice = ?", newTrx.Invoice).Model(&TransactionComplete{}).Find(&resQry).Error; er != nil {
		return domain.TransactionDetail{}, er
	}
	res := ToDomDetail(resQry)
	return res, nil
}

func (rq *repoQuery) PutScss(ID uint) error {
	if err := rq.db.Exec("UPDATE transactions SET status = ? WHERE invoice = ?", 3, ID).Error; err != nil {
		return err
	}
	return nil
}

func (rq *repoQuery) PutStts(updateStts domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	updateStts.ID = ID
	var resQry Transaction

	if err := rq.db.Exec("UPDATE transactions SET total = total + ?, other = ?, additional = ?, status = ? WHERE id = ?",
		updateStts.Additional, updateStts.Other, updateStts.Additional, updateStts.Status, ID).Error; err != nil {
		return domain.TransactionCore{}, err
	}

	if er := rq.db.Table("transactions").Select("invoice", "total").Where("id = ?", ID).Model(&TransactionComplete{}).Find(&resQry).Error; er != nil {
		return domain.TransactionCore{}, er
	}

	v := strconv.Itoa(resQry.Invoice)

	// 1. Initiate Snap client
	s.New("SB-Mid-server-eKSCGMJJG-IEL_LscFEV9-nP", midtrans.Sandbox)

	// 2. Initiate Snap request param
	req := &snap.Request{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  v,
			GrossAmt: int64(resQry.Total),
		},
		CreditCard: &snap.CreditCardDetails{
			Secure: true,
		},
	}

	// 3. Execute request create Snap transaction to Midtrans Snap API
	snapResp, _ := s.CreateTransaction(req)

	newPaymentToken := snapResp.Token
	newPaymentLink := snapResp.RedirectURL

	if err := rq.db.Exec("UPDATE transactions SET payment_token = ?, payment_link = ? WHERE id = ?",
		newPaymentToken, newPaymentLink, ID).Error; err != nil {
		return domain.TransactionCore{}, err
	}

	if er := rq.db.Table("transactions").Select("id", "payment_token", "payment_link", "status").Where("id = ?", ID).Model(&Transaction{}).Find(&resQry).Error; er != nil {
		return domain.TransactionCore{}, er
	}

	res := ToDomStts(resQry)
	return res, nil
}

func (rq *repoQuery) PutCmmt(updateCmmt domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	updateCmmt.ID = ID
	var cnv Transaction = FromDomCmmt(updateCmmt)
	if err := rq.db.Exec("UPDATE transactions SET comment = ? WHERE id = ?",
		updateCmmt.Comment, ID).Error; err != nil {
		return domain.TransactionCore{}, err
	}
	res := ToDomCmmt(cnv)
	return res, nil
}

func (rq *repoQuery) GetAll() ([]domain.TransactionAll, error) {
	var resQry []TransactionComplete
	if err := rq.db.Table("transactions").Select("transactions.id", "transactions.schedule", "transactions.invoice", "transactions.total", "transactions.status", "users.fullname").Joins("join users on users.id=transactions.user_id").Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
		return nil, err
	}
	res := ToDomainAll(resQry)
	return res, nil
}

func (rq *repoQuery) GetMy(userID uint) (domain.TransactionHistory, error) {
	var resQry TransactionComplete
	if err := rq.db.Table("transactions").Select("transactions.id", "transactions.schedule", "transactions.invoice", "transactions.total").Where("user_id = ? && status != 3", userID).Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
		return domain.TransactionHistory{}, err
	}
	res := ToDomHistory(resQry)
	return res, nil
}

func (rq *repoQuery) GetHistory(userID uint) ([]domain.TransactionHistory, error) {
	var resQry []TransactionComplete
	if err := rq.db.Table("transactions").Select("transactions.id", "transactions.schedule", "transactions.invoice", "transactions.total").Where("user_id = ?", userID).Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
		return nil, err
	}
	res := ToDomainHistory(resQry)
	return res, nil
}

func (rq *repoQuery) GetDetail(ID uint) (domain.TransactionDetail, error) {
	var resQry TransactionComplete
	if err := rq.db.Table("transactions").Select("transactions.id", "transactions.location", "transactions.schedule", "transactions.phone", "transactions.address", "transactions.invoice", "transactions.total", "transactions.payment_token", "transactions.payment_link", "transactions.other", "transactions.status", "users.fullname", "vehicles.name_vehicle", "services.service_name").Joins("join users on users.id=transactions.user_id").Joins("join details on details.transaction_id=transactions.invoice").Joins("join vehicles on vehicles.id=details.vehicle_id").Joins("join services on services.id=details.service_id").Where("transactions.id = ?", ID).Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
		return domain.TransactionDetail{}, err
	}
	res := ToDomDetail(resQry)
	return res, nil
}

func (rq *repoQuery) Delete(ID uint) error {
	var resQry Transaction
	if err := rq.db.Where("id = ?", ID).Delete(&resQry).Error; err != nil {
		return err
	}
	return nil
}
