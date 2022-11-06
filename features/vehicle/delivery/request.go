package delivery

import "bengcall/features/vehicle/domain"

type AddFormat struct {
	Name_vehicle string `json:"name_vehicle" form:"name_vehicle"`
}

func ToDomain(i interface{}) domain.VehicleCore {
	switch i.(type) {
	case AddFormat:
		cnv := i.(AddFormat)
		return domain.VehicleCore{Name_vehicle: cnv.Name_vehicle}
	}
	return domain.VehicleCore{}
}
