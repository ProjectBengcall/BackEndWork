package delivery

import "bengcall/features/vehicle/domain"

type AddFormat struct {
	Name string `json:"name" form:"name"`
}

func ToDomain(i interface{}) domain.VehicleCore {
	switch i.(type) {
	case AddFormat:
		cnv := i.(AddFormat)
		return domain.VehicleCore{Name: cnv.Name}
	}
	return domain.VehicleCore{}
}
