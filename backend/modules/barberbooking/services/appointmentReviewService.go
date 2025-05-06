package barberBookingService

import (
	"context"
	"fmt"
	"time"
	"errors"
	"strings"

	"gorm.io/gorm"

	barberBookingModels "myapp/modules/barberbooking/models"
	coreModels "myapp/modules/core/models"
)

// AppointmentReviewService handles creation and retrieval of appointment reviews.
type AppointmentReviewService struct {
	DB                 *gorm.DB
	appointmentService *appointmentService
}

func (s *AppointmentReviewService) GetByID(ctx context.Context, id uint) (*barberBookingModels.AppointmentReview, error) {
    var rev barberBookingModels.AppointmentReview
    err := s.DB.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", id).
        First(&rev).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, fmt.Errorf("review with ID %d not found", id)
        }
        return nil, fmt.Errorf("failed to fetch review: %w", err)
    }
    return &rev, nil
}

// NewAppointmentReviewService constructs a new review service.
// You can inject the existing AppointmentService to validate appointments.
func NewAppointmentReviewService(db *gorm.DB, apptSvc *appointmentService) *AppointmentReviewService {
	return &AppointmentReviewService{
		DB:                 db,
		appointmentService: apptSvc,
	}
}

// CreateReview creates a new review for a completed appointment.
func (s *AppointmentReviewService) CreateReview(ctx context.Context, review *barberBookingModels.AppointmentReview) (*barberBookingModels.AppointmentReview, error) {
	// 1. Ensure required fields are present
	if review.AppointmentID == 0 || review.Rating < 1 || review.Rating > 5 {
		return nil, fmt.Errorf("invalid review input: appointmentID and rating (1-5) are required")
	}

	// 2. Validate appointment exists and is completed
	appt, err := s.appointmentService.GetByID(ctx, review.AppointmentID)
	if err != nil {
		return nil, fmt.Errorf("appointment lookup failed: %w", err)
	}
	if appt.Status != barberBookingModels.StatusComplete {
		return nil, fmt.Errorf("cannot review appointment: status is %s", appt.Status)
	}

	// 3. Populate timestamps
	now := time.Now()
	review.CreatedAt = now
	review.UpdatedAt = now

	// 4. Persist the review
	if err := s.DB.WithContext(ctx).Create(review).Error; err != nil {
		return nil, fmt.Errorf("failed to create review: %w", err)
	}

	return review, nil
}

// (Optional) GetReviews fetches all reviews for a given appointment.
func (s *AppointmentReviewService) GetReviews(ctx context.Context, appointmentID uint) ([]barberBookingModels.AppointmentReview, error) {
	var reviews []barberBookingModels.AppointmentReview
	err := s.DB.WithContext(ctx).
		Where("appointment_id = ? AND deleted_at IS NULL", appointmentID).
		Find(&reviews).Error
	return reviews, err
}

func (s *AppointmentReviewService) UpdateReview(ctx context.Context, reviewID uint, input *barberBookingModels.AppointmentReview) (*barberBookingModels.AppointmentReview, error) {
	if input.Rating < 1 || input.Rating > 5 {
		return nil, errors.New("invalid rating: must be between 1 and 5")
	}

	var existing barberBookingModels.AppointmentReview
	if err := s.DB.WithContext(ctx).
		Where("id = ?", reviewID).
		First(&existing).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("review with ID %d not found", reviewID)
		}
		return nil, err
	}

	existing.Rating = input.Rating
	existing.Comment = input.Comment
	existing.UpdatedAt = time.Now()

	if err := s.DB.WithContext(ctx).Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
}

func (s *AppointmentReviewService) GetReviewByAppointment(ctx context.Context, appointmentID uint) (*barberBookingModels.AppointmentReview, error) {
	var review barberBookingModels.AppointmentReview

	err := s.DB.WithContext(ctx).
		Where("appointment_id = ?", appointmentID).
		First(&review).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("review for appointment ID %d not found", appointmentID)
		}
		return nil, err
	}

	return &review, nil
}

func (s *AppointmentReviewService) DeleteReview(ctx context.Context, reviewID uint, actorUserID uint, actorRole string) error {
	var review barberBookingModels.AppointmentReview

	err := s.DB.WithContext(ctx).
		Where("id = ?", reviewID).
		First(&review).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return fmt.Errorf("review with ID %d not found", reviewID)
		}
		return err
	}

	// 🛡️ Check Permission
	switch strings.ToUpper(actorRole) {
	case string(coreModels.RoleNameUser):
		// USER ต้องเป็นเจ้าของ review (customer_id)
		if review.CustomerID == nil || *review.CustomerID != actorUserID {
			return errors.New("you are not authorized to delete this review")
		}

	case string(coreModels.RoleNameBranchAdmin),string(coreModels.RoleNameSaaSSuperAdmin):
		// ADMIN สามารถลบได้ทุกรีวิว
		// เพิ่มกรณี role เพิ่มเติมได้ที่นี่

	default:
		return fmt.Errorf("role %s is not authorized to delete reviews", actorRole)
	}

	// 🧼 Soft delete
	return s.DB.WithContext(ctx).Delete(&review).Error
}

