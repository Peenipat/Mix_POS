package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"time"

	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"

	"gorm.io/gorm"
)

type WorkingDayOverrideService struct {
	DB *gorm.DB
}


func (s *WorkingDayOverrideService) GetAll(ctx context.Context) ([]barberBookingModels.WorkingDayOverride, error) {
	panic("unimplemented")
}

func (s *WorkingDayOverrideService) GetByFilter(ctx context.Context, filter barberBookingPort.WorkingDayOverrideFilter) ([]barberBookingModels.WorkingDayOverride, error) {
	panic("unimplemented")
}



func NewWorkingDayOverrideService(db *gorm.DB) barberBookingPort.IWorkingDayOverrideService {
	return &WorkingDayOverrideService{DB: db}
}

func (s *WorkingDayOverrideService) Create(
	ctx context.Context,
	input barberBookingPort.WorkingDayOverrideInput,
) (*barberBookingModels.WorkingDayOverride, error) {

	workDate, err := time.Parse("2006-01-02", input.WorkDate)
	if err != nil {
		return nil, fmt.Errorf("invalid work_date format (expected YYYY-MM-DD): %w", err)
	}
	startTime, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return nil, fmt.Errorf("invalid start_time format (expected HH:mm): %w", err)
	}
	endTime, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return nil, fmt.Errorf("invalid end_time format (expected HH:mm): %w", err)
	}

	// 2. ตรวจสอบว่า override เดิมมีอยู่แล้วหรือไม่ (branch + work_date ต้องไม่ซ้ำ)
	var existing barberBookingModels.WorkingDayOverride
	err = s.DB.WithContext(ctx).
		Where("branch_id = ? AND work_date = ?", input.BranchID, workDate).
		First(&existing).Error

	if err == nil {
		return nil, fmt.Errorf("override for branch %d on %s already exists", input.BranchID, input.WorkDate)
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, fmt.Errorf("failed to check existing override: %w", err)
	}

	// 3. สร้าง record ใหม่
	newOverride := &barberBookingModels.WorkingDayOverride{
		BranchID:  input.BranchID,
		WorkDate:  workDate,
		StartTime: helperFunc.TimeOnly{Time: startTime},
		EndTime:   helperFunc.TimeOnly{Time: endTime},
		IsClosed: input.IsClosed,
	}

	if err := s.DB.WithContext(ctx).Create(newOverride).Error; err != nil {
		return nil, fmt.Errorf("failed to create working_day_override: %w", err)
	}

	return newOverride, nil
}

func (s *WorkingDayOverrideService) Update(
	ctx context.Context,
	id uint,
	input barberBookingPort.WorkingDayOverrideInput,
) error {

	workDate, err := time.Parse("2006-01-02", input.WorkDate)
	if err != nil {
		return fmt.Errorf("invalid work_date format (expected YYYY-MM-DD): %w", err)
	}
	startTime, err := time.Parse("15:04", input.StartTime)
	if err != nil {
		return fmt.Errorf("invalid start_time format (expected HH:mm): %w", err)
	}
	endTime, err := time.Parse("15:04", input.EndTime)
	if err != nil {
		return fmt.Errorf("invalid end_time format (expected HH:mm): %w", err)
	}

	var override barberBookingModels.WorkingDayOverride
	if err := s.DB.WithContext(ctx).First(&override, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("override with ID %d not found", id)
		}
		return fmt.Errorf("failed to fetch existing override: %w", err)
	}

	if override.WorkDate != workDate || override.BranchID != input.BranchID {
		var conflict barberBookingModels.WorkingDayOverride
		err = s.DB.WithContext(ctx).
			Where("branch_id = ? AND work_date = ? AND id != ?", input.BranchID, workDate, id).
			First(&conflict).Error

		if err == nil {
			return fmt.Errorf("another override already exists for branch %d on %s", input.BranchID, input.WorkDate)
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("failed to check for conflicting override: %w", err)
		}
	}

	override.BranchID = input.BranchID
	override.WorkDate = workDate
	override.StartTime =  helperFunc.TimeOnly{Time: startTime}
	override.EndTime =  helperFunc.TimeOnly{Time: endTime}
	override.IsClosed = input.IsClosed

	if err := s.DB.WithContext(ctx).Save(&override).Error; err != nil {
		return fmt.Errorf("failed to update working_day_override: %w", err)
	}

	return nil
}

func (s *WorkingDayOverrideService) GetByID(
	ctx context.Context,
	id uint,
) (*barberBookingModels.WorkingDayOverride, error) {

	var override barberBookingModels.WorkingDayOverride

	if err := s.DB.WithContext(ctx).First(&override, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("working day override with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch working day override: %w", err)
	}

	return &override, nil
}

func (s *WorkingDayOverrideService) Delete(ctx context.Context, id uint) error {
	var record barberBookingModels.WorkingDayOverride

	// ตรวจสอบว่า record นี้มีอยู่จริงไหม
	if err := s.DB.WithContext(ctx).First(&record, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return fmt.Errorf("working_day_override with id %d not found", id)
		}
		return fmt.Errorf("failed to fetch override: %w", err)
	}

	// ลบข้อมูล (soft delete หากใช้ gorm.DeletedAt)
	if err := s.DB.WithContext(ctx).Delete(&record).Error; err != nil {
		return fmt.Errorf("failed to delete working_day_override: %w", err)
	}

	return nil
}

func (s *WorkingDayOverrideService) GetOverridesByDateRange(ctx context.Context, branchID uint, startDate, endDate time.Time) ([]barberBookingModels.WorkingDayOverride, error) {
	var overrides []barberBookingModels.WorkingDayOverride
	err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND work_date BETWEEN ? AND ? AND deleted_at IS NULL", branchID, startDate, endDate).
		Order("work_date ASC").
		Find(&overrides).Error

	if err != nil {
		return nil, err
	}
	return overrides, nil
}


