package delivery

import "bengcall/features/service/domain"

type ServiceFormat struct {
	ServiceName string `json:"service_name" form:"service_name" validate:"required"`
	Price       int    `json:"price" form:"price" validate:"required"`
	VehicleID   uint   `json:"vehicle_id" form:"vehicle_id" validate:"required"`
}

func ToDomain(i interface{}) domain.Core {
	switch i.(type) {
	case ServiceFormat:
		cnv := i.(ServiceFormat)
		return domain.Core{ServiceName: cnv.ServiceName, Price: cnv.Price, VehicleID: cnv.VehicleID}
	}
	return domain.Core{}
}
