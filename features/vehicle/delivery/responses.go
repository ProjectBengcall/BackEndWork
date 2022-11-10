package delivery

import "bengcall/features/vehicle/domain"

func SuccessResponse(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		"data":    data,
	}
}

func FailResponse(msg string) map[string]string {
	return map[string]string{
		"message": msg,
	}
}

type AddResponse struct {
	ID           uint   `json:"id"`
	Name_vehicle string `json:"name_vehicle"`
}

type GetResponse struct {
	ID           uint     `json:"id"`
	Name_vehicle string   `json:"name_vehicle"`
	Services     []GetRes `json:"services"`
}

type GetRes struct {
	ID          uint   `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	VehicleID   uint   `json:"vehicle_id"`
}

func ToResponse(basic interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "reg":
		cnv := basic.(domain.VehicleCore)
		res = AddResponse{ID: cnv.ID, Name_vehicle: cnv.Name_vehicle}
	case "all":
		var arr []AddResponse
		cnv := basic.([]domain.VehicleCore)
		for _, val := range cnv {
			arr = append(arr, AddResponse{ID: val.ID, Name_vehicle: val.Name_vehicle})
		}
		res = arr
	}

	return res
}

func ToResponseGet(vehicle interface{}, service interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "vs":
		var arr []GetResponse
		cnv := vehicle.([]domain.VehicleCore)
		for _, val := range cnv {

			var ar []GetRes
			cnr := service.([]domain.ServiceVehicle)
			for _, vol := range cnr {
				if vol.VehicleID == val.ID {
					ar = append(ar, GetRes{ID: vol.ID, ServiceName: vol.ServiceName, Price: vol.Price})
				}
			}

			arr = append(arr, GetResponse{ID: val.ID, Name_vehicle: val.Name_vehicle, Services: ar})

		}
		res = arr
	}
	return res
}
