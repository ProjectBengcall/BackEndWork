package delivery

import "bengcall/features/service/domain"

func SuccessResponse(msg string, data interface{}) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
		"data":    data,
	}
}

func FailResponse(msg string) map[string]interface{} {
	return map[string]interface{}{
		"message": msg,
	}
}

type Response struct {
	ID          uint   `json:"id"`
	ServiceName string `json:"service_name"`
	Price       int    `json:"price"`
	VehicleID   uint   `json:"vehicle_id"`
}

func ToResponse(core interface{}, code string) interface{} {
	var res interface{}
	switch code {
	case "service":
		cnv := core.(domain.Core)
		res = Response{ID: cnv.ID, ServiceName: cnv.ServiceName, Price: cnv.Price, VehicleID: cnv.VehicleID}
	case "servray":
		var arr []Response
		cnv := core.([]domain.Core)
		for _, val := range cnv {
			arr = append(arr, Response{ID: val.ID, ServiceName: val.ServiceName, Price: val.Price, VehicleID: val.VehicleID})
		}
		res = arr
	}
	return res
}
