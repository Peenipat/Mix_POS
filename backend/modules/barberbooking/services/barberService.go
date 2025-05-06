package barberBookingService

import (
	"context"
	"errors"
	"time"
	"fmt"

	"gorm.io/gorm"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type BarberService struct {
	DB *gorm.DB
}

func NewBarberService(db *gorm.DB) *BarberService {
	return &BarberService{DB: db}
}

// CreateBarber creates a new barber
func (s *BarberService) CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error {
	// Validation
	if barber.BranchID == 0 {
		return fmt.Errorf("branch_id is required")
	}
	if barber.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	// ลบ record เดิมที่ถูก soft-delete ไปแล้ว (ถ้ามี user_id เดิม)
	var existing barberBookingModels.Barber
	err := s.DB.WithContext(ctx).
		Unscoped().
		Where("user_id = ?", barber.UserID).
		First(&existing).Error

	if err == nil && existing.DeletedAt.Valid {
		// ถ้ามีและถูก soft-delete → ลบทิ้งจริงก่อน
		if err := s.DB.WithContext(ctx).Unscoped().Delete(&existing).Error; err != nil {
			return fmt.Errorf("failed to purge existing deleted barber: %w", err)
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing barber: %w", err)
	}

	// สร้างใหม่
	barber.CreatedAt = time.Now()
	barber.UpdatedAt = time.Now()
	return s.DB.WithContext(ctx).Create(barber).Error
}



// GetBarberByID fetches a single barber by ID
func (s *BarberService) GetBarberByID(ctx context.Context, id uint) (*barberBookingModels.Barber, error) {
	var barber barberBookingModels.Barber
	if err := s.DB.WithContext(ctx).First(&barber, id).Error; err != nil {
		return nil, err
	}
	return &barber, nil
}

// ListBarbers optionally filters by branch_id
func (s *BarberService) ListBarbers(ctx context.Context, branchID *uint) ([]barberBookingModels.Barber, error) {
	var barbers []barberBookingModels.Barber
	query := s.DB.WithContext(ctx).Model(&barberBookingModels.Barber{})
	if branchID != nil {
		query = query.Where("branch_id = ?", *branchID)
	}
	if err := query.Find(&barbers).Error; err != nil {
		return nil, err
	}
	return barbers, nil
}

// UpdateBarber updates barber info
func (s *BarberService) UpdateBarber(ctx context.Context, id uint, updated *barberBookingModels.Barber) (*barberBookingModels.Barber, error) {
	var barber barberBookingModels.Barber
	if err := s.DB.WithContext(ctx).First(&barber, id).Error; err != nil {
		return nil, err
	}

	barber.BranchID = updated.BranchID
	barber.UserID = updated.UserID
	barber.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Save(&barber).Error; err != nil {
		return nil, err
	}

	return &barber, nil
}

// DeleteBarber performs soft delete
func (s *BarberService) DeleteBarber(ctx context.Context, id uint) error {
	result := s.DB.WithContext(ctx).Delete(&barberBookingModels.Barber{}, id)
	if result.RowsAffected == 0 {
		return errors.New("barber not found")
	}
	return result.Error
}

func (s *BarberService) GetBarberByUser(ctx context.Context, userID uint) (*barberBookingModels.Barber, error) {
	var barber barberBookingModels.Barber
	err := s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&barber).Error
	if err != nil {
		return nil, err
	}
	return &barber, nil
}

func (s *BarberService) ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) {
	var barbers []barberBookingModels.Barber

	err := s.DB.WithContext(ctx).
		Joins("JOIN branches ON branches.id = barbers.branch_id").
		Where("branches.tenant_id = ?", tenantID).
		Where("barbers.deleted_at IS NULL").
		Find(&barbers).Error

	return barbers, err
}



