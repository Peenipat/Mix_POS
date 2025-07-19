package barberBookingService

import (
	"context"
	"errors"

	"fmt"
	"gorm.io/gorm"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"strings"
)

type WorkingHourService struct {
	DB *gorm.DB
}

func NewWorkingHourService(db *gorm.DB) barberBookingPort.IWorkingHourService {
	return &WorkingHourService{DB: db}
}

func (s *WorkingHourService) GetWorkingHours(ctx context.Context, branchID uint, tenantID uint) ([]barberBookingModels.WorkingHour, error) {
	var hours []barberBookingModels.WorkingHour
	err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND tenant_id = ? AND deleted_at IS NULL", branchID, tenantID).
		Order("weekday asc").
		Find(&hours).Error
	if err != nil {
		return nil, err
	}
	return hours, nil
}

func (s *WorkingHourService) UpdateWorkingHours(ctx context.Context, branchID uint, tenantID uint, input []barberBookingDto.WorkingHourInput) error {
	tx := s.DB.WithContext(ctx).Begin()

	for _, wh := range input {
		if wh.Weekday < 0 || wh.Weekday > 6 {
			tx.Rollback()
			return fmt.Errorf("invalid weekday: %d", wh.Weekday)
		}

		var existing barberBookingModels.WorkingHour
		err := tx.
			Where("branch_id = ? AND tenant_id = ? AND weekday = ?", branchID, tenantID, wh.Weekday).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newWH := barberBookingModels.WorkingHour{
				BranchID:  branchID,
				TenantID:  tenantID,
				Weekday:   wh.Weekday,
				StartTime: wh.StartTime,
				EndTime:   wh.EndTime,
				IsClosed:  wh.IsClosed,
			}
			if err := tx.Create(&newWH).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if err == nil {
			// อัปเดต
			existing.StartTime = wh.StartTime
			existing.EndTime = wh.EndTime
			existing.IsClosed = wh.IsClosed
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}


func (s *WorkingHourService) CreateWorkingHours(ctx context.Context, branchID uint, input barberBookingDto.WorkingHourInput) error {
	if input.Weekday < 0 || input.Weekday > 6 {
		return fmt.Errorf("invalid weekday: %d", input.Weekday)
	}
	if input.StartTime.After(input.EndTime) || input.StartTime.Equal(input.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}

	wh := barberBookingModels.WorkingHour{
		BranchID:  branchID,
		Weekday:   input.Weekday,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
	}

	if err := s.DB.WithContext(ctx).Create(&wh).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("working hour for weekday %d already exists", input.Weekday)
		}
		return err
	}

	return nil
}
