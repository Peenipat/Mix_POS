package barberBookingPort
import (
	"context"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IAppointmentStatusLogService interface {
	LogStatusChange(ctx context.Context, appointmentID uint, oldStatus, newStatus string, userID *uint, customerID *uint, notes string) error
	GetLogsForAppointment(ctx context.Context, appointmentID uint) ([]barberBookingModels.AppointmentStatusLog, error)
	DeleteLogsByAppointmentID(ctx context.Context, appointmentID uint) error
}

