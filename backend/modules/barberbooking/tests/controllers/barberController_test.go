package barberBookingControllers_test

import (
	// "bytes"
	"encoding/json"
	"mime/multipart"
	// "errors"
	"context"
	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBarberService struct {
	mock.Mock
}

// UpdateBarber implements barberBookingPort.IBarber.
func (m *MockBarberService) UpdateBarber(ctx context.Context, barberID uint, payload *barberBookingPort.UpdateBarberRequest, file *multipart.FileHeader) (*barberBookingModels.Barber, error) {
	panic("unimplemented")
}

// UpdateBarber implements barberBookingPort.IBarber.

// GetBarberByID implements barberBookingPort.IBarber.
func (m *MockBarberService) GetBarberByID(ctx context.Context, id uint) (*barberBookingPort.BarberDetailResponse, error) {
	panic("unimplemented")
}

// ListUserNotBarber implements barberBookingPort.IBarber.
func (m *MockBarberService) ListUserNotBarber(ctx context.Context, branchID *uint) ([]barberBookingPort.UserNotBarber, error) {
	panic("unimplemented")
}

// FindAvailableUsers implements barberBookingPort.IBarber.
func (m *MockBarberService) FindAvailableUsers(ctx context.Context, tenantID uint, branchID uint) ([]coreModels.User, error) {
	panic("unimplemented")
}

func (m *MockBarberService) CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error {
	args := m.Called(ctx, barber)
	return args.Error(0)
}

func (m *MockBarberService) ListBarbersByBranch(ctx context.Context, branchID *uint) ([]barberBookingPort.BarberWithUser, error) {
	args := m.Called(ctx, branchID)
	return args.Get(0).([]barberBookingPort.BarberWithUser), args.Error(1)
}

func (m *MockBarberService) DeleteBarber(ctx context.Context, barberID uint) error {
	args := m.Called(ctx, barberID)
	return args.Error(0)
}

func (m *MockBarberService) GetBarberByUser(ctx context.Context, userID uint) (*barberBookingModels.Barber, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Barber), args.Error(1)
}

func (m *MockBarberService) ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) {
	args := m.Called(ctx, tenantID)
	return args.Get(0).([]barberBookingModels.Barber), args.Error(1)
}

func setupBarberTestApp(mockSvc barberBookingPort.IBarber) *fiber.App {
	app := fiber.New()
	controller := barberBookingControllers.NewBarberController(mockSvc)

	app.Use(func(c *fiber.Ctx) error {
		role := c.Get("X-Mock-Role")
		if role != "" {
			c.Locals("role", role)
		}
		return c.Next()
	})

	app.Post("/barbers", controller.CreateBarber)
	app.Get("/barbers/:barber_id", controller.GetBarberByID)
	app.Get("/branches/:branch_id/barbers", controller.ListBarbersByBranch)
	// app.Put("/barbers/:barber_id", controller.UpdateBarber)
	app.Delete("/barbers/:barber_id", controller.DeleteBarber)
	app.Get("/users/:user_id/barber", controller.GetBarberByUser)
	app.Get("/tenants/:tenant_id/barbers", controller.ListBarbersByTenant)

	return app
}

func TestCreateBarber(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/barbers", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/barbers", strings.NewReader("invalid-json"))
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn201", func(t *testing.T) {
		payload := `{"branch_id": 1, "user_id": 10}`
		mockSvc.ExpectedCalls = nil
		mockSvc.On("CreateBarber", mock.Anything, mock.AnythingOfType("*barberBookingModels.Barber")).Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/barbers", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var body struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Barber created", body.Message)
	})
}

