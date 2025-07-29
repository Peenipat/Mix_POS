package barberBookingPort

import (
	"context"
	"mime/multipart"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type CreateServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
}

type UpdateServiceRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Duration    int     `json:"duration"`
	Price       float64 `json:"price"`
}

type IServiceService interface {
	// Public APIs
	GetAllServices(tenantID uint, branchID uint) ([]barberBookingModels.Service, error)
	GetServiceByID(id uint) (*barberBookingModels.Service, error)

	// Protected APIs
	CreateService(
		ctx context.Context,
		tenantID uint,
		branchID uint,
		payload *CreateServiceRequest,
		file *multipart.FileHeader, 
	) (*barberBookingModels.Service, error)
	UpdateService(
		ctx context.Context,
		serviceID uint,
		payload *UpdateServiceRequest,
		file *multipart.FileHeader, 
	) (*barberBookingModels.Service, error)
	DeleteService(id uint) error
}
