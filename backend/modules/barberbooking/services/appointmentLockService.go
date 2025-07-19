package barberBookingService

import (
	"context"
	"gorm.io/gorm"
	"errors"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"time"
	"fmt"
	"log"
)

type appointmentLockService struct {
	db *gorm.DB
}

func NewAppointmentLockService(db *gorm.DB) barberBookingPort.IAppointmentLock {
	return &appointmentLockService{db}
}

// CleanupExpiredLocks implements barberBookingPort.IAppointmentLock.
func (a *appointmentLockService) CleanupExpiredLocks(ctx context.Context) error {
	panic("unimplemented")
}

// CreateAppointmentLock implements barberBookingPort.IAppointmentLock.
func (a *appointmentLockService) CreateAppointmentLock(
	ctx context.Context,
	input barberBookingPort.AppointmentLockInput,
) (*barberBookingModels.AppointmentLock, error) {

	db := a.db.WithContext(ctx)
	var count int64

	if err := db.Model(&barberBookingModels.Appointment{}).
		Where("tenant_id = ? AND branch_id = ? AND barber_id = ? AND start_time < ? AND end_time > ? AND status IN ?",
			input.TenantID, input.BranchID, input.BarberID, input.EndTime, input.StartTime,
			[]string{"PENDING", "CONFIRMED"}).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("slot is already booked")
	}

	if err := db.Model(&barberBookingModels.AppointmentLock{}).
		Where("tenant_id = ? AND branch_id = ? AND barber_id = ? AND start_time < ? AND end_time > ? AND is_active = ? AND expires_at > ?",
			input.TenantID, input.BranchID, input.BarberID, input.EndTime, input.StartTime,
			true, time.Now()).
		Count(&count).Error; err != nil {
		return nil, err
	}
	if count > 0 {
		return nil, errors.New("slot is currently being locked by another user")
	}
	lock := &barberBookingModels.AppointmentLock{
		TenantID:   input.TenantID,
		BranchID:   input.BranchID,
		BarberID:   input.BarberID,
		CustomerID: input.CustomerID,
		StartTime:  input.StartTime,
		EndTime:    input.EndTime,
		ExpiresAt:  time.Now().Add(7 * time.Minute),
		IsActive:   true,
	}

	if err := db.Create(lock).Error; err != nil {
		return nil, err
	}

	return lock, nil
}

// GetAppointmentLocks implements barberBookingPort.IAppointmentLock.
func (a *appointmentLockService) GetAppointmentLocks(
	ctx context.Context,
	branchID uint,
	barberID uint,
	date time.Time,
) ([]barberBookingModels.AppointmentLock, error) {

	db := a.db.WithContext(ctx)

	// เริ่มต้นและสิ้นสุดของวัน (00:00 - 23:59)
	startOfDay := time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())
	endOfDay := startOfDay.Add(24 * time.Hour).Add(-time.Nanosecond)
	log.Println(": ",branchID,"server time now: ", barberID)
	var locks []barberBookingModels.AppointmentLock

	err := db.Model(&barberBookingModels.AppointmentLock{}).
		Where("branch_id = ? AND barber_id = ? AND start_time BETWEEN ? AND ? AND is_active = ? AND expires_at > ?", 
			branchID, barberID, startOfDay, endOfDay, true, time.Now()).
		Order("start_time ASC").
		Find(&locks).Error

	if err != nil {
		return nil, err
	}

	return locks, nil
}


// IsSlotAvailable implements barberBookingPort.IAppointmentLock.
func (a *appointmentLockService) IsSlotAvailable(ctx context.Context, tenantID uint, branchID uint, barberID uint, start time.Time, end time.Time) (bool, error) {
	panic("unimplemented")
}

func (a *appointmentLockService) ReleaseAppointmentLock(ctx context.Context, lockID uint) error {
	db := a.db.WithContext(ctx)

	// ตรวจสอบว่า lock นี้มีอยู่และยัง active อยู่ไหม
	var lock barberBookingModels.AppointmentLock
	if err := db.First(&lock, "id = ? AND is_active = ?", lockID, true).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("lock not found or already released")
		}
		return err
	}

	// ปรับ is_active = false เพื่อปลดล็อก
	if err := db.Model(&lock).Update("is_active", false).Error; err != nil {
		return err
	}

	return nil
}

