package barberBookingDto

import (
	barberBookingModels "myapp/modules/barberbooking/models"
	"time"
)

type AppointmentFilter struct {
	TenantID   uint
	BranchID   *uint
	BarberID   *uint
	CustomerID *uint
	Status     *barberBookingModels.AppointmentStatus
	StartDate  *time.Time
	EndDate    *time.Time
	Limit      *int
	Offset     *int
	SortBy     *string
}
