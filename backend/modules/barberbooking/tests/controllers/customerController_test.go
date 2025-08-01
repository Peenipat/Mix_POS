package barberBookingControllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/gofiber/fiber/v2"

	"github.com/stretchr/testify/mock"
)

type MockCustomerService struct {
	mock.Mock
}

// GetCustomers implements barberBookingPort.ICustomer.
func (m *MockCustomerService) GetCustomers(ctx context.Context, filter barberBookingPort.GetCustomersFilter) ([]barberBookingModels.Customer, int64, error) {
	panic("unimplemented")
}

// GetPendingAndCancelledCount implements barberBookingPort.ICustomer.
func (m *MockCustomerService) GetPendingAndCancelledCount(ctx context.Context, tenantID uint, branchID uint, customerID uint) ([]barberBookingPort.CountByCustomerStatus, error) {
	panic("unimplemented")
}

func (m *MockCustomerService) GetCustomerByID(ctx context.Context, tenantID, customerID uint) (*barberBookingModels.Customer, error) {
	args := m.Called(ctx, tenantID, customerID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Customer), args.Error(1)
}

func (m *MockCustomerService) CreateCustomer(ctx context.Context, cus *barberBookingModels.Customer) error {
	args := m.Called(ctx, cus)
	return args.Error(0)
}

func (m *MockCustomerService) UpdateCustomer(ctx context.Context, tenantID, customerID uint, customer *barberBookingModels.Customer) (*barberBookingModels.Customer, error) {
	args := m.Called(ctx, tenantID, customerID, customer)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Customer), args.Error(1)
}

func (m *MockCustomerService) DeleteCustomer(ctx context.Context, tenantID, customerID uint) error {
	args := m.Called(ctx, tenantID, customerID)
	return args.Error(0)
}

func (m *MockCustomerService) FindCustomerByEmail(ctx context.Context, tenantID uint, email string) (*barberBookingModels.Customer, error) {
	args := m.Called(ctx, tenantID, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Customer), args.Error(1)
}

func setupCustomerTestApp(mockSvc barberBookingPort.ICustomer) *fiber.App {
	app := fiber.New()
	controller := barberBookingControllers.NewCustomerController(mockSvc)

	app.Use(func(c *fiber.Ctx) error {
		role := c.Get("X-Mock-Role")
		if role != "" {
			c.Locals("role", role)
		}
		return c.Next()
	})
	app.Get("/tenants/:tenant_id/customers", controller.GetAllCustomers)
	app.Get("/tenants/:tenant_id/customers/:cus_id", controller.GetCustomerByID)
	app.Post("/tenants/:tenant_id/customers/find-by-email", controller.FindCustomerByEmail)
	app.Post("/tenants/:tenant_id/customers", controller.CreateCustomer)
	app.Put("/tenants/:tenant_id/customers/:cus_id", controller.UpdateCustomer)
	app.Delete("/tenants/:tenant_id/customers/:cus_id", controller.DeleteCustomer)
	app.Post("/tenants/:tenant_id/customers/find-by-email", controller.FindCustomerByEmail)
	return app
}

func TestCustomerController_GetAllCustomers(t *testing.T) {
	mockSvc := new(MockCustomerService)

	app := setupCustomerTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers", nil)
		req.Header.Set("X-Mock-Role", "USER")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Success_ShouldReturnCustomerList", func(t *testing.T) {
		mockSvc.ExpectedCalls = nil
		expected := []barberBookingModels.Customer{
			{ID: 1, TenantID: 1, Name: "John Doe"},
			{ID: 2, TenantID: 1, Name: "Jane Doe"},
		}
		mockSvc.On("GetAllCustomers", mock.Anything, uint(1)).Return(expected, nil)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[2]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                         `json:"status"`
			Message string                         `json:"message"`
			Data    []barberBookingModels.Customer `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, 2, len(body.Data))
		assert.Equal(t, "Jane Doe", body.Data[1].Name)

		mockSvc.AssertCalled(t, "GetAllCustomers", mock.Anything, uint(1))
	})

	//มีปํญหา
	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/tenants/1234/customers", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[1]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})
}
func TestGetCustomerByID(t *testing.T) {
	// Note: context is not asserted on, so mock.Anything is used instead
	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers/1", nil)
		req.Header.Set("X-Mock-Role", "USER")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("BadTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/abc/customers/1", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("BadCustomerID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers/xyz", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("GetCustomerByID", mock.Anything, uint(1), uint(5)).
			Return(nil, errors.New("db error"))

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers/5", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("GetCustomerByID", mock.Anything, uint(1), uint(6)).
			Return(nil, nil)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers/6", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnCustomer", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		want := &barberBookingModels.Customer{
			ID:       42,
			TenantID: 1,
			Name:     "Alice Example",
			Email:    "alice@example.com",
		}
		mockSvc.
			On("GetCustomerByID", mock.Anything, uint(1), uint(42)).
			Return(want, nil)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/customers/42", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                        `json:"status"`
			Message string                        `json:"message"`
			Data    *barberBookingModels.Customer `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Customer retrieved", body.Message)
		assert.Equal(t, want, body.Data)

		mockSvc.AssertExpectations(t)
	})
}

