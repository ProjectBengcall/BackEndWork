package repository

import (
	"bengcall/features/vehicle/domain"

	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	Name_vehicle string
}

type Service struct {
	gorm.Model
	ServiceName string
	Price       int
	VehicleID   uint
}

func FromDomain(dv domain.VehicleCore) Vehicle {
	return Vehicle{
		Model:        gorm.Model{ID: dv.ID},
		Name_vehicle: dv.Name_vehicle,
	}
}

func ToDomain(v Vehicle) domain.VehicleCore {
	return domain.VehicleCore{
		ID:           v.ID,
		Name_vehicle: v.Name_vehicle,
	}
}

func ToDomainArray(av []Vehicle) []domain.VehicleCore {
	var res []domain.VehicleCore
	for _, val := range av {
		res = append(res, domain.VehicleCore{ID: val.ID, Name_vehicle: val.Name_vehicle})
	}

	return res
}

func ToDomainArraySer(av []Service) []domain.ServiceVehicle {
	var res []domain.ServiceVehicle
	for _, val := range av {
		res = append(res, domain.ServiceVehicle{ID: val.ID, ServiceName: val.ServiceName, Price: val.Price, VehicleID: val.VehicleID})
	}

	return res
}
