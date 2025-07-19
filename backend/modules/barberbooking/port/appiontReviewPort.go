package barberBookingPort
import (
	"context"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type IAppointmentReview interface{
	GetByID(ctx context.Context, id uint) (*barberBookingModels.AppointmentReview, error) 
	CreateReview(ctx context.Context, review *barberBookingModels.AppointmentReview) (*barberBookingModels.AppointmentReview, error)
	GetReviews(ctx context.Context, appointmentID uint) ([]barberBookingModels.AppointmentReview, error)
	UpdateReview(ctx context.Context, reviewID uint, input *barberBookingModels.AppointmentReview) (*barberBookingModels.AppointmentReview, error)
	GetReviewByAppointment(ctx context.Context, appointmentID uint) (*barberBookingModels.AppointmentReview, error)
	DeleteReview(ctx context.Context, reviewID uint, actorUserID uint, actorRole string) error
	GetAverageRatingByBarber(ctx context.Context, barberID uint) (float64, error)
}

type CreateAppointmentReviewRequest struct {
    CustomerID *uint  `json:"customer_id,omitempty" example:"4"`
    Rating     int    `json:"rating" example:"5"`
    Comment    string `json:"comment,omitempty" example:"Great service!"`
}

type UpdateAppointmentReviewRequest struct {
    Rating  int    `json:"rating" example:"4"`
    Comment string `json:"comment,omitempty" example:"Updated comment"`
}