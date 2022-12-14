package repository

import (
	"bengcall/features/transaction/domain"
	"errors"
	"math/rand"
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

	if err := rq.db.Exec("INSERT INTO transactions (id, created_at, updated_at, deleted_at, location, phone, address, invoice, total, other, payment_token, payment_link, schedule, status, user_id) values (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)",
		nil, time.Now(), time.Now(), nil, newTrx.Location, newTrx.Phone, newTrx.Address, newTrx.Invoice, total, newTrx.Other, newTrx.PaymentToken, newTrx.PaymentLink, newTrx.Schedule, 1, newTrx.UserID).Error; err != nil {
		return domain.TransactionDetail{}, err
	}

	if er := rq.db.Table("transactions").Select("transactions.id", "transactions.invoice", "transactions.total", "transactions.status", "transactions.payment_token", "transactions.payment_link", "users.email").Joins("join users on users.id=transactions.user_id").Where("invoice = ?", newTrx.Invoice).Model(&TransactionComplete{}).Find(&resQry).Error; er != nil {
		return domain.TransactionDetail{}, er
	}
	res := ToDomDetail(resQry)
	return res, nil
}

func (rq *repoQuery) PutScss(ID uint) error {
	if err := rq.db.Exec("UPDATE transactions SET status = 6 WHERE invoice = ?", ID).Error; err != nil {
		return err
	}
	return nil
}

func (rq *repoQuery) PutStts(updateStts domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	updateStts.ID = ID
	var resQry Transaction
	var detQry Detail

	if updateStts.Additional > 0 {
		if updateStts.Status != 0 {
			if err := rq.db.Exec("UPDATE transactions SET total = total + ?, other = ?, additional = ?, status = ? WHERE id = ?",
				updateStts.Additional, updateStts.Other, updateStts.Additional, updateStts.Status, ID).Error; err != nil {
				return domain.TransactionCore{}, err
			}
		} else {
			if err := rq.db.Exec("UPDATE transactions SET total = total + ?, other = ?, additional = ? WHERE id = ?",
				updateStts.Additional, updateStts.Other, updateStts.Additional, ID).Error; err != nil {
				return domain.TransactionCore{}, err
			}
		}

		if er := rq.db.Table("transactions").Select("total", "invoice").Where("id = ?", ID).Model(&TransactionComplete{}).Find(&resQry).Error; er != nil {
			return domain.TransactionCore{}, er
		}

		if er := rq.db.Table("details").Select("id").Where("transaction_id = ?", resQry.Invoice).Model(&Detail{}).Find(&detQry).Error; er != nil {
			return domain.TransactionCore{}, er
		}

		var invo int
		rand.Seed(time.Now().UnixNano())
		invo = rand.Intn(100000)
		v := strconv.Itoa(invo)

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

		if err := rq.db.Exec("UPDATE transactions SET invoice = ?, payment_token = ?, payment_link = ? WHERE id = ?",
			invo, newPaymentToken, newPaymentLink, ID).Error; err != nil {
			return domain.TransactionCore{}, err
		}

		if err := rq.db.Exec("UPDATE details SET transaction_id = ? WHERE id = ?",
			invo, detQry.ID).Error; err != nil {
			return domain.TransactionCore{}, err
		}
	} else {
		if err := rq.db.Exec("UPDATE transactions SET status = ? WHERE id = ?",
			updateStts.Status, ID).Error; err != nil {
			return domain.TransactionCore{}, err
		}
	}

	if er := rq.db.Table("transactions").Select("id", "invoice", "payment_token", "payment_link", "status").Where("id = ?", ID).Model(&Transaction{}).Find(&resQry).Error; er != nil {
		return domain.TransactionCore{}, er
	}

	res := ToDomStts(resQry)
	return res, nil
}

func (rq *repoQuery) PutCmmt(updateCmmt domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	updateCmmt.ID = ID
	var resQry Transaction
	var cnv Transaction = FromDomCmmt(updateCmmt)

	if er := rq.db.Table("transactions").Select("id").Where("id = ?", ID).Model(&Transaction{}).Find(&resQry).Error; er != nil {
		return domain.TransactionCore{}, er
	}

	if resQry.ID == ID {
		if err := rq.db.Exec("UPDATE transactions SET comment = ? WHERE id = ?",
			updateCmmt.Comment, ID).Error; err != nil {
			return domain.TransactionCore{}, err
		}
		res := ToDomCmmt(cnv)
		return res, nil
	} else {
		return domain.TransactionCore{}, errors.New("id not recognize")
	}
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
	if err := rq.db.Table("transactions").Select("transactions.id", "transactions.schedule", "transactions.invoice", "transactions.total").Where("user_id = ? && status != 6", userID).Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
		return domain.TransactionHistory{}, err
	}
	res := ToDomHistory(resQry)
	return res, nil
}

func (rq *repoQuery) GetHistory(userID uint) ([]domain.TransactionHistory, error) {
	var resQry []TransactionComplete
	if err := rq.db.Table("transactions").Select("transactions.id", "transactions.schedule", "transactions.invoice", "transactions.total").Where("user_id = ? && status = 6", userID).Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
		return nil, err
	}
	res := ToDomainHistory(resQry)
	return res, nil
}

func (rq *repoQuery) GetDetail(ID uint) (domain.TransactionDetail, []domain.DetailCores, error) {
	var resQry TransactionComplete
	var dtlQry []Details

	if er := rq.db.Table("transactions").Select("id").Where("id = ?", ID).Model(&TransactionComplete{}).Find(&resQry).Error; er != nil {
		return domain.TransactionDetail{}, nil, er
	}

	if resQry.ID == uint(ID) {

		if err := rq.db.Table("transactions").Select("transactions.id", "transactions.location", "transactions.schedule", "transactions.phone", "transactions.address", "transactions.invoice", "transactions.total", "transactions.payment_token", "transactions.payment_link", "transactions.other", "transactions.additional", "transactions.status", "users.fullname").Joins("join users on users.id=transactions.user_id").Where("transactions.id = ?", ID).Model(&TransactionComplete{}).Find(&resQry).Error; err != nil {
			return domain.TransactionDetail{}, nil, err
		}
		if err := rq.db.Table("details").Select("details.id", "vehicles.name_vehicle", "services.service_name", "details.transaction_id", "details.sub_total").Joins("join vehicles on vehicles.id=details.vehicle_id").Joins("join services on services.id=details.service_id").Where("details.transaction_id = ?", resQry.Invoice).Model(&Details{}).Find(&dtlQry).Error; err != nil {
			return domain.TransactionDetail{}, nil, err
		}

		dtl := ToDomDetails(dtlQry)
		res := ToDomDetail(resQry)
		return res, dtl, nil
	} else {
		return domain.TransactionDetail{}, nil, errors.New("id not recognize")
	}
}

func (rq *repoQuery) Delete(ID uint) error {
	var resQry Transaction
	if err := rq.db.Where("id = ?", ID).Delete(&resQry).Error; err != nil {
		return err
	}
	return nil
}
