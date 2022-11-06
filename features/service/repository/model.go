package repository

import (
	"bengcall/features/service/domain"

	"gorm.io/gorm"
)

type Service struct {
	gorm.Model
	ServiceName string
	Price       int
	VehicleID   int
}

func FromDomain(du domain.Core) Service {
	return Service{
		Model:       gorm.Model{ID: du.ID},
		ServiceName: du.ServiceName,
		Price:       du.Price,
		VehicleID:   du.VehicleID,
	}
}

func ToDomain(s Service) domain.Core {
	return domain.Core{
		ID:          s.ID,
		ServiceName: s.ServiceName,
		Price:       s.Price,
		VehicleID:   s.VehicleID,
	}
}

func ToDomainArray(sa []Service) []domain.Core {
	var res []domain.Core
	for _, val := range sa {
		res = append(res, domain.Core{ID: val.ID, ServiceName: val.ServiceName, Price: val.Price, VehicleID: val.VehicleID})
	}

	return res
}
