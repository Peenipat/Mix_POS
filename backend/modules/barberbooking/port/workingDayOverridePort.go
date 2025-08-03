package barberBookingPort

import (
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type WorkingDayOverrideInput struct {
	BranchID  uint      `json:"branch_id" validate:"required"`
	WorkDate  string    `json:"work_date" validate:"required"`    
	StartTime string    `json:"start_time" validate:"required"`    
	EndTime   string    `json:"end_time" validate:"required"`  
	IsClosed  bool		`json:"is_closed"` 
	Reason    string 	`json:"reason"`   
}

type WorkingDayOverrideFilter struct {
	BranchID *uint     `form:"branch_id"` // optional
	FromDate *string   `form:"from_date"` // optional: YYYY-MM-DD
	ToDate   *string   `form:"to_date"`   // optional: YYYY-MM-DD
}

type IWorkingDayOverrideService interface {
	// ดึง override ทั้งหมดในระบบ
	GetAll(ctx context.Context) ([]barberBookingModels.WorkingDayOverride, error)

	// ดึง override ตาม ID
	GetByID(ctx context.Context, id uint) (*barberBookingModels.WorkingDayOverride, error)

	// ดึง override โดยใช้ filter เช่น branch_id, วันที่
	GetByFilter(ctx context.Context, filter WorkingDayOverrideFilter) ([]barberBookingModels.WorkingDayOverride, error)

	// สร้าง override ใหม่
	Create(ctx context.Context, input WorkingDayOverrideInput) (*barberBookingModels.WorkingDayOverride, error)

	// แก้ไข override
	Update(ctx context.Context, id uint, input WorkingDayOverrideInput) error

	// ลบ override
	Delete(ctx context.Context, id uint) error
	GetOverridesByDateRange(ctx context.Context, branchID uint, startDate, endDate time.Time) ([]barberBookingModels.WorkingDayOverride, error)
}


