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

func (s *appointmentService) checkBarberAvailabilityTx(tx *gorm.DB, tenantID, barberID uint, start, end time.Time) (bool, error) {
	var count int64
	err := tx.Model(&barberBookingModels.Appointment{}).
		Where("tenant_id = ? AND barber_id = ? AND status IN ? AND start_time < ? AND end_time > ?",
			tenantID, barberID,
			[]barberBookingModels.AppointmentStatus{
				barberBookingModels.StatusPending,
				barberBookingModels.StatusConfirmed,
			},
			end, start,
		).Count(&count).Error

	if err != nil {
		return false, fmt.Errorf("failed to check barber availability: %w", err)
	}
	return count == 0, nil
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
