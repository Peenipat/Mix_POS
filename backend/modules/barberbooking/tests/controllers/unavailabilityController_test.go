package barberBookingControllers

import (
	// "errors"
	"encoding/json"
	"time"
	"context"
	"gorm.io/gorm"
	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUnavailabilityService struct {
	mock.Mock
}

func (m *MockUnavailabilityService) CreateUnavailability(ctx context.Context, input *barberBookingModels.Unavailability) (*barberBookingModels.Unavailability, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Unavailability), args.Error(1)
}

func (m *MockUnavailabilityService) GetUnavailabilitiesByBranch(ctx context.Context, branchID uint, from, to time.Time) ([]barberBookingModels.Unavailability, error) {
	args := m.Called(ctx, branchID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]barberBookingModels.Unavailability), args.Error(1)
}

func (m *MockUnavailabilityService) GetUnavailabilitiesByBarber(ctx context.Context, barberID uint, from, to time.Time) ([]barberBookingModels.Unavailability, error) {
	args := m.Called(ctx, barberID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]barberBookingModels.Unavailability), args.Error(1)
}

func (m *MockUnavailabilityService) UpdateUnavailability(ctx context.Context, id uint, updates map[string]interface{}) error {
	args := m.Called(ctx, id, updates)
	return args.Error(0)
}

func (m *MockUnavailabilityService) DeleteUnavailability(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)
	return args.Error(0)
}


func setupUnavailabilityTestApp(mockSvc barberBookingPort.IUnavailabilitySerivce) *fiber.App {
	app := fiber.New()
	ctrl := barberBookingControllers.NewUnavailabilityController(mockSvc)

	app.Use(func(c *fiber.Ctx) error {
		role := c.Get("X-Mock-Role")
		if role != "" {
			c.Locals("role", role)
		}
		return c.Next()
	})

	app.Post("/unavailabilities", ctrl.CreateUnavailability)
	app.Get("/branches/:branch_id/unavailabilities", ctrl.GetUnavailabilitiesByBranch)
	app.Get("/barbers/:barber_id/unavailabilities", ctrl.GetUnavailabilitiesByBarber)
	app.Patch("/unavailabilities/:id", ctrl.UpdateUnavailability)
	app.Delete("/unavailabilities/:id", ctrl.DeleteUnavailability)


	return app
}

func TestCreateUnavailability(t *testing.T) {
	mockSvc := new(MockUnavailabilityService)
	app := setupUnavailabilityTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/unavailabilities", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/unavailabilities", strings.NewReader("invalid-json"))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn201", func(t *testing.T) {
		payload := `{"branch_id": 1, "date": "2025-06-01T00:00:00Z", "reason": "Public holiday"}`
		expected := &barberBookingModels.Unavailability{
			ID:       1,
			BranchID: ptrUint(1),
			Date:     time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC),
			Reason:   "Public holiday",
		}

		mockSvc.On("CreateUnavailability", mock.Anything, mock.AnythingOfType("*barberBookingModels.Unavailability")).Return(expected, nil)

		req := httptest.NewRequest(http.MethodPost, "/unavailabilities", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var body struct {
			Status  string                            `json:"status"`
			Message string                            `json:"message"`
			Data    *barberBookingModels.Unavailability `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Unavailability created", body.Message)
		assert.Equal(t, uint(1), body.Data.ID)
	})

	// t.Run("DuplicateDate_ShouldReturn409", func(t *testing.T) {
	// 	payload := `{"barber_id": 1, "branch_id": 1, "date": "2025-06-01", "reason": "duplicate test"}`
	// 	mockSvc.On("CreateUnavailability", mock.Anything, mock.AnythingOfType("*barberBookingModels.Unavailability")).
	// 		Return(nil, errors.New("unavailability already exists for this date")).Once()
	
	// 	req := httptest.NewRequest(http.MethodPost, "/unavailabilities", strings.NewReader(payload))
	// 	req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
	// 	req.Header.Set("Content-Type", "application/json")
	// 	resp, err := app.Test(req)
	
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, http.StatusConflict, resp.StatusCode)
	// })
	
	t.Run("MissingRequiredField_ShouldReturn400", func(t *testing.T) {
		payload := `{"barber_id": 1}` // ขาด branch_id และ date
		req := httptest.NewRequest(http.MethodPost, "/unavailabilities", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
	
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// t.Run("InternalServerError_ShouldReturn500", func(t *testing.T) {
	// 	payload := `{"barber_id": 1, "branch_id": 1, "date": "2025-06-01", "reason": "internal error"}`
	// 	mockSvc.On("CreateUnavailability", mock.Anything, mock.AnythingOfType("*barberBookingModels.Unavailability")).
	// 		Return(nil, errors.New("DB connection lost")).Once()
	
	// 	req := httptest.NewRequest(http.MethodPost, "/unavailabilities", strings.NewReader(payload))
	// 	req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
	// 	req.Header.Set("Content-Type", "application/json")
	// 	resp, err := app.Test(req)
	
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	// })
	
}

func TestGetUnavailabilitiesByBranch(t *testing.T) {
	mockSvc := new(MockUnavailabilityService)
	app := setupUnavailabilityTestApp(mockSvc)

	route := "/branches/1/unavailabilities?from=2025-05-01&to=2025-05-31"

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/branches/abc/unavailabilities?from=2025-05-01&to=2025-05-31", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("GetUnavailabilitiesByBranch", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodGet, route, nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturnList", func(t *testing.T) {
		unavails := []barberBookingModels.Unavailability{
			{ID: 1, BranchID: uintPtr(1), Date: time.Date(2025, 5, 10, 0, 0, 0, 0, time.UTC)},
			{ID: 2, BranchID: uintPtr(1), Date: time.Date(2025, 5, 12, 0, 0, 0, 0, time.UTC)},
		}
		mockSvc.On("GetUnavailabilitiesByBranch", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(unavails, nil).Once()

		req := httptest.NewRequest(http.MethodGet, route, nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                                `json:"status"`
			Message string                                `json:"message"`
			Data    []barberBookingModels.Unavailability `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Unavailabilities retrieved", body.Message)
		assert.Len(t, body.Data, 2)
	})
}

