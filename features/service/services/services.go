package services

import (
	"bengcall/features/service/domain"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
)

type servService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &servService{
		qry: repo,
	}
}

func (ss *servService) GetSpesific(vehicleID int) ([]domain.Core, error) {
	res, err := ss.qry.Get(vehicleID)
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

func (ss *servService) AddService(newService domain.Core) (domain.Core, error) {
	res, err := ss.qry.Add(newService)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.Core{}, errors.New("rejected from database")
		}
		return domain.Core{}, errors.New("some problem on database")
	}

	return res, nil
}

func (ss *servService) Delete(ID uint) error {
	err := ss.qry.Del(ID)
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return errors.New("no data")
		}
	}
	return nil
}
