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

type UserNotBarber struct {
    UserID      uint      `json:"user_id"`
    Username    string    `json:"username"`
    Email       string    `json:"email"`
    PhoneNumber string    `json:"phone_number"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type UpdateBarberRequest struct {
    BranchID    uint   `json:"branch_id"`
    UserID      uint   `json:"user_id"`
    PhoneNumber string `json:"phone_number"`
    Username    string `json:"username"` // ถ้าต้องการอัปเดตชื่อผู้ใช้ด้วย
    Email       string `json:"email"`    // ถ้าต้องการอัปเดตอีเมลด้วย
}

type IBarber interface{
	CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error
	GetBarberByID(ctx context.Context, id uint) (*barberBookingModels.Barber, error) 
	ListBarbersByBranch(ctx context.Context, branchID *uint) ([]BarberWithUser, error)
    UpdateBarber(ctx context.Context,barberID uint,updated *barberBookingModels.Barber,updatedUsername string,updatedEmail string,) (*barberBookingModels.Barber, error)
	DeleteBarber(ctx context.Context, id uint) error
	GetBarberByUser(ctx context.Context, userID uint)(*barberBookingModels.Barber, error) 
	ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) 
	ListUserNotBarber(ctx context.Context, branchID *uint) ([]UserNotBarber, error)
}

