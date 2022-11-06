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
