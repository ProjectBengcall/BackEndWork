package domain

type VehicleCore struct {
	ID           uint
	Name_vehicle string
}

type ServiceVehicle struct {
	ID           uint
	ServiceName  string
	Price        int
	Name_vehicle string
}

type ServiceVehicleDet struct {
	ServiceVehicleDet []ServiceVehicle
}

type Repository interface {
	Delete(vehicleID uint) error
	Add(newItem VehicleCore) (VehicleCore, error)
	GetAll() ([]VehicleCore, error)

	Get() ([]ServiceVehicle, error)
}

type Service interface {
	DeleteVehicle(vehicleID uint) error
	AddVehicle(newItem VehicleCore) (VehicleCore, error)
	GetVehicle() ([]VehicleCore, error)

	GetService() ([]ServiceVehicle, error)
}
