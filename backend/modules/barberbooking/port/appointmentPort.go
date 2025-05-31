package barberBookingPort
import (
	"time"
	"context"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingDto "myapp/modules/barberbooking/dto"
)

type IAppointment interface {
	CheckBarberAvailability(ctx context.Context, tenantID, barberID uint, start, end time.Time) (bool, error)
	CreateAppointment(ctx context.Context, input *barberBookingModels.Appointment) (*barberBookingModels.Appointment, error)
	GetAvailableBarbers(ctx context.Context, tenantID, branchID uint, start, end time.Time) ([]barberBookingModels.Barber, error)
	UpdateAppointment(ctx context.Context, id uint, tenantID uint, input *barberBookingModels.Appointment) (*barberBookingModels.Appointment, error)
	GetAppointmentByID(ctx context.Context, id uint) (*barberBookingModels.Appointment, error)
	ListAppointments(ctx context.Context, filter barberBookingDto.AppointmentFilter) ([]barberBookingModels.Appointment, error)
	CancelAppointment(ctx context.Context,appointmentID uint,actorUserID *uint,actorCustomerID *uint,) error
	RescheduleAppointment( ctx context.Context,appointmentID uint,newStartTime time.Time,actorUserID *uint, actorCustomerID *uint,) error
	CalculateAppointmentEndTime(ctx context.Context, serviceID uint, startTime time.Time) (time.Time, error)
	GetAppointmentsByBarber(ctx context.Context, barberID uint, start *time.Time, end *time.Time) ([]barberBookingModels.Appointment, error)
	DeleteAppointment(ctx context.Context, appointmentID uint) error
	GetUpcomingAppointmentsByCustomer(ctx context.Context, customerID uint) (*barberBookingModels.Appointment, error)
}

// CreateAppointmentRequest is the payload for creating a new appointment
// swagger:model CreateAppointmentRequest
// example:
// {
//   "branch_id": 1,
//   "service_id": 2,
//   "barber_id": 3,
//   "customer_id": 4,
//   "start_time": "2025-05-30T10:00:00Z",
//   "notes": "Please be on time"
// }
type CreateAppointmentRequest struct {
    BranchID   uint   `json:"branch_id" example:"1"`
    ServiceID  uint   `json:"service_id" example:"2"`
    BarberID   *uint  `json:"barber_id,omitempty" example:"3"`
    CustomerID uint   `json:"customer_id" example:"4"`
    StartTime  string `json:"start_time" example:"2025-05-30T10:00:00Z"`
    Notes      string `json:"notes,omitempty" example:"Preferred barber: John"`
}