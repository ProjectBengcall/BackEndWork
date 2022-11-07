package domain

type Core struct {
	ID          uint
	ServiceName string
	Price       int
	VehicleID   uint
}

type Repository interface {
	Add(newService Core) (Core, error)
	Get(vehicleID int) ([]Core, error)
	Del(ID uint) error
}

type Service interface {
	AddService(newService Core) (Core, error)
	GetSpesific(vehicleID int) ([]Core, error)
	Delete(ID uint) error
}
