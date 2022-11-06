package repository

import (
	"bengcall/features/service/domain"
	"time"

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

func (rq *repoQuery) Get(vehicleID int) ([]domain.Core, error) {
	var resQry []Service
	if err := rq.db.Table("services").Select("id", "service_name", "price", "vehicle_id").Where("vehicle_id = ?", vehicleID).Model(&Service{}).Find(&resQry).Error; err != nil {
		return nil, err
	}
	res := ToDomainArray(resQry)
	return res, nil
}

func (rq *repoQuery) Add(newService domain.Core) (domain.Core, error) {
	var resQry Service
	if err := rq.db.Exec("INSERT INTO services (id, created_at, updated_at, deleted_at, service_name, price, vehicle_id) values (?,?,?,?,?,?,?)",
		nil, time.Now(), time.Now(), nil, newService.ServiceName, newService.Price, newService.VehicleID).Error; err != nil {
		return domain.Core{}, err
	}
	if er := rq.db.Table("services").Select("id", "service_name", "price", "vehicle_id").Where("service_name = ? && vehicle_id = ?", newService.ServiceName, newService.VehicleID).Model(&Service{}).Find(&resQry).Error; er != nil {
		return domain.Core{}, er
	}
	res := ToDomain(resQry)
	return res, nil
}

func (rq *repoQuery) Del(ID uint) error {
	var resQry Service
	if err := rq.db.Where("id = ?", ID).Delete(&resQry).Error; err != nil {
		return err
	}
	return nil
}
