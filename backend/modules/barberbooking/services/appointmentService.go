package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingDto "myapp/modules/barberbooking/dto"
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

func (s *appointmentService) GetAppointmentByID(ctx context.Context, id uint) (*barberBookingModels.Appointment, error) {
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

func (s *appointmentService) ListAppointments(ctx context.Context, filter barberBookingDto.AppointmentFilter) ([]barberBookingModels.Appointment, error) {
	var appointments []barberBookingModels.Appointment

	tx := s.DB.WithContext(ctx).Model(&barberBookingModels.Appointment{})

	tx = tx.Where("tenant_id = ?", filter.TenantID)

	if filter.BranchID != nil {
		tx = tx.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.BarberID != nil {
		tx = tx.Where("barber_id = ?", *filter.BarberID)
	}
	if filter.CustomerID != nil {
		tx = tx.Where("customer_id = ?", *filter.CustomerID)
	}
	if filter.Status != nil {
		tx = tx.Where("status = ?", *filter.Status)
	}
	if filter.StartDate != nil {
		tx = tx.Where("start_time >= ?", *filter.StartDate)
	}
	if filter.EndDate != nil {
		tx = tx.Where("end_time <= ?", *filter.EndDate)
	}

	if err := tx.Order("start_time asc").Find(&appointments).Error; err != nil {
		return nil, err
	}

	if filter.Limit != nil {
		tx = tx.Limit(*filter.Limit)
	}
	if filter.Offset != nil {
		tx = tx.Offset(*filter.Offset)
	}

	// Sorting
	if filter.SortBy != nil && *filter.SortBy != "" {
		tx = tx.Order(*filter.SortBy)
	} else {
		tx = tx.Order("start_time asc") // default sort
	}

	if err := tx.Find(&appointments).Error; err != nil {
		return nil, err
	}
	return appointments, nil
}

func (s *appointmentService) CancelAppointment(ctx context.Context, appointmentID uint, actorUserID uint) error {
	var ap barberBookingModels.Appointment

	err := s.DB.WithContext(ctx).
		Where("id = ?", appointmentID).
		First(&ap).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("appointment with ID %d not found", appointmentID)
		}
		return err
	}

	//  ไม่อนุญาตให้ยกเลิกถ้าเป็น COMPLETED หรือ CANCELLED แล้ว
	if ap.Status == barberBookingModels.StatusComplete || ap.Status == barberBookingModels.StatusCancelled {
		return errors.New("appointment cannot be cancelled in its current status")
	}

	//  เปลี่ยนสถานะเป็น CANCELLED
	ap.Status = barberBookingModels.StatusCancelled
	ap.UpdatedAt = time.Now()
	ap.UserID = &actorUserID // ใครบันทึกการยกเลิก

	if err := s.DB.WithContext(ctx).Save(&ap).Error; err != nil {
		return err
	}

	//  (Optional) log ไปยัง appointment_status_logs ในอนาคต

	return nil
}

func (s *appointmentService) RescheduleAppointment( ctx context.Context,appointmentID uint,newStartTime time.Time,actorUserID uint,) error {
	// 1. ดึง appointment
	var ap barberBookingModels.Appointment
	if err := s.DB.WithContext(ctx).
		Preload("Service").
		First(&ap, "id = ?", appointmentID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("appointment with ID %d not found", appointmentID)
		}
		return err
	}

	// 2. ตรวจสอบสถานะ
	if ap.Status == barberBookingModels.StatusComplete || ap.Status == barberBookingModels.StatusCancelled {
		return fmt.Errorf("cannot reschedule a completed or cancelled appointment")
	}

	// 3. ตรวจสอบการชน (simplified logic)
	newEndTime := newStartTime.Add(time.Duration(ap.Service.Duration) * time.Minute)
	var conflict int64
	err := s.DB.WithContext(ctx).
		Model(&barberBookingModels.Appointment{}).
		Where("barber_id = ? AND branch_id = ? AND id != ? AND status IN ? AND start_time < ? AND end_time > ?",
			ap.BarberID, ap.BranchID, ap.ID,
			[]barberBookingModels.AppointmentStatus{
				barberBookingModels.StatusPending,
				barberBookingModels.StatusConfirmed,
			},
			newEndTime, newStartTime,
		).Count(&conflict).Error
	if err != nil {
		return err
	}
	if conflict > 0 {
		return fmt.Errorf("cannot reschedule: time slot conflicts with another appointment")
	}

	// 4. อัปเดตเวลาและผู้แก้ไข
	ap.StartTime = newStartTime
	ap.EndTime = newEndTime
	ap.UserID = &actorUserID
	ap.UpdatedAt = time.Now()
	ap.Status = barberBookingModels.StatusConfirmed // หรือเก็บสถานะเดิมไว้ก็ได้

	return s.DB.WithContext(ctx).Save(&ap).Error
}



