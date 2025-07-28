package barberBookingPort

import (
	"context"
	"mime/multipart"
	barberBookingModels "myapp/modules/barberbooking/models"
	"time"
)

type BarberWithUser struct {
    ID        uint      `json:"id"`
    BranchID  uint      `json:"branch_id"`
    UserID    uint      `json:"user_id"`
    Username  string    `json:"username"`
    Email     string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
    Description string    `json:"description"`
    RoleUser    string    `json:"role_user"`
    ImgPath     *string   `json:"img_path"` 
    ImgName     *string   `json:"img_name"`
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
    BranchID     uint   `form:"branch_id"`  
    Description  string `form:"description"`
    RoleUser     string `form:"role_user"`
    Username     string `form:"username"`
    Email        string `form:"email"`
    PhoneNumber  string `form:"phone_number"`
    ImgPath      string `form:"img_path"`
    ImgName      string `form:"img_name"`
}

type CreateBarberInput struct {
	UserID          uint        `json:"user_id"`
    Description     string      `gorm:"type:varchar(100);not null" json:"description"`
}

type BarberDetailResponse struct {
	ID          uint   `json:"id"`
	BranchID    uint   `json:"branch_id"`
	TenantID    uint   `json:"tenant_id"`
	UserID      uint   `json:"user_id"`
	RoleUser    string `json:"role_user"`
	Description string `json:"description"`

	User struct {
		ID          uint   `json:"id"`
		Username    string `json:"username"`
		Email       string `json:"email"`
		PhoneNumber string `json:"phone_number"`
		BranchID    uint   `json:"branch_id"`
		ImgPath     string `json:"Img_path"`
		ImgName     string `json:"Img_name"`
	} `json:"user"`
}

type BarberUserMinimal struct {
	ID          uint   `json:"id"`
	Username    string `json:"username"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
	BranchID    uint   `json:"branch_id"`
	ImgPath     string `json:"Img_path"`
	ImgName     string `json:"Img_name"`
}

type BarberDetailMinimalResponse struct {
	ID          uint              `json:"id"`
	User        BarberUserMinimal `json:"user"`
	TenantID    uint              `json:"tenant_id"`
	RoleUser    string            `json:"role_user"`
	Description string            `json:"description"`
}


type IBarber interface{
	CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error
	GetBarberByID(ctx context.Context, id uint) (*BarberDetailResponse, error)
	ListBarbersByBranch(ctx context.Context, branchID *uint) ([]BarberWithUser, error)
	UpdateBarber(
		ctx context.Context,
		barberID uint,
		payload *UpdateBarberRequest,
		file *multipart.FileHeader,
	) (*barberBookingModels.Barber, error)
	DeleteBarber(ctx context.Context, id uint) error
	GetBarberByUser(ctx context.Context, userID uint)(*barberBookingModels.Barber, error) 
	ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) 
	ListUserNotBarber(ctx context.Context, branchID *uint) ([]UserNotBarber, error)
}

