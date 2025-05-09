package barberBookingPort
import (
	"context"
	barberBookingModels "myapp/modules/barberbooking/models"
)
type ICustomer interface {
	GetAllCustomers(ctx context.Context, tenantID uint) ([]barberBookingModels.Customer, error)
	GetCustomerByID(ctx context.Context, tenantID, customerID uint) (*barberBookingModels.Customer, error)
	CreateCustomer(ctx context.Context, customer *barberBookingModels.Customer) error
	UpdateCustomer(ctx context.Context, tenantID, customerID uint, updateData *barberBookingModels.Customer) (*barberBookingModels.Customer, error)
	DeleteCustomer(ctx context.Context, tenantID, customerID uint) error
	FindCustomerByEmail(ctx context.Context, tenantID uint, email string) (*barberBookingModels.Customer, error)
}