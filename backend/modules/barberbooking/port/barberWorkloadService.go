package barberBookingPort
import (
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IbarberWorkload interface{
	GetWorkloadByBarber(ctx context.Context, barberID uint, date time.Time) (*barberBookingModels.BarberWorkload, error)
	GetWorkloadByDate(ctx context.Context, date time.Time) ([]barberBookingModels.BarberWorkload, error)
	UpsertBarberWorkload(ctx context.Context, barberID uint, date time.Time, appointments int, hours int) error 
}