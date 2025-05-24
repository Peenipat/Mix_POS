package barberBookingControllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	barberBookingController "myapp/modules/barberbooking/controllers"
	models "myapp/modules/barberbooking/models"
	// port "myapp/modules/barberbooking/port"
	// helperFunc "myapp/modules/barberbooking"
)

// MockReviewService implements port.IAppointmentReview
type MockReviewService struct {
	mock.Mock
}

// CreateReview implements barberBookingPort.IAppointmentReview.
func (m *MockReviewService) CreateReview(ctx context.Context, review *models.AppointmentReview) (*models.AppointmentReview, error) {
    args := m.Called(ctx, review)
    r := args.Get(0)
    if r != nil {
        return r.(*models.AppointmentReview), args.Error(1)
    }
    return nil, args.Error(1)
}

// DeleteReview implements barberBookingPort.IAppointmentReview.

func (m *MockReviewService) DeleteReview(
    ctx context.Context,
    reviewID uint,
    actorCustomerID uint,
    actorRole string,
) error {
    args := m.Called(ctx, reviewID, actorCustomerID, actorRole)
    return args.Error(0)
}

// GetAverageRatingByBarber implements barberBookingPort.IAppointmentReview.
func (m *MockReviewService) GetAverageRatingByBarber(ctx context.Context, barberID uint) (float64, error) {
    args := m.Called(ctx, barberID)
    return args.Get(0).(float64), args.Error(1)
}

// GetReviewByAppointment implements barberBookingPort.IAppointmentReview.
func (m *MockReviewService) GetReviewByAppointment(ctx context.Context, appointmentID uint) (*models.AppointmentReview, error) {
	panic("unimplemented")
}

// GetReviews implements barberBookingPort.IAppointmentReview.
func (m *MockReviewService) GetReviews(ctx context.Context, appointmentID uint) ([]models.AppointmentReview, error) {
	panic("unimplemented")
}

// UpdateReview implements barberBookingPort.IAppointmentReview.

func (m *MockReviewService) UpdateReview(ctx context.Context, reviewID uint, input *models.AppointmentReview) (*models.AppointmentReview, error) {
    args := m.Called(ctx, reviewID, input)
    r := args.Get(0)
    if r != nil {
        return r.(*models.AppointmentReview), args.Error(1)
    }
    return nil, args.Error(1)
}

func (m *MockReviewService) GetByID(ctx context.Context, id uint) (*models.AppointmentReview, error) {
	args := m.Called(ctx, id)
	r := args.Get(0)
	if r != nil {
		return r.(*models.AppointmentReview), args.Error(1)
	}
	return nil, args.Error(1)
}

// setup Fiber app with authorization middleware and route
func setupAppointmentReviewApp(mockSvc *MockReviewService) *fiber.App {
	app := fiber.New()
	ctrl := barberBookingController.NewAppointmentReviewController(mockSvc)
	app.Get("/tenants/:tenant_id/reviews/:review_id", ctrl.GetReviewByID)

	group := app.Group("/tenants/:tenant_id/appointments/:appointment_id")
	group.Post("/reviews", ctrl.CreateReview)
	group.Put("/reviews/:review_id", ctrl.UpdateReview)
	group.Delete("/reviews/:review_id", ctrl.DeleteReview)
    group.Get("/barbers/:barber_id/average-rating", ctrl.GetAverageRatingByBarber)

    tenants := app.Group("/tenants/:tenant_id")
    tenants.Get("/barbers/:barber_id/average-rating", ctrl.GetAverageRatingByBarber)
	return app
}