func TestGetBarberByID(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("BarberNotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("GetBarberByID", mock.Anything, uint(1)).Return(nil, nil)
		req := httptest.NewRequest(http.MethodGet, "/barbers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturnBarber", func(t *testing.T) {
		mockBarber := &barberBookingModels.Barber{ID: 1, BranchID: 2, UserID: 3}
		mockSvc.ExpectedCalls = nil
		mockSvc.On("GetBarberByID", mock.Anything, uint(1)).Return(mockBarber, nil)

		req := httptest.NewRequest(http.MethodGet, "/barbers/1", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                      `json:"status"`
			Message string                      `json:"message"`
			Data    *barberBookingModels.Barber `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Barber retrieved", body.Message)
		assert.Equal(t, mockBarber.ID, body.Data.ID)
	})
}

func TestListBarbersByBranch(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("Success_ShouldReturnBarberList", func(t *testing.T) {
		branchID := uint(1)
		mockList := []barberBookingModels.Barber{
			{ID: 1, BranchID: branchID, UserID: 10},
			{ID: 2, BranchID: branchID, UserID: 11},
		}
		mockSvc.ExpectedCalls = nil
		mockSvc.On("ListBarbersByBranch", mock.Anything, &branchID).Return(mockList, nil)

		req := httptest.NewRequest(http.MethodGet, "/branches/1/barbers", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                       `json:"status"`
			Message string                       `json:"message"`
			Data    []barberBookingModels.Barber `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "List Barber retrieved", body.Message)
		assert.Len(t, body.Data, 2)
		assert.Equal(t, uint(1), body.Data[0].ID)
	})
}

func TestUpdateBarber(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/barbers/1", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/barbers/1", strings.NewReader("invalid-json"))
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("BarberNotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("GetBarberByID", mock.Anything, uint(1)).Return(nil, nil)

		req := httptest.NewRequest(http.MethodPut, "/barbers/1", strings.NewReader(`{"branch_id": 2}`))
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		existing := &barberBookingModels.Barber{ID: 1, BranchID: 1, UserID: 2}
		updated := &barberBookingModels.Barber{ID: 1, BranchID: 2, UserID: 2}

		mockSvc.ExpectedCalls = nil
		mockSvc.On("GetBarberByID", mock.Anything, uint(1)).Return(existing, nil)
		mockSvc.On("UpdateBarber", mock.Anything, uint(1), mock.Anything).Return(updated, nil)

		req := httptest.NewRequest(http.MethodPut, "/barbers/1", strings.NewReader(`{"branch_id": 2}`))
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                      `json:"status"`
			Message string                      `json:"message"`
			Data    *barberBookingModels.Barber `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Barber Updated", body.Message)
		assert.Equal(t, updated.BranchID, body.Data.BranchID)
	})
}

func TestGetBarberByUser(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("Success_ShouldReturnBarber", func(t *testing.T) {
		barber := &barberBookingModels.Barber{
			ID:       1,
			BranchID: 2,
			UserID:   3,
		}
		mockSvc.On("GetBarberByUser", mock.Anything, uint(1)).Return(barber, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1/barber", nil)
		req.RequestURI = "/users/1/barber"
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                      `json:"status"`
			Message string                      `json:"message"`
			Data    *barberBookingModels.Barber `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Barber retrieved", body.Message)
		assert.Equal(t, uint(1), body.Data.ID)
	})

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/1/barber", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/users/abc/barber", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("BarberNotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("GetBarberByUser", mock.Anything, uint(1)).Return(nil, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/users/1/barber", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

}

func TestListBarbersByTenant(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tenants/abc/barbers", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NoBarbers_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("ListBarbersByTenant", mock.Anything, uint(1)).Return([]barberBookingModels.Barber{}, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturnList", func(t *testing.T) {
		barbers := []barberBookingModels.Barber{
			{ID: 1, BranchID: 1, UserID: 1},
			{ID: 2, BranchID: 2, UserID: 2},
		}
		mockSvc.On("ListBarbersByTenant", mock.Anything, uint(1)).Return(barbers, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                       `json:"status"`
			Message string                       `json:"message"`
			Data    []barberBookingModels.Barber `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Barbers retrieved", body.Message)
		assert.Len(t, body.Data, 2)
	})
}

func TestDeleteBarber(t *testing.T) {
	mockSvc := new(MockBarberService)
	app := setupBarberTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/barbers/1", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodDelete, "/barbers/abc", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("DeleteFailed_ShouldReturn500", func(t *testing.T) {
		mockSvc.On("DeleteBarber", mock.Anything, uint(1)).Return(assert.AnError)

		req := httptest.NewRequest(http.MethodDelete, "/barbers/1", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc.ExpectedCalls = nil
		mockSvc.On("DeleteBarber", mock.Anything, uint(1)).Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/barbers/1", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageBarber[0]))
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
		assert.Equal(t, "Barber delete successfully", body.Message)
	})
}
