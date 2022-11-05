package domain

type VehicleCore struct {
	ID   uint
	Name string
}

type Repository interface {
	Delete(vehicleID uint) error
	Add(newItem VehicleCore) (VehicleCore, error)
	GetAll() ([]VehicleCore, error)
}

type Service interface {
	DeleteVehicle(vehicleID uint) error
	AddVehicle(newItem VehicleCore) (VehicleCore, error)
	GetVehicle() ([]VehicleCore, error)
}