func TestGetReviewByID(t *testing.T) {
	validReview := &models.AppointmentReview{
		ID:            7,
		AppointmentID: 18,
		CustomerID:    uintPtr(3),
		Rating:        5,
		Comment:       "Excellent",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	t.Run("InvalidTenantID_Returns400", func(t *testing.T) {
		mockSvc := new(MockReviewService)
		app := setupAppointmentReviewApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/abc/reviews/7", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidReviewID_Returns400", func(t *testing.T) {
		mockSvc := new(MockReviewService)
		app := setupAppointmentReviewApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/reviews/xyz", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidIDZero_Returns400", func(t *testing.T) {
		// service stub not needed since id=0 is checked before service call
		mockSvc := new(MockReviewService)
		app := setupAppointmentReviewApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/reviews/0", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceNotFound_Returns404", func(t *testing.T) {
		mockSvc := new(MockReviewService)
		mockSvc.
			On("GetByID", mock.Anything, uint(9)).
			Return(nil, errors.New("review with ID 9 not found"))

		app := setupAppointmentReviewApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/tenants/1/reviews/9", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError_Returns500", func(t *testing.T) {
		mockSvc := new(MockReviewService)
		mockSvc.
			On("GetByID", mock.Anything, uint(5)).
			Return(nil, errors.New("db failure"))

		app := setupAppointmentReviewApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/tenants/1/reviews/5", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_Returns200", func(t *testing.T) {
		mockSvc := new(MockReviewService)
		mockSvc.
			On("GetByID", mock.Anything, uint(7)).
			Return(validReview, nil)
	
		app := setupAppointmentReviewApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/tenants/1/reviews/7", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	
		var body struct {
			Status string                     `json:"status"`
			Data   *models.AppointmentReview  `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
	
		// เปรียบเทียบเฉพาะฟิลด์สำคัญ
		got := body.Data
		assert.Equal(t, validReview.ID, got.ID)
		assert.Equal(t, validReview.AppointmentID, got.AppointmentID)
		assert.Equal(t, *validReview.CustomerID, *got.CustomerID)
		assert.Equal(t, validReview.Rating, got.Rating)
		assert.Equal(t, validReview.Comment, got.Comment)
		// เปรียบเทียบเวลาภายใน 1 วินาที
		assert.WithinDuration(t, validReview.CreatedAt, got.CreatedAt, time.Second)
		assert.WithinDuration(t, validReview.UpdatedAt, got.UpdatedAt, time.Second)
	
		mockSvc.AssertExpectations(t)
	})
	
}

func TestCreateReviewController(t *testing.T) {
    now := time.Now().Truncate(time.Second)

    t.Run("InvalidTenantID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodPost, "/tenants/abc/appointments/18/reviews", strings.NewReader(`{"rating":5}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidAppointmentID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/xyz/reviews", strings.NewReader(`{"rating":5}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidJSON", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/18/reviews", strings.NewReader(`notjson`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidInput_RatingOutOfRange", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        body := `{"rating":10}`
        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/18/reviews", strings.NewReader(body))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("AppointmentLookupFailed", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        // stub service to return wrapped lookup error
        mockSvc.
            On("CreateReview", mock.Anything, mock.MatchedBy(func(r *models.AppointmentReview) bool {
                return r.AppointmentID == 18
            })).
            Return(nil, errors.New("appointment lookup failed: sql err"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":5}`
        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/18/reviews", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("CannotReviewAppointment_StatusNotComplete", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("CreateReview", mock.Anything, mock.Anything).
            Return(nil, errors.New("cannot review appointment: status is PENDING"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":4}`
        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/18/reviews", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("ServiceError_Returns500", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("CreateReview", mock.Anything, mock.Anything).
            Return(nil, errors.New("db error"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":3}`
        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/18/reviews", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("Success_Returns201", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        created := &models.AppointmentReview{
            ID:            7,
            AppointmentID: 18,
            CustomerID:    uintPtr(3),
            Rating:        5,
            Comment:       "Great!",
            CreatedAt:     now,
            UpdatedAt:     now,
        }
        mockSvc.
            On("CreateReview", mock.Anything, mock.MatchedBy(func(r *models.AppointmentReview) bool {
                // ensure payload fields mapped
                return r.AppointmentID == 18 && r.Rating == 5
            })).
            Return(created, nil)

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"customer_id":3,"rating":5,"comment":"Great!"}`
        req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/18/reviews", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusCreated, resp.StatusCode)

        var respBody struct {
            Status string                  `json:"status"`
            Data   *models.AppointmentReview `json:"data"`
        }
        err := json.NewDecoder(resp.Body).Decode(&respBody)
        assert.NoError(t, err)
        assert.Equal(t, "success", respBody.Status)
        got := respBody.Data
        assert.Equal(t, created.ID, got.ID)
        assert.Equal(t, created.AppointmentID, got.AppointmentID)
        assert.Equal(t, *created.CustomerID, *got.CustomerID)
        assert.Equal(t, created.Rating, got.Rating)
        assert.Equal(t, created.Comment, got.Comment)
        assert.WithinDuration(t, created.CreatedAt, got.CreatedAt, time.Second)
        assert.WithinDuration(t, created.UpdatedAt, got.UpdatedAt, time.Second)

        mockSvc.AssertExpectations(t)
    })
}

func TestUpdateReviewController(t *testing.T) {
    now := time.Now().Truncate(time.Second)

    t.Run("InvalidTenantID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/tenants/abc/appointments/18/reviews/7", strings.NewReader(`{"rating":5}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidReviewID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/xyz/reviews/7", strings.NewReader(`{"rating":5}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidJSON", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/18/reviews/7", strings.NewReader(`not-json`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidRating_Returns400", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        // service returns rating error
        mockSvc.
            On("UpdateReview", mock.Anything, uint(7), mock.MatchedBy(func(r *models.AppointmentReview) bool {
                return r.Rating == 10
            })).
            Return(nil, errors.New("invalid rating: must be between 1 and 5"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":10,"comment":"oops"}`
        req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/18/reviews/7", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("NotFound_Returns404", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("UpdateReview", mock.Anything, uint(7), mock.Anything).
            Return(nil, errors.New("review with ID 7 not found"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":4,"comment":"late"}` 
        req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/18/reviews/7", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("ServiceError_Returns500", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("UpdateReview", mock.Anything, uint(7), mock.Anything).
            Return(nil, errors.New("db error"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":3}`
        req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/18/reviews/7", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("Success_Returns200", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        updated := &models.AppointmentReview{
            ID:            7,
            AppointmentID: 18,
            CustomerID:    uintPtr(3),
            Rating:        5,
            Comment:       "updated",
            CreatedAt:     now.Add(-time.Hour),
            UpdatedAt:     now,
        }
        // match input Rating and Comment
        mockSvc.
            On("UpdateReview", mock.Anything, uint(7), mock.MatchedBy(func(r *models.AppointmentReview) bool {
                return r.Rating == 5 && r.Comment == "updated"
            })).
            Return(updated, nil)

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"rating":5,"comment":"updated"}`
        req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/18/reviews/7", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var respBody struct {
            Status string                   `json:"status"`
            Data   *models.AppointmentReview `json:"data"`
        }
        err := json.NewDecoder(resp.Body).Decode(&respBody)
        assert.NoError(t, err)
        assert.Equal(t, "success", respBody.Status)

        got := respBody.Data
        assert.Equal(t, updated.ID, got.ID)
        assert.Equal(t, updated.Rating, got.Rating)
        assert.Equal(t, updated.Comment, got.Comment)
        assert.WithinDuration(t, updated.UpdatedAt, got.UpdatedAt, time.Second)

        mockSvc.AssertExpectations(t)
    })
}

func TestDeleteReviewController(t *testing.T) {
    t.Run("InvalidTenantID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/tenants/abc/appointments/10/reviews/5", strings.NewReader(`{"actor_customer_id":3}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidAppointmentID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/xyz/reviews/5", strings.NewReader(`{"actor_customer_id":3}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidReviewID", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/10/reviews/abc", strings.NewReader(`{"actor_customer_id":3}`))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidJSON", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/10/reviews/5", strings.NewReader(`not-json`))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("NotFound_Returns404", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("DeleteReview", mock.Anything, uint(5), uint(3),  mock.Anything).
            Return(errors.New("review with ID 5 not found"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"actor_customer_id":3}`
        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/10/reviews/5", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("Forbidden_Returns403", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("DeleteReview", mock.Anything, uint(5), uint(3),  mock.Anything).
            Return(errors.New("you are not authorized to delete this review"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"actor_customer_id":3}`
        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/10/reviews/5", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusForbidden, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("ServiceError_Returns500", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("DeleteReview", mock.Anything, uint(5), uint(3),  mock.Anything).
            Return(errors.New("db error"))

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"actor_customer_id":3}`
        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/10/reviews/5", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("Success_Returns200", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("DeleteReview", mock.Anything, uint(5), uint(3),  mock.Anything).
            Return(nil)

        app := setupAppointmentReviewApp(mockSvc)
        body := `{"actor_customer_id":3}`
        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/appointments/10/reviews/5", strings.NewReader(body))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var bodyResp map[string]string
        json.NewDecoder(resp.Body).Decode(&bodyResp)
        assert.Equal(t, "success", bodyResp["status"])
        assert.Equal(t, "Review deleted", bodyResp["message"])
        mockSvc.AssertExpectations(t)
    })
}

func TestGetAverageRatingByBarberController(t *testing.T) {
    t.Run("InvalidTenantID_Returns400", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodGet, "/tenants/abc/barbers/2/average-rating", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidBarberID_Returns400", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        app := setupAppointmentReviewApp(mockSvc)

        req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/xyz/average-rating", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ServiceError_Returns500", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("GetAverageRatingByBarber", mock.Anything, uint(5)).
            Return(0.0, errors.New("db failure"))

        app := setupAppointmentReviewApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/5/average-rating", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("Success_Returns200", func(t *testing.T) {
        mockSvc := new(MockReviewService)
        mockSvc.
            On("GetAverageRatingByBarber", mock.Anything, uint(7)).
            Return(4.25, nil)

        app := setupAppointmentReviewApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/7/average-rating", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string  `json:"status"`
            Data   float64 `json:"data"`
        }
        err := json.NewDecoder(resp.Body).Decode(&body)
        assert.NoError(t, err)
        assert.Equal(t, "success", body.Status)
        assert.Equal(t, 4.25, body.Data)
        mockSvc.AssertExpectations(t)
    })
}

