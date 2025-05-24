package barberBookingPort
import (
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IServiceService interface {
	// Public APIs
	GetAllServices() ([]barberBookingModels.Service, error)
	GetServiceByID(id uint) (*barberBookingModels.Service, error)

	// Protected APIs
	CreateService(service *barberBookingModels.Service) error
	UpdateService(id uint, service *barberBookingModels.Service) (*barberBookingModels.Service, error)
	DeleteService(id uint) error
}
