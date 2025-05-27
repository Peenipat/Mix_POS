package barberBookingPort

import(
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type BarberWithUser struct {
    ID        uint      `json:"id"`
    BranchID  uint      `json:"branch_id"`
    UserID    uint      `json:"user_id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}

type IBarber interface{
	CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error
	GetBarberByID(ctx context.Context, id uint) (*barberBookingModels.Barber, error) 
	ListBarbersByBranch(ctx context.Context, branchID *uint) ([]BarberWithUser, error)
	UpdateBarber(ctx context.Context, id uint, updated *barberBookingModels.Barber) (*barberBookingModels.Barber, error)
	DeleteBarber(ctx context.Context, id uint) error
	GetBarberByUser(ctx context.Context, userID uint)(*barberBookingModels.Barber, error) 
	ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) 
}