package services

import (
	"bengcall/features/vehicle/domain"
	"errors"
	"strings"

	"github.com/labstack/gommon/log"
)

type vehicleService struct {
	qry domain.Repository
}

func New(repo domain.Repository) domain.Service {
	return &vehicleService{
		qry: repo,
	}
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
	if err != nil {
		if strings.Contains(err.Error(), "table") {
			return errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return errors.New("no data")
		}
	}
	return nil
}

// GetVehicle implements domain.Service
func (vs *vehicleService) GetVehicle() ([]domain.VehicleCore, error) {
	res, err := vs.qry.GetAll()
	if err != nil {
		log.Error(err.Error())
		if strings.Contains(err.Error(), "table") {
			return nil, errors.New("database error")
		} else if strings.Contains(err.Error(), "found") {
			return nil, errors.New("no data")
		}
	}

	if len(res) == 0 {
		log.Info("no data")
		return nil, errors.New("no data")
	}

	return res, nil
}
