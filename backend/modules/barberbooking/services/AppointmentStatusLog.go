package barberBookingService

import (
	"context"
	"gorm.io/gorm"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"time"
	"fmt"
)

type appointmentStatusLogService struct {
	DB *gorm.DB
}


func NewAppointmentStatusLogService(db *gorm.DB) barberBookingPort.IAppointmentStatusLogService {
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

func (s *appointmentStatusLogService) GetLogsForAppointment(
    ctx context.Context,
    appointmentID uint,
) ([]barberBookingModels.AppointmentStatusLog, error) {
    // 0. appointmentID ต้องไม่ใช่ 0
    if appointmentID == 0 {
        return nil, fmt.Errorf("invalid appointment id: %d", appointmentID)
    }

    // 1. เช็คว่ามีนัดจริงไหม
    var count int64
    if err := s.DB.WithContext(ctx).
        Model(&barberBookingModels.Appointment{}).
        Where("id = ? AND deleted_at IS NULL", appointmentID).
        Count(&count).Error; err != nil {
        return nil, fmt.Errorf("failed checking appointment existence: %w", err)
    }
    // ถ้าไม่เจอ คืน slice ว่าง ไม่ถือเป็น error
    if count == 0 {
        return []barberBookingModels.AppointmentStatusLog{}, nil
    }

    // 2. ดึง logs ตามลำดับเวลา
    var logs []barberBookingModels.AppointmentStatusLog
    if err := s.DB.WithContext(ctx).
        Where("appointment_id = ?", appointmentID).
        Order("changed_at ASC").
        Find(&logs).Error; err != nil {
        return nil, fmt.Errorf("failed retrieving logs: %w", err)
    }

    // 3. กรองกรณี changed_at เป็น zero
    sanitized := logs[:0]
    for _, lg := range logs {
        if lg.ChangedAt.IsZero() {
            continue
        }
        sanitized = append(sanitized, lg)
    }
    return sanitized, nil
}



func (s *appointmentStatusLogService) DeleteLogsByAppointmentID(ctx context.Context, appointmentID uint) error {
	return s.DB.WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		Delete(&barberBookingModels.AppointmentStatusLog{}).Error
}
