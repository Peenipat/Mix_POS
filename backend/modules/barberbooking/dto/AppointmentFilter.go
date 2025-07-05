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


type AppointmentResponseDTO struct {
	ID         uint      `json:"id"`
	TenantID   uint      `json:"tenant_id"`
	BranchID   uint      `json:"branch_id"`
	ServiceID  uint      `json:"service_id"`
	BarberID   uint      `json:"barber_id"`
	CustomerID uint      `json:"customer_id"`
	StartTime  time.Time `json:"start_time"`
	EndTime    time.Time `json:"end_time"`
	Status     string    `json:"status"`
	Notes      string    `json:"notes,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

