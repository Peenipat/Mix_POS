package barberBookingPort
import (
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingDto "myapp/modules/barberbooking/dto"
)

type IbarberWorkload interface{
	GetWorkloadByBarber(ctx context.Context, barberID uint, date time.Time) (*barberBookingModels.BarberWorkload, error)
	GetWorkloadSummaryByBranch(ctx context.Context,date time.Time,tenantID uint,branchID uint,) ([]barberBookingDto.BranchWorkloadSummary, error)
	UpsertBarberWorkload(ctx context.Context, barberID uint, date time.Time, appointments int, hours int) error 
}

type UpsertBarberWorkloadRequest struct {
    // วันที่ในรูปแบบ YYYY-MM-DD
    Date         string `json:"date" example:"2025-05-30"`
    // จำนวนการนัดหมาย
    Appointments int    `json:"appointments" example:"10"`
    // จำนวนชั่วโมงทำงาน
    Hours        int    `json:"hours" example:"8"`
}
