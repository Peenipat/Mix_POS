package barberBookingService
import (
	"gorm.io/gorm"
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"

)



type appointmentStatusLogService struct {
	DB *gorm.DB
}


func NewAppointmentStatusLogService(db *gorm.DB) IAppointmentStatusLogService {
	return &appointmentStatusLogService{DB: db}
}

func (s *appointmentStatusLogService) LogStatusChange(ctx context.Context, appointmentID uint, oldStatus, newStatus string, userID *uint, customerID *uint, notes string) error {
	log := barberBookingModels.AppointmentStatusLog{
		AppointmentID:       appointmentID,
		OldStatus:           oldStatus,
		NewStatus:           newStatus,
		ChangedByUserID:     userID,
		ChangedByCustomerID: customerID,
		ChangedAt:           time.Now().UTC(),
		Notes:               notes,
	}
	return s.DB.WithContext(ctx).Create(&log).Error
}

func (s *appointmentStatusLogService) GetLogsForAppointment(ctx context.Context, appointmentID uint) ([]barberBookingModels.AppointmentStatusLog, error) {
	var logs []barberBookingModels.AppointmentStatusLog
	err := s.DB.WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		Order("changed_at ASC").
		Find(&logs).Error
	return logs, err
}

func (s *appointmentStatusLogService) DeleteLogsByAppointmentID(ctx context.Context, appointmentID uint) error {
	return s.DB.WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		Delete(&barberBookingModels.AppointmentStatusLog{}).Error
}
