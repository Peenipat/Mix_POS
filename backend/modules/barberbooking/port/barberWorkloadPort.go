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