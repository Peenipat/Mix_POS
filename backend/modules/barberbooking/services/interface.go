package barberBookingService

import (
	barberBookingModels "myapp/modules/barberbooking/models"
	"context"
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

type ICustomerService interface {
	GetAllCustomers(ctx context.Context, tenantID uint) ([]barberBookingModels.Customer, error)
	GetCustomerByID(ctx context.Context, tenantID, customerID uint) (*barberBookingModels.Customer, error)
	CreateCustomer(ctx context.Context, customer *barberBookingModels.Customer) error
	UpdateCustomer(ctx context.Context, tenantID, customerID uint, updateData map[string]interface{}) error
	DeleteCustomer(ctx context.Context, tenantID, customerID uint) error
	FindCustomerByEmail(ctx context.Context, tenantID uint, email string) (*barberBookingModels.Customer, error)
}

type IAppointmentStatusLogService interface {
	LogStatusChange(ctx context.Context, appointmentID uint, oldStatus, newStatus string, userID *uint, customerID *uint, notes string) error
	GetLogsForAppointment(ctx context.Context, appointmentID uint) ([]barberBookingModels.AppointmentStatusLog, error)
	DeleteLogsByAppointmentID(ctx context.Context, appointmentID uint) error
}
