package barberBookingPort
import (
	"context"
	barberBookingModels "myapp/modules/barberbooking/models"
)
type ICustomer interface {
	GetCustomers(ctx context.Context, filter GetCustomersFilter) ([]barberBookingModels.Customer, int64, error)
	GetCustomerByID(ctx context.Context, tenantID, customerID uint) (*barberBookingModels.Customer, error)
	CreateCustomer(ctx context.Context, customer *barberBookingModels.Customer) error
	UpdateCustomer(ctx context.Context, tenantID, customerID uint, updateData *barberBookingModels.Customer) (*barberBookingModels.Customer, error)
	DeleteCustomer(ctx context.Context, tenantID, customerID uint) error
	FindCustomerByEmail(ctx context.Context, tenantID uint, email string) (*barberBookingModels.Customer, error)
	GetPendingAndCancelledCount( ctx context.Context,tenantID uint,branchID uint,customerID uint) ([]CountByCustomerStatus, error)
}
type CountByCustomerStatus struct {
    CustomerID uint                  `gorm:"column:customer_id"`
    Status     barberBookingModels.AppointmentStatus `gorm:"column:status"`
    Total      int64                 `gorm:"column:total"`
}

type GetCustomersFilter struct {
	TenantID  uint
	BranchID  uint
	Name      string
	Phone     string
	SortBy    string // created_at / updated_at
	SortOrder string // asc / desc
	Page      int
	Limit     int
}
