package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type appointmentService struct {
	DB *gorm.DB
}



func NewAppointmentService(db *gorm.DB) *appointmentService {
	return &appointmentService{DB: db}
}

// checkBarberAvailabilityTx...
func (s *appointmentService) checkBarberAvailabilityTx(tx *gorm.DB, tenantID, barberID uint, start, end time.Time) (bool, error) {
	var barber barberBookingModels.Barber
	if err := tx.Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", barberID, tenantID).
		First(&barber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // ❗️ไม่เจอ barber = ไม่ว่าง
		}
		return false, err // error อื่น ๆ
	}

	// ตรวจสอบว่า barber มีการจองคิวช่วงเวลานี้หรือไม่
	var count int64
	err := tx.Model(&barberBookingModels.Appointment{}).
		Where("tenant_id = ? AND barber_id = ? AND status IN ? AND deleted_at IS NULL", tenantID, barberID,
			[]string{
				string(barberBookingModels.StatusPending),
				string(barberBookingModels.StatusConfirmed)}).
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count == 0, nil
}

func (s *appointmentService) CheckBarberAvailability(ctx context.Context, tenantID, barberID uint, start, end time.Time) (bool, error) {
	tx := s.DB.WithContext(ctx)
	return s.checkBarberAvailabilityTx(tx, tenantID, barberID, start, end)
}

func (s *appointmentService) CreateAppointment(ctx context.Context, input *barberBookingModels.Appointment) (*barberBookingModels.Appointment, error) {
	if input == nil {
		return nil, errors.New("input appointment data is required")
	}
	if input.TenantID == 0 || input.ServiceID == 0 || input.CustomerID == 0 || input.StartTime.IsZero() {
		return nil, errors.New("missing required fields")
	}

	var result *barberBookingModels.Appointment

	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. ดึง service เพื่อใช้ duration + validate tenant
		var service barberBookingModels.Service
		if err := tx.Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", input.ServiceID, input.TenantID).
			First(&service).Error; err != nil {
			return fmt.Errorf("service not found or access denied")
		}

		// หลังจากดึง service ได้สำเร็จ
		if service.Duration <= 0 {
			return fmt.Errorf("duration must be > 0")
		}

		// 2. คำนวณ EndTime จาก StartTime + Duration
		startTime := input.StartTime
		endTime := startTime.Add(time.Duration(service.Duration) * time.Minute)
		input.EndTime = endTime

		// 3. ถ้ามี barber → ตรวจสอบ availability
		if input.BarberID != nil {
			// ตรวจสอบ barber และ branch
			var barber barberBookingModels.Barber
			if err := tx.Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", *input.BarberID, input.TenantID).
				First(&barber).Error; err != nil {
				return fmt.Errorf("barber not found or mismatched branch")
			}
			if barber.BranchID != input.BranchID {
				return fmt.Errorf("barber not found or mismatched branch")
			}

			// ตรวจสอบ availability
			available, err := s.checkBarberAvailabilityTx(tx, input.TenantID, *input.BarberID, startTime, endTime)
			if err != nil {
				return fmt.Errorf("check barber availability failed: %w", err)
			}
			if !available {
				return fmt.Errorf("barber is not available during this time")
			}
		}

		// 4. ตั้งค่าข้อมูลและสร้าง appointment
		if input.Status == "" {
			input.Status = barberBookingModels.StatusPending
		}
		input.CreatedAt = time.Now()
		input.UpdatedAt = time.Now()

		if err := tx.Create(input).Error; err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		result = input
		return nil
	})

	return result, err
}

func (s *appointmentService) GetAvailableBarbers(ctx context.Context, tenantID, branchID uint, start, end time.Time) ([]barberBookingModels.Barber, error) {
	var barbers []barberBookingModels.Barber

	activeStatuses := []barberBookingModels.AppointmentStatus{
		barberBookingModels.StatusPending,
		barberBookingModels.StatusConfirmed,
		barberBookingModels.StatusRescheduled,
	}

	tx := s.DB.WithContext(ctx)

	subQuery := tx.Model(&barberBookingModels.Appointment{}).
		Select("barber_id").
		Where(`
			tenant_id = ?
			AND barber_id IS NOT NULL
			AND status IN ?
			AND NOT (end_time <= ? OR start_time >= ?)
		`, tenantID, activeStatuses, start, end)

	err := tx.Model(&barberBookingModels.Barber{}).
		Where("tenant_id = ? AND branch_id = ? AND id NOT IN (?)", tenantID, branchID, subQuery).
		Find(&barbers).Error

	return barbers, err
}

func (s *appointmentService) UpdateAppointment(ctx context.Context, id uint, tenantID uint, input *barberBookingModels.Appointment) (*barberBookingModels.Appointment, error) {
	if input == nil {
		return nil, errors.New("input appointment data is required")
	}

	var ap barberBookingModels.Appointment
	tx := s.DB.WithContext(ctx)

	//  ดึง appointment เดิมมา
	if err := tx.Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID).First(&ap).Error; err != nil {
		return nil, fmt.Errorf("appointment not found")
	}

	//  ตรวจสอบว่า service ใหม่ถูกต้อง (ถ้ามีการแก้)
	if input.ServiceID != 0 && input.ServiceID != ap.ServiceID {
		var svc barberBookingModels.Service
		if err := tx.Where("id = ? AND tenant_id = ?", input.ServiceID, tenantID).First(&svc).Error; err != nil {
			return nil, fmt.Errorf("service not found or access denied")
		}
		if svc.Duration <= 0 {
			return nil, fmt.Errorf("duration must be > 0")
		}
		ap.ServiceID = input.ServiceID
		ap.EndTime = input.StartTime.Add(time.Duration(svc.Duration) * time.Minute)
	}

	//  ตรวจสอบ barber ใหม่ (ถ้ามีการเปลี่ยน)
	if input.BarberID != nil {
		var barber barberBookingModels.Barber
		if err := tx.Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", *input.BarberID, tenantID).First(&barber).Error; err != nil {
			return nil, fmt.Errorf("barber not found or access denied")
		}
		if input.BranchID != 0 && barber.BranchID != input.BranchID {
			return nil, fmt.Errorf("barber mismatched branch")
		}

		available, err := s.checkBarberAvailabilityTx(tx, tenantID, *input.BarberID, input.StartTime, ap.EndTime)
		if err != nil {
			return nil, fmt.Errorf("check barber availability failed: %w", err)
		}
		if !available {
			return nil, fmt.Errorf("barber is not available during this time")
		}
		ap.BarberID = input.BarberID
	}

	//  อัปเดตฟิลด์ทั่วไป
	if !input.StartTime.IsZero() {
		ap.StartTime = input.StartTime
	}
	if input.Status != "" {
		ap.Status = input.Status
	}
	ap.UpdatedAt = time.Now()

	//  Save
	if err := tx.Save(&ap).Error; err != nil {
		return nil, fmt.Errorf("failed to update appointment: %w", err)
	}

	return &ap, nil
}

func (s *appointmentService) GetByID(ctx context.Context, id uint) (*barberBookingModels.Appointment, error) {
    var appt barberBookingModels.Appointment
    // ดึงเฉพาะเรคคอร์ดที่ยังไม่ถูกลบ (deleted_at IS NULL)
    err := s.DB.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", id).
        First(&appt).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("appointment with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to fetch appointment: %w", err)
    }
    return &appt, nil
}
