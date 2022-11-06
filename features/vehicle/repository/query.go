package repository

import (
	"bengcall/features/vehicle/domain"

	"github.com/labstack/gommon/log"
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

// Add implements domain.Repository
func (rq *repoQuery) Add(newItem domain.VehicleCore) (domain.VehicleCore, error) {
	var cnv Vehicle
	cnv = FromDomain(newItem)
	if err := rq.db.Create(&cnv).Error; err != nil {
		log.Error("Error on insert user", err.Error())
		return domain.VehicleCore{}, err
	}
	// selesai dari DB
	newItem = ToDomain(cnv)
	return newItem, nil
}

// Delete implements domain.Repository
func (rq *repoQuery) Delete(vehicleID uint) error {
	var data Vehicle
	if err := rq.db.Delete(&data, "id = ?", vehicleID).Error; err != nil {
		log.Error("error on deleting user", err.Error())
		return err
	}
	return nil
}

// GetAll implements domain.Repository
func (rq *repoQuery) GetAll() ([]domain.VehicleCore, error) {
	var resQry []Vehicle
	if err := rq.db.Find(&resQry).Error; err != nil {
		log.Error("Error on All user", err.Error())
		return nil, err
	}
	// selesai dari DB
	res := ToDomainArray(resQry)
	return res, nil
}