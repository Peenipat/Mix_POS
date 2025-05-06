package barberBookingService
import (
	"context"
	"errors"
	"time"
	"strings"

	"gorm.io/gorm"

	barberModels "myapp/modules/barberbooking/models"
)

type UnavailabilityService struct {
	DB *gorm.DB
}

func NewUnavailabilityService(db *gorm.DB) *UnavailabilityService {
	return &UnavailabilityService{DB: db}
}

// CreateUnavailability creates a new unavailability entry
func (s *UnavailabilityService) CreateUnavailability(ctx context.Context, input *barberModels.Unavailability) (*barberModels.Unavailability, error) {
	var existing barberModels.Unavailability
	err := s.DB.WithContext(ctx).Where("date = ? AND barber_id = ? AND branch_id = ?", input.Date, input.BarberID, input.BranchID).
		First(&existing).Error
	if err == nil {
		return nil, errors.New("unavailability already exists for this date")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	if err := s.DB.WithContext(ctx).Create(&input).Error; err != nil {
        if strings.Contains(err.Error(), "duplicate key") {
            return nil, errors.New("unavailability already exists for this date")
        }
        return nil, err
    }
	return input, nil
}

// GetUnavailabilitiesByBranch returns all unavailabilities for a specific branch within a date range
func (s *UnavailabilityService) GetUnavailabilitiesByBranch(ctx context.Context, branchID uint, from, to time.Time) ([]barberModels.Unavailability, error) {
	var results []barberModels.Unavailability
	err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND date BETWEEN ? AND ?", branchID, from, to).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}


// GetUnavailabilitiesByBarber returns all unavailabilities for a specific barber within a date range
func (s *UnavailabilityService) GetUnavailabilitiesByBarber(ctx context.Context, barberID uint, from, to time.Time) ([]barberModels.Unavailability, error) {
	var results []barberModels.Unavailability
	err := s.DB.WithContext(ctx).
		Where("barber_id = ? AND date BETWEEN ? AND ?", barberID, from, to).
		Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

// UpdateUnavailability updates an existing unavailability entry
func (s *UnavailabilityService) UpdateUnavailability(ctx context.Context, id uint, updates map[string]interface{}) error {
	result := s.DB.WithContext(ctx).
		Model(&barberModels.Unavailability{}).
		Where("id = ?", id).
		Updates(updates)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// DeleteUnavailability soft deletes an unavailability
func (s *UnavailabilityService) DeleteUnavailability(ctx context.Context, id uint) error {
	result := s.DB.WithContext(ctx).
		Where("id = ?", id).
		Delete(&barberModels.Unavailability{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