func TestCreateCustomer(t *testing.T) {
	t.Run("PermissionDenied_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers", nil)
		req.Header.Set("X-Mock-Role", "USER") // ไม่ใช่ role ที่อนุญาต
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers", strings.NewReader("invalid-json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidCustomerInput_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		body := `{"name": "     ", "phone": "123"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		payload := barberBookingModels.Customer{
			Name:  "Alice",
			Phone: "0123456789",
		}

		jsonPayload, _ := json.Marshal(payload)

		mockSvc.
			On("CreateCustomer", mock.Anything, mock.MatchedBy(func(cus *barberBookingModels.Customer) bool {
				return cus.Name == "Alice" && cus.Phone == "0123456789"
			})).
			Return(errors.New("db error"))

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers", bytes.NewReader(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturn201", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		payload := barberBookingModels.Customer{
			Name:  "Alice",
			Phone: "0123456789",
		}

		jsonPayload, _ := json.Marshal(payload)

		mockSvc.
			On("CreateCustomer", mock.Anything, mock.AnythingOfType("*barberBookingModels.Customer")).
			Return(nil)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers", bytes.NewReader(jsonPayload))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})
}

func TestUpdateCustomer(t *testing.T) {
	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/tenants/abc/customers/1", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidCustomerID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/tenants/1/customers/xyz", nil)
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/tenants/1/customers/1", strings.NewReader("invalid-json"))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidInput_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		body := `{"name": "", "phone": "123"}`
		req := httptest.NewRequest(http.MethodPut, "/tenants/1/customers/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CustomerNotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("GetCustomerByID", mock.Anything, uint(1), uint(99)).
			Return(nil, nil)

		body := `{"name": "Updated", "phone": "0123456789"}`
		req := httptest.NewRequest(http.MethodPut, "/tenants/1/customers/99", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("UpdateError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		existing := &barberBookingModels.Customer{
			ID:       1,
			TenantID: 1,
			Name:     "Old Name",
			Phone:    "0987654321",
		}

		mockSvc.
			On("GetCustomerByID", mock.Anything, uint(1), uint(1)).
			Return(existing, nil)
		mockSvc.
			On("UpdateCustomer", mock.Anything, uint(1), uint(1), mock.Anything).
			Return(nil, errors.New("update error"))

		body := `{"name": "Updated Name", "phone": "0123456789"}`
		req := httptest.NewRequest(http.MethodPut, "/tenants/1/customers/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		existing := &barberBookingModels.Customer{
			ID:       1,
			TenantID: 1,
			Name:     "Old Name",
			Phone:    "0987654321",
			Email:    "old@example.com",
		}
		updated := *existing
		updated.Name = "New Name"
		updated.Phone = "0123456789"

		mockSvc.
			On("GetCustomerByID", mock.Anything, uint(1), uint(1)).
			Return(existing, nil)
		mockSvc.
			On("UpdateCustomer", mock.Anything, uint(1), uint(1), mock.Anything).
			Return(&updated, nil)

		body := `{"name": "New Name", "phone": "0123456789", "email": "new@example.com"}`
		req := httptest.NewRequest(http.MethodPut, "/tenants/1/customers/1", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Status  string                        `json:"status"`
			Message string                        `json:"message"`
			Data    *barberBookingModels.Customer `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&result)
		assert.NoError(t, err)
		assert.Equal(t, "success", result.Status)
		assert.Equal(t, "Customer Updated", result.Message)
		assert.Equal(t, "New Name", result.Data.Name)

		mockSvc.AssertExpectations(t)
	})
}

func TestDeleteCustomer(t *testing.T) {
	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodDelete, "/tenants/1/customers/1", nil)
		req.Header.Set("X-Mock-Role", "USER") // ไม่มีสิทธิ์

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodDelete, "/tenants/abc/customers/1", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidCustomerID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodDelete, "/tenants/1/customers/xyz", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("DeleteCustomer", mock.Anything, uint(1), uint(5)).
			Return(errors.New("delete error"))

		req := httptest.NewRequest(http.MethodDelete, "/tenants/1/customers/5", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("DeleteCustomer", mock.Anything, uint(1), uint(7)).
			Return(nil)

		req := httptest.NewRequest(http.MethodDelete, "/tenants/1/customers/7", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body map[string]string
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Customer delete successfully", body["message"])

		mockSvc.AssertExpectations(t)
	})
}

func TestFindCustomerByEmail(t *testing.T) {
	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers/find-by-email", nil)
		req.Header.Set("X-Mock-Role", "USER") // ไม่มีสิทธิ์

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/abc/customers/find-by-email", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers/find-by-email", strings.NewReader("bad-json"))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("CustomerNotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("FindCustomerByEmail", mock.Anything, uint(1), "missing@example.com").
			Return(nil, nil)

		body := `{"email": "missing@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers/find-by-email", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		mockSvc.
			On("FindCustomerByEmail", mock.Anything, uint(1), "error@example.com").
			Return(nil, errors.New("db error"))

		body := `{"email": "error@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers/find-by-email", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnCustomer", func(t *testing.T) {
		mockSvc := new(MockCustomerService)
		app := setupCustomerTestApp(mockSvc)

		expected := &barberBookingModels.Customer{
			ID:       99,
			TenantID: 1,
			Name:     "Alice",
			Email:    "alice@example.com",
			Phone:    "0123456789",
		}

		mockSvc.
			On("FindCustomerByEmail", mock.Anything, uint(1), "alice@example.com").
			Return(expected, nil)

		body := `{"email": "alice@example.com"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/customers/find-by-email", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", string(barberBookingControllers.RolesCanManageCustomer[0]))

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var bodyResp struct {
			Status  string                        `json:"status"`
			Message string                        `json:"message"`
			Data    *barberBookingModels.Customer `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&bodyResp)
		assert.NoError(t, err)
		assert.Equal(t, "success", bodyResp.Status)
		assert.Equal(t, "Customer retrieved", bodyResp.Message)
		assert.Equal(t, expected.Email, bodyResp.Data.Email)

		mockSvc.AssertExpectations(t)
	})
}
