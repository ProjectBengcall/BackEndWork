package services

import (
	"bengcall/features/transaction/domain"
	"errors"
	"math/rand"
	"strings"
	"time"

	"github.com/labstack/gommon/log"
)

type transactionService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &transactionService{
		qry: repo,
	}
}

func (ss *transactionService) Transaction(newTrx domain.TransactionCore, newDtl []domain.DetailCore) (domain.TransactionDetail, error) {
	var invo int
	rand.Seed(time.Now().UnixNano())
	invo = rand.Intn(100000)
	newTrx.Invoice = invo

	res, err := ss.qry.Post(newTrx, newDtl)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionDetail{}, errors.New("rejected from database")
		}
		return domain.TransactionDetail{}, errors.New("some problem on database")
	}

	return res, nil
}

func (ss *transactionService) Success(ID uint) error {
	err := ss.qry.PutScss(ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return errors.New("Rejected from Database")
		}
		return errors.New("Some Problem on Database")
	}

	return nil
}

func (ss *transactionService) Status(updateStts domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	res, err := ss.qry.PutStts(updateStts, ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionCore{}, errors.New("Rejected from Database")
		}
		return domain.TransactionCore{}, errors.New("Some Problem on Database")
	}

	return res, nil
}

func (ss *transactionService) Comment(updateCmmt domain.TransactionCore, ID uint) (domain.TransactionCore, error) {
	res, err := ss.qry.PutCmmt(updateCmmt, ID)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.TransactionCore{}, errors.New("Rejected from Database")
		} else if strings.Contains(err.Error(), "id") {
			return domain.TransactionCore{}, errors.New("There's no ID")
		}
		return domain.TransactionCore{}, errors.New("Some Problem on Database")
	}

	return res, nil
}

func (ss *transactionService) All() ([]domain.TransactionAll, error) {
	res, err := ss.qry.GetAll()
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return nil, errors.New("No Data")
		}
	}

	return res, nil
}

func (ss *transactionService) History(userID uint) ([]domain.TransactionHistory, error) {
	res, err := ss.qry.GetHistory(userID)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return nil, errors.New("No Data")
		}
	}
	return res, nil
}

func (ss *transactionService) My(userID uint) (domain.TransactionHistory, error) {
	res, err := ss.qry.GetMy(userID)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return domain.TransactionHistory{}, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.TransactionHistory{}, errors.New("No Data")
		}
	}
	return res, nil
}

func (ss *transactionService) Detail(ID uint) (domain.TransactionDetail, []domain.DetailCores, error) {
	res, dtl, err := ss.qry.GetDetail(ID)
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return domain.TransactionDetail{}, nil, errors.New("Database Error")
		} else if strings.Contains(err.Error(), "found") {
			return domain.TransactionDetail{}, nil, errors.New("No Data")
		}
	}

	return res, dtl, nil
}

func (ss *transactionService) Cancel(ID uint) error {
	err := ss.qry.Delete(ID)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return errors.New("no data")
		}
	}
	return nil
}
