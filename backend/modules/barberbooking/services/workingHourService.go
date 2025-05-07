package barberBookingService


import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"strings"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingDto "myapp/modules/barberbooking/dto"

)

type WorkingHourService struct {
	DB *gorm.DB
}



func NewWorkingHourService(db *gorm.DB) *WorkingHourService {
	return &WorkingHourService{DB: db}
}



func (s *WorkingHourService) GetWorkingHours(ctx context.Context, branchID uint) ([]barberBookingModels.WorkingHour, error) {
	var hours []barberBookingModels.WorkingHour
	err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND deleted_at IS NULL", branchID).
		Order("weekday asc").
		Find(&hours).Error
	if err != nil {
		return nil, err
	}
	return hours, nil
}

func (s *WorkingHourService) UpdateWorkingHours(ctx context.Context, branchID uint, input []barberBookingDto.WorkingHourInput) error {
	tx := s.DB.WithContext(ctx).Begin()

	for _, wh := range input {
		if wh.Weekday < 0 || wh.Weekday > 6 {
			tx.Rollback()
			return fmt.Errorf("invalid weekday: %d", wh.Weekday)
		}
		data := barberBookingModels.WorkingHour{
			BranchID:  branchID,
			Weekday:   wh.Weekday,
			StartTime: wh.StartTime,
			EndTime:   wh.EndTime,
		}
		err := tx.
			Where("branch_id = ? AND weekday = ?", branchID, wh.Weekday).
			Assign(data). // UPDATE IF EXISTS
			FirstOrCreate(&data).Error
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}


func (s *WorkingHourService) GetBranchOpenStatus(ctx context.Context, branchID uint, weekday int, now time.Time) (bool, error) {
	var wh barberBookingModels.WorkingHour
	err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND weekday = ?", branchID, weekday).
		First(&wh).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return false, nil // ไม่มีข้อมูลถือว่าปิด
		}
		return false, err
	}

	start := time.Date(now.Year(), now.Month(), now.Day(), wh.StartTime.Hour(), wh.StartTime.Minute(), 0, 0, now.Location())
	end := time.Date(now.Year(), now.Month(), now.Day(), wh.EndTime.Hour(), wh.EndTime.Minute(), 0, 0, now.Location())

	isOpen := now.After(start) && now.Before(end)
	return isOpen, nil
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

	// สร้างใหม่เท่านั้น ถ้ามีแล้วห้ามซ้ำ
	if err := s.DB.WithContext(ctx).Create(&wh).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("working hour for weekday %d already exists", input.Weekday)
		}
		return err
	}

	return nil
}

