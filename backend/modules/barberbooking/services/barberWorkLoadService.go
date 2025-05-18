package barberBookingService

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	barberBookingDto "myapp/modules/barberbooking/dto"
)

type barberWorkloadService struct {
	DB *gorm.DB
}

func NewBarberWorkloadService(db *gorm.DB) barberBookingPort.IbarberWorkload {
	return &barberWorkloadService{DB: db}
}

// GetWorkloadByBarber: ดึง workload รายวันของช่าง
func (s *barberWorkloadService) GetWorkloadByBarber(ctx context.Context, barberID uint, date time.Time) (*barberBookingModels.BarberWorkload, error) {
	var workload barberBookingModels.BarberWorkload
	err := s.DB.
    WithContext(ctx).
    Where("barber_id = ? AND DATE(date) = ?", barberID, date.Format("2006-01-02")).
    First(&workload).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // not found is not an error
		}
		return nil, err
	}
	return &workload, nil
}

func (s *barberWorkloadService) GetWorkloadSummaryByBranch(
    ctx context.Context,
    date time.Time,
    tenantID uint,
    branchID uint,
) ([]barberBookingDto.BranchWorkloadSummary, error) {
    start := date.Truncate(24 * time.Hour)
    end := start.Add(24 * time.Hour)

    var sums []barberBookingDto.BranchWorkloadSummary
    query := s.DB.WithContext(ctx).
        Table("branches br").
        Select(
            `br.tenant_id,
             br.id AS branch_id,
             COUNT(DISTINCT bw.barber_id) AS num_worked,
             COUNT(b.id) AS total_barbers`,
        ).
        Joins("JOIN barbers b ON b.branch_id = br.id").
        Joins(`LEFT JOIN barber_workloads bw
               ON bw.barber_id = b.id
              AND bw.date >= ?
              AND bw.date < ?`, start, end)

    // Apply tenant filter if provided
    if tenantID != 0 {
        query = query.Where("br.tenant_id = ?", tenantID)
    }
    // Apply branch filter if provided
    if branchID != 0 {
        query = query.Where("br.id = ?", branchID)
    }

    query = query.Group("br.tenant_id, br.id").Scan(&sums)
    if query.Error != nil {
        return nil, query.Error
    }
    return sums, nil
}


func (s *barberWorkloadService) UpsertBarberWorkload(ctx context.Context, barberID uint, date time.Time, appointments int, hours int) error {
	//  Truncate เวลาออกให้เหลือแค่วัน เพื่อความแม่นยำในการเปรียบเทียบ
	date = date.Truncate(24 * time.Hour)

	var workload barberBookingModels.BarberWorkload
	err := s.DB.WithContext(ctx).
		Where("barber_id = ? AND date = ?", barberID, date).
		First(&workload).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		// Insert ใหม่
		workload = barberBookingModels.BarberWorkload{
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
	workload.TotalAppointments += appointments
	workload.TotalHours += hours
	return s.DB.WithContext(ctx).Save(&workload).Error
}
