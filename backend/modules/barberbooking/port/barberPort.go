package barberBookingPort

import(
	"context"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IBarber interface{
	CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error
	GetBarberByID(ctx context.Context, id uint) (*barberBookingModels.Barber, error) 
	ListBarbersByBranch(ctx context.Context, branchID *uint) ([]barberBookingModels.Barber, error)
	UpdateBarber(ctx context.Context, id uint, updated *barberBookingModels.Barber) (*barberBookingModels.Barber, error)
	DeleteBarber(ctx context.Context, id uint) error
	GetBarberByUser(ctx context.Context, userID uint)(*barberBookingModels.Barber, error) 
	ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) 
}