package services

import (
	"bengcall/config"
	"bengcall/features/vehicle/domain"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type vehicleService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &vehicleService{
		qry: repo,
	}
}

// GetService implements domain.Service
func (vs *vehicleService) GetService() ([]domain.VehicleCore, []domain.ServiceVehicle, error) {
	resVehicle, err := vs.qry.GetAll()
	if err != nil {
		return nil, nil, errors.New("no data")
	}
	resService, err := vs.qry.Get()
	if err != nil {
		return nil, nil, errors.New("get service error")
	}

	return resVehicle, resService, nil
}

// AddVehicle implements domain.Service
func (vs *vehicleService) AddVehicle(newItem domain.VehicleCore) (domain.VehicleCore, error) {
	res, err := vs.qry.Add(newItem)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate") {
			return domain.VehicleCore{}, errors.New("rejected from database")
		}

		return domain.VehicleCore{}, errors.New("some problem on database")
	}

	return res, nil
}

// DeleteVehicle implements domain.Service
func (vs *vehicleService) DeleteVehicle(vehicleID uint) error {
	err := vs.qry.Delete(vehicleID)
	if err == gorm.ErrRecordNotFound {
		log.Error(err.Error())
		return gorm.ErrRecordNotFound
	} else if err != nil {
		log.Error(err.Error())
		return errors.New(config.DATABASE_ERROR)
	}
	return nil
}

// GetVehicle implements domain.Service
func (vs *vehicleService) GetVehicle() ([]domain.VehicleCore, error) {
	res, err := vs.qry.GetAll()
	if err == gorm.ErrRecordNotFound {
		log.Error(err.Error())
		return nil, gorm.ErrRecordNotFound
	} else if err != nil {
		log.Error(err.Error())
		return nil, errors.New(config.DATABASE_ERROR)
	}

	if len(res) == 0 {
		log.Info("no data")
		return nil, errors.New(config.DATA_NOTFOUND)
	}

	return res, nil
}
