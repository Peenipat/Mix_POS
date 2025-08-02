package barberBookingPort

import (
	"context"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	"time"
)

type IAppointment interface {
	CheckBarberAvailability(ctx context.Context, tenantID, barberID uint, start, end time.Time) (bool, error)
	CreateAppointment(ctx context.Context, input *barberBookingModels.Appointment) (*barberBookingDto.AppointmentResponseDTO, error)
	GetAvailableBarbers(ctx context.Context, tenantID, branchID uint, start, end time.Time) ([]barberBookingModels.Barber, error)
	UpdateAppointment(ctx context.Context, id uint, tenantID uint, input *barberBookingModels.Appointment) (*barberBookingModels.Appointment, error)
	GetAppointmentByID(ctx context.Context, id uint) (*barberBookingModels.Appointment, error)
	ListAppointments(ctx context.Context, filter barberBookingDto.AppointmentFilter) ([]barberBookingModels.Appointment, error)
	CancelAppointment(ctx context.Context, appointmentID uint, actorUserID *uint, actorCustomerID *uint) error
	RescheduleAppointment(ctx context.Context, appointmentID uint, newStartTime time.Time, actorUserID *uint, actorCustomerID *uint) error
	CalculateAppointmentEndTime(ctx context.Context, serviceID uint, startTime time.Time) (time.Time, error)
	DeleteAppointment(ctx context.Context, appointmentID uint) error
	GetUpcomingAppointmentsByCustomer(ctx context.Context, customerID uint) (*barberBookingModels.Appointment, error)
	GetAppointments(
		ctx context.Context,
		filter GetAppointmentsFilter,
	) ([]AppointmentBrief, int64, error)
	GetAppointmentsByBarber(ctx context.Context, barberID uint, filter AppointmentFilter) ([]AppointmentBrief, error)
	GetAppointmentsByPhone(
		ctx context.Context,
		phone string,
	) ([]AppointmentBrief, error)
}

// CreateAppointmentRequest is the payload for creating a new appointment
// swagger:model CreateAppointmentRequest
// example:
//
//	{
//	  "branch_id": 1,
//	  "service_id": 2,
//	  "barber_id": 3,
//	  "customer_id": 4,
//	  "start_time": "2025-05-30T10:00:00Z",
//	  "notes": "Please be on time"
//	}
type CustomerInput struct {
	Name  string `json:"name" example:"John"`
	Phone string `json:"phone" example:"0123456789"`
}

type CreateAppointmentRequest struct {
	BranchID   uint           `json:"branch_id" example:"1"`
	ServiceID  uint           `json:"service_id" example:"2"`
	BarberID   *uint          `json:"barber_id,omitempty" example:"3"`
	CustomerID uint           `json:"customer_id" example:"4"`
	StartTime  string         `json:"start_time" example:"2025-05-30T10:00:00Z"`
	Notes      string         `json:"notes,omitempty" example:"Preferred barber: John"`
	Customer   *CustomerInput `json:"customer,omitempty"`
}

type ServiceBrief struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Duration    int    `json:"duration"`
	Price       int    `json:"price"`
}

type BarberBrief struct {
	Username string `json:"username"`
}

type CustomerBrief struct {
	Name  string `json:"name"`
	Phone string `json:"phone"`
}

type AppointmentFilter struct {
	Start    *time.Time
	End      *time.Time
	Status   []barberBookingModels.AppointmentStatus
	TimeMode string
}

type GetAppointmentsFilter struct {
	TenantID  *uint
	BranchID  *uint
	Search    string
	Statuses  []barberBookingModels.AppointmentStatus
	BarberID  *uint
	ServiceID *uint
	StartDate *time.Time
	EndDate   *time.Time
	Page      int
	Limit     int
}

type AppointmentBrief struct {
	ID         uint          `json:"id"`
	BranchID   uint          `json:"branch_id"`
	ServiceID  uint          `json:"service_id"`
	Service    ServiceBrief  `json:"service"`
	BarberID   uint          `json:"barber_id"`
	Barber     BarberBrief   `json:"barber"`
	CustomerID uint          `json:"customer_id"`
	Customer   CustomerBrief `json:"customer"`
	StartTime  time.Time     `json:"start_time"`
	EndTime    time.Time     `json:"end_time"`
	Status     string        `json:"status"`
}
