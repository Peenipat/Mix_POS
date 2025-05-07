package barberBookingService

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	barberbookingmodels "myapp/modules/barberbooking/models"
)

type barberWorkloadService struct {
	DB *gorm.DB
}

func NewBarberWorkloadService(db *gorm.DB) *barberWorkloadService {
	return &barberWorkloadService{DB: db}
}

// GetWorkloadByBarber: ดึง workload รายวันของช่าง
func (s *barberWorkloadService) GetWorkloadByBarber(ctx context.Context, barberID uint, date time.Time) (*barberbookingmodels.BarberWorkload, error) {
	var workload barberbookingmodels.BarberWorkload
	err := s.DB.WithContext(ctx).
		Where("barber_id = ? AND strftime('%Y-%m-%d', date) = ?", barberID, date.Format("2006-01-02")).
		First(&workload).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &workload, nil
}

// GetWorkloadByDate: ดึง workload ช่างทั้งหมดในวันนั้น
func (s *barberWorkloadService) GetWorkloadByDate(ctx context.Context, date time.Time) ([]barberbookingmodels.BarberWorkload, error) {
	var workloads []barberbookingmodels.BarberWorkload
	start := date.Truncate(24 * time.Hour)
	end := start.Add(24 * time.Hour)

	err := s.DB.WithContext(ctx).
		Where("date >= ? AND date < ?", start, end).
		Find(&workloads).Error
	if err != nil {
		return nil, err
	}
	return workloads, nil
}

func (s *barberWorkloadService) UpsertBarberWorkload(ctx context.Context, barberID uint, date time.Time, appointments int, hours int) error {
	//  Truncate เวลาออกให้เหลือแค่วัน เพื่อความแม่นยำในการเปรียบเทียบ
	date = date.Truncate(24 * time.Hour)

	var workload barberbookingmodels.BarberWorkload
	err := s.DB.WithContext(ctx).
		Where("barber_id = ? AND date = ?", barberID, date).
		First(&workload).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Insert ใหม่
		workload = barberbookingmodels.BarberWorkload{
			BarberID:          barberID,
			Date:              date,
			TotalAppointments: appointments,
			TotalHours:        hours,
		}
		return s.DB.WithContext(ctx).Create(&workload).Error
	} else if err != nil {
		return err
	}

	// Update
	workload.TotalAppointments = appointments
	workload.TotalHours = hours
	return s.DB.WithContext(ctx).Save(&workload).Error
}
