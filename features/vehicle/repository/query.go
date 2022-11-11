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

// Get implements domain.Repository
func (rq *repoQuery) Get() ([]domain.VehicleCore, []domain.ServiceVehicle, error) {
	var verQry []Vehicle
	var serQry []Service

	if err := rq.db.Find(&verQry).Error; err != nil {
		log.Error("Error on All user", err.Error())
		return nil, nil, err
	}

	if err := rq.db.Find(&serQry).Error; err != nil {
		log.Error("Error on All user", err.Error())
		return nil, nil, err
	}

	// if err := rq.db.Table("vehicles").Select("vehicles.id", "vehicles.name_vehicle", "services.service_name", "services.price").Joins("join services on services.vehicle_id=vehicles.id").Order("vehicles.name_vehicle asc").Where("services.deleted_at IS NULL").Model(&Service{}).Find(&serQry).Error; err != nil {
	// 	return nil, nil, err
	// }

	ver := ToDomainArray(verQry)
	ser := ToDomainArraySer(serQry)
	return ver, ser, nil
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
func (rq *repoQuery) Delete(vehicleID uint) (domain.VehicleCore, error) {
	var resQry Vehicle

	if err := rq.db.First(&resQry, "ID = ?", vehicleID).Error; err != nil {
		log.Error(err.Error())
		return ToDomain(resQry), err
	}
	if err := rq.db.Delete(&resQry).Error; err != nil {
		log.Error(err.Error())
		return ToDomain(resQry), err
	}

	res := ToDomain(resQry)
	return res, nil
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