func TestGetUnavailabilitiesByBarber(t *testing.T) {
	mockSvc := new(MockUnavailabilityService)
	app := setupUnavailabilityTestApp(mockSvc)

	route := "/barbers/1/unavailabilities?from=2025-05-01&to=2025-05-31"

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/barbers/abc/unavailabilities?from=2025-05-01&to=2025-05-31", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingQueryParams_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/barbers/1/unavailabilities", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidDateFormat_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/barbers/1/unavailabilities?from=bad&to=2025-05-31", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("GetUnavailabilitiesByBarber", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodGet, route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturnList", func(t *testing.T) {
		unavails := []barberBookingModels.Unavailability{
			{ID: 1, BarberID: uintPtr(1), Date: time.Date(2025, 5, 5, 0, 0, 0, 0, time.UTC)},
			{ID: 2, BarberID: uintPtr(1), Date: time.Date(2025, 5, 20, 0, 0, 0, 0, time.UTC)},
		}
		mockSvc.On("GetUnavailabilitiesByBarber", mock.Anything, uint(1), mock.Anything, mock.Anything).Return(unavails, nil).Once()

		req := httptest.NewRequest(http.MethodGet, route, nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                                `json:"status"`
			Message string                                `json:"message"`
			Data    []barberBookingModels.Unavailability `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Unavailabilities retrieved", body.Message)
		assert.Len(t, body.Data, 2)
	})
}

func TestUpdateUnavailability(t *testing.T) {
	mockSvc := new(MockUnavailabilityService)
	app := setupUnavailabilityTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/unavailabilities/1", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/unavailabilities/abc", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/unavailabilities/1", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("EmptyUpdateFields_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPatch, "/unavailabilities/1", strings.NewReader(`{}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("UpdateUnavailability", mock.Anything, uint(1), mock.Anything).
			Return(gorm.ErrRecordNotFound).Once()

		req := httptest.NewRequest(http.MethodPatch, "/unavailabilities/1", strings.NewReader(`{"reason":"updated"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc.On("UpdateUnavailability", mock.Anything, uint(1), mock.Anything).
			Return(nil).Once()

		req := httptest.NewRequest(http.MethodPatch, "/unavailabilities/1", strings.NewReader(`{"reason":"Rescheduled"}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Unavailability updated", body.Message)
	})
}

func TestDeleteUnavailability(t *testing.T) {
	mockSvc := new(MockUnavailabilityService)
	app := setupUnavailabilityTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/unavailabilities/1", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/unavailabilities/abc", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("DeleteUnavailability", mock.Anything, uint(1)).
			Return(gorm.ErrRecordNotFound).Once()

		req := httptest.NewRequest(http.MethodDelete, "/unavailabilities/1", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc.On("DeleteUnavailability", mock.Anything, uint(1)).
			Return(nil).Once()

		req := httptest.NewRequest(http.MethodDelete, "/unavailabilities/1", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Unavailability deleted", body.Message)
	})
}



func uintPtr(v uint) *uint {
	return &v
}


func ptrUint(v uint) *uint {
	return &v
}
