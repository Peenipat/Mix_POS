package barberBookingService


import (
	"context"
	"fmt"
	"time"
	"errors"
	"database/sql"

	"gorm.io/gorm"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
)

// AppointmentReviewService handles creation and retrieval of appointment reviews.
type appointmentReviewService struct {
	DB                 *gorm.DB
	appointmentService barberBookingPort.IAppointment
}

// NewAppointmentReviewService constructs a new review service.
// You can inject the existing AppointmentService to validate appointments.
func NewAppointmentReviewService(db *gorm.DB) *appointmentReviewService {
	return &appointmentReviewService{DB:db}
}
func (s *appointmentReviewService) GetByID(ctx context.Context, id uint) (*barberBookingModels.AppointmentReview, error) {
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

// CreateReview creates a new review for a completed appointment.
func (s *appointmentReviewService) CreateReview(ctx context.Context, review *barberBookingModels.AppointmentReview) (*barberBookingModels.AppointmentReview, error) {
	// 1. Ensure required fields are present
	if review.AppointmentID == 0 || review.Rating < 1 || review.Rating > 5 {
		return nil, fmt.Errorf("invalid review input: appointmentID and rating (1-5) are required")
	}

	// 2. Validate appointment exists and is completed
	appt, err := s.appointmentService.GetAppointmentByID(ctx, review.AppointmentID)
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
func (s *appointmentReviewService) GetReviews(ctx context.Context, appointmentID uint) ([]barberBookingModels.AppointmentReview, error) {
	var reviews []barberBookingModels.AppointmentReview
	err := s.DB.WithContext(ctx).
		Where("appointment_id = ? AND deleted_at IS NULL", appointmentID).
		Find(&reviews).Error
	return reviews, err
}

func (s *appointmentReviewService) UpdateReview(ctx context.Context, reviewID uint, input *barberBookingModels.AppointmentReview) (*barberBookingModels.AppointmentReview, error) {
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

// (Optional)
func (s *appointmentReviewService) GetReviewByAppointment(ctx context.Context, appointmentID uint) (*barberBookingModels.AppointmentReview, error) {
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


func (s *appointmentReviewService) DeleteReview(
    ctx context.Context,
    reviewID uint,
    actorCustomerID uint,
) error {
    // 1. โหลดรีวิวจาก DB
    var review barberBookingModels.AppointmentReview
    if err := s.DB.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", reviewID).
        First(&review).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return fmt.Errorf("review with ID %d not found", reviewID)
        }
        return fmt.Errorf("failed fetching review: %w", err)
    }

    // 2. เช็คว่า actorCustomerID เป็นเจ้าของรีวิว
    if review.CustomerID == nil || *review.CustomerID != actorCustomerID {
        return errors.New("you are not authorized to delete this review")
    }

    // 3. Soft delete
    if err := s.DB.WithContext(ctx).Delete(&review).Error; err != nil {
        return fmt.Errorf("failed deleting review: %w", err)
    }
    return nil
}


func (s *appointmentReviewService) GetAverageRatingByBarber(ctx context.Context, barberID uint) (float64, error) {
	var avg sql.NullFloat64

	err := s.DB.WithContext(ctx).
		Table("appointment_reviews").
		Select("AVG(appointment_reviews.rating)").
		Joins("JOIN appointments ON appointment_reviews.appointment_id = appointments.id").
		Where("appointments.barber_id = ? AND appointment_reviews.deleted_at IS NULL AND appointments.deleted_at IS NULL", barberID).
		Scan(&avg).Error

	if err != nil {
		return 0, err
	}
	if !avg.Valid {
		return 0, nil // คืน 0 ถ้า NULL (ไม่มีรีวิว)
	}
	return avg.Float64, nil
}




