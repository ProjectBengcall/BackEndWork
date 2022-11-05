package repository

import (
	"bengcall/features/vehicle/domain"

	"gorm.io/gorm"
)

type Vehicle struct {
	gorm.Model
	Name string
}

func FromDomain(dv domain.VehicleCore) Vehicle {
	return Vehicle{
		Model: gorm.Model{ID: dv.ID},
		Name:  dv.Name,
	}
}

func ToDomain(v Vehicle) domain.VehicleCore {
	return domain.VehicleCore{
		ID:   v.ID,
		Name: v.Name,
	}
}

func ToDomainArray(av []Vehicle) []domain.VehicleCore {
	var res []domain.VehicleCore
	for _, val := range av {
		res = append(res, domain.VehicleCore{ID: val.ID, Name: val.Name})
	}

	return res
}
