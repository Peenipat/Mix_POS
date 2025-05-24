package barberBookingPort

import (
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IUnavailabilitySerivce interface{
	CreateUnavailability(ctx context.Context, input *barberBookingModels.Unavailability) (*barberBookingModels.Unavailability, error)
	GetUnavailabilitiesByBranch(ctx context.Context, branchID uint, from, to time.Time) ([]barberBookingModels.Unavailability, error) 
	GetUnavailabilitiesByBarber(ctx context.Context, barberID uint, from, to time.Time) ([]barberBookingModels.Unavailability, error)
	UpdateUnavailability(ctx context.Context, id uint, updates map[string]interface{}) error
	DeleteUnavailability(ctx context.Context, id uint) error 
}