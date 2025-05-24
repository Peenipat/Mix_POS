package Core_Controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	coreControllers "myapp/modules/core/controllers"
	coreModels "myapp/modules/core/models"
	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"
)


type MockTenantService struct {
	mock.Mock
}

func (m *MockTenantService) DeleteTenant(ctx context.Context, id uint) error {
    args := m.Called(ctx, id)
    return args.Error(0)
}

func (m *MockTenantService) GetTenantByID(ctx context.Context, id uint) (*coreModels.Tenant, error) {
    args := m.Called(ctx, id)
    tenant, _ := args.Get(0).(*coreModels.Tenant)
    return tenant, args.Error(1)
}


func (m *MockTenantService) ListTenants(ctx context.Context, onlyActive bool) ([]coreModels.Tenant, error) {
    args := m.Called(ctx, onlyActive)
    if args.Get(0) == nil {
        return nil, args.Error(1)
    }
    return args.Get(0).([]coreModels.Tenant), args.Error(1)
}

func (m *MockTenantService) UpdateTenant(ctx context.Context, input corePort.UpdateTenantInput) error {
    args := m.Called(ctx, input)
    return args.Error(0)
}

func (m *MockTenantService) CreateTenant(ctx context.Context, input corePort.CreateTenantInput) (*coreModels.Tenant, error) {
	args := m.Called(ctx, input)
	tenant, _ := args.Get(0).(*coreModels.Tenant)
	return tenant, args.Error(1)
}

// Setup app with CreateTenant route
func setupTenantApp(svc *MockTenantService) *fiber.App {
	app := fiber.New()
	ctrl := coreControllers.NewTenantController(svc)
	app.Post("/tenants", ctrl.CreateTenant)
	app.Get("/tenants/:id", ctrl.GetTenantByID)
	app.Get("/tenants", ctrl.ListTenants)
	app.Put("/tenants/:id", ctrl.UpdateTenant)
	app.Delete("/tenants/:id", ctrl.DeleteTenant)
	return app
}

func TestCreateTenantController(t *testing.T) {
	t.Run("MalformedJSON", func(t *testing.T) {
		svc := new(MockTenantService)
		app := setupTenantApp(svc)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewReader([]byte(`{"name":`)))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidInput", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("CreateTenant", mock.Anything, corePort.CreateTenantInput{Name: "", Domain: "d"}).
			Return(nil, coreServices.ErrInvalidTenantInput).
			Once()

		app := setupTenantApp(svc)
		body := map[string]string{"name": "", "domain": "d"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("DomainTaken", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("CreateTenant", mock.Anything, corePort.CreateTenantInput{Name: "N", Domain: "d"}).
			Return(nil, coreServices.ErrDomainTaken).
			Once()

		app := setupTenantApp(svc)
		body := map[string]string{"name": "N", "domain": "d"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("CreateTenant", mock.Anything, corePort.CreateTenantInput{Name: "N", Domain: "d"}).
			Return(nil, errors.New("db down")).
			Once()

		app := setupTenantApp(svc)
		body := map[string]string{"name": "N", "domain": "d"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		svc := new(MockTenantService)
		expected := &coreModels.Tenant{ID: 42, Name: "Acme", Domain: "acme.local", IsActive: true}
		svc.
			On("CreateTenant", mock.Anything, corePort.CreateTenantInput{Name: "Acme", Domain: "acme.local"}).
			Return(expected, nil).
			Once()

		app := setupTenantApp(svc)
		body := map[string]string{"name": "Acme", "domain": "acme.local"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPost, "/tenants", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respBody struct {
			Status string            `json:"status"`
			Data   coreModels.Tenant `json:"data"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&respBody))
		assert.Equal(t, "success", respBody.Status)
		assert.Equal(t, expected.ID, respBody.Data.ID)
		assert.Equal(t, expected.Name, respBody.Data.Name)
		assert.Equal(t, expected.Domain, respBody.Data.Domain)

		svc.AssertExpectations(t)
	})
}

func TestGetTenantByIDController(t *testing.T) {
    t.Run("InvalidID_Format", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodGet, "/tenants/abc", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidID_Zero", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodGet, "/tenants/0", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ServiceError_InvalidID", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("GetTenantByID", mock.Anything, uint(5)).
            Return((*coreModels.Tenant)(nil), coreServices.ErrInvalidTenantID).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants/5", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_NotFound", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("GetTenantByID", mock.Anything, uint(7)).
            Return((*coreModels.Tenant)(nil), coreServices.ErrTenantNotFound).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants/7", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_Internal", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("GetTenantByID", mock.Anything, uint(9)).
            Return((*coreModels.Tenant)(nil), errors.New("db down")).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants/9", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        svc.AssertExpectations(t)
    })

    t.Run("Success", func(t *testing.T) {
        expected := &coreModels.Tenant{
            ID:       11,
            Name:     "TestTenant",
            Domain:   "test.local",
            IsActive: true,
        }
        svc := new(MockTenantService)
        svc.
            On("GetTenantByID", mock.Anything, uint(11)).
            Return(expected, nil).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants/11", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string            `json:"status"`
            Data   coreModels.Tenant `json:"data"`
        }
        err := json.NewDecoder(resp.Body).Decode(&body)
        require.NoError(t, err)

        assert.Equal(t, "success", body.Status)
        assert.Equal(t, expected.ID, body.Data.ID)
        assert.Equal(t, expected.Name, body.Data.Name)
        assert.Equal(t, expected.Domain, body.Data.Domain)
        assert.Equal(t, expected.IsActive, body.Data.IsActive)

        svc.AssertExpectations(t)
    })
}

func TestListTenantsController(t *testing.T) {
    ctx := mock.Anything // we match any context

    t.Run("InvalidActiveQuery", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodGet, "/tenants?active=notbool", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ServiceError", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.On("ListTenants", ctx, true).
            Return(nil, errors.New("db failure")).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        svc.AssertExpectations(t)
    })

    t.Run("Success_DefaultActive", func(t *testing.T) {
        expected := []coreModels.Tenant{
            {ID: 1, Name: "A", Domain: "a.local"},
        }
        svc := new(MockTenantService)
        svc.On("ListTenants", ctx, true).
            Return(expected, nil).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string                 `json:"status"`
            Data   []coreModels.Tenant    `json:"data"`
        }
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "success", body.Status)
        assert.Equal(t, expected, body.Data)
        svc.AssertExpectations(t)
    })

    t.Run("Success_ActiveFalse", func(t *testing.T) {
        expected := []coreModels.Tenant{
            {ID: 2, Name: "B", Domain: "b.local"},
            {ID: 3, Name: "C", Domain: "c.local"},
        }
        svc := new(MockTenantService)
        svc.On("ListTenants", ctx, false).
            Return(expected, nil).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodGet, "/tenants?active=false", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string                 `json:"status"`
            Data   []coreModels.Tenant    `json:"data"`
        }
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "success", body.Status)
        assert.Equal(t, expected, body.Data)
        svc.AssertExpectations(t)
    })
}

func TestUpdateTenantController(t *testing.T) {
    validInput := map[string]interface{}{
        "name":      "New Name",
        "domain":    "new.domain",
        "is_active": false,
    }
    validBody, _ := json.Marshal(validInput)


    t.Run("InvalidID_Format", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodPut, "/tenants/abc", bytes.NewReader(validBody))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidID_Zero", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodPut, "/tenants/0", bytes.NewReader(validBody))
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("MalformedJSON", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodPut, "/tenants/1", bytes.NewReader([]byte(`{"name":`)))
        req.Header.Set("Content-Type", "application/json")
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ServiceError_InvalidID", func(t *testing.T) {
		svc := new(MockTenantService)
		// Expect UpdateTenant called with matching ID, Name and Domain – ignore IsActive
		svc.
			On("UpdateTenant", mock.Anything, mock.MatchedBy(func(in corePort.UpdateTenantInput) bool {
				return in.ID == 5 &&
					in.Name != nil && *in.Name == "New Name" &&
					in.Domain != nil && *in.Domain == "new.domain"
			})).
			Return(coreServices.ErrInvalidTenantID).
			Once()
	
		app := setupTenantApp(svc)
		// Build a body that omits is_active so IsActive remains nil
		body := map[string]interface{}{
			"name":   "New Name",
			"domain": "new.domain",
		}
		buf, _ := json.Marshal(body)
	
		req := httptest.NewRequest(http.MethodPut, "/tenants/5", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
	
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})

    t.Run("ServiceError_InvalidInput", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("UpdateTenant", mock.Anything, mock.MatchedBy(func(in corePort.UpdateTenantInput) bool {
				return in.ID == 7 &&
					in.Name != nil && *in.Name == "New Name" &&
					in.Domain != nil && *in.Domain == "new.domain"
				// ไม่ต้องเช็ค in.IsActive
			})).
			Return(coreServices.ErrInvalidTenantInput).
			Once()
	
		app := setupTenantApp(svc)
		body := map[string]interface{}{
			"name":   "New Name",
			"domain": "new.domain",
			// omit "is_active"
		}
		buf, _ := json.Marshal(body)
	
		req := httptest.NewRequest(http.MethodPut, "/tenants/7", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
	
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})
	
	t.Run("ServiceError_DomainTaken", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("UpdateTenant", mock.Anything, mock.MatchedBy(func(in corePort.UpdateTenantInput) bool {
				return in.ID == 9 &&
					in.Name != nil && *in.Name == "New Name" &&
					in.Domain != nil && *in.Domain == "new.domain"
			})).
			Return(coreServices.ErrDomainTaken).
			Once()
	
		app := setupTenantApp(svc)
		body := map[string]interface{}{
			"name":   "New Name",
			"domain": "new.domain",
		}
		buf, _ := json.Marshal(body)
	
		req := httptest.NewRequest(http.MethodPut, "/tenants/9", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
	
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})
	
	t.Run("ServiceError_NotFound", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("UpdateTenant", mock.Anything, mock.MatchedBy(func(in corePort.UpdateTenantInput) bool {
				return in.ID == 11 &&
					in.Name != nil && *in.Name == "New Name" &&
					in.Domain != nil && *in.Domain == "new.domain"
			})).
			Return(coreServices.ErrTenantNotFound).
			Once()
	
		app := setupTenantApp(svc)
		body := map[string]interface{}{
			"name":   "New Name",
			"domain": "new.domain",
		}
		buf, _ := json.Marshal(body)
	
		req := httptest.NewRequest(http.MethodPut, "/tenants/11", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
	
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Internal", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("UpdateTenant", mock.Anything, mock.MatchedBy(func(in corePort.UpdateTenantInput) bool {
				return in.ID == 13 &&
					in.Name != nil && *in.Name == "New Name" &&
					in.Domain != nil && *in.Domain == "new.domain"
			})).
			Return(errors.New("db down")).
			Once()
	
		app := setupTenantApp(svc)
		req := httptest.NewRequest(http.MethodPut, "/tenants/13", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})
	
	t.Run("Success", func(t *testing.T) {
		svc := new(MockTenantService)
		svc.
			On("UpdateTenant", mock.Anything, mock.MatchedBy(func(in corePort.UpdateTenantInput) bool {
				return in.ID == 15 &&
					in.Name != nil && *in.Name == "New Name" &&
					in.Domain != nil && *in.Domain == "new.domain"
			})).
			Return(nil).
			Once()
	
		app := setupTenantApp(svc)
		req := httptest.NewRequest(http.MethodPut, "/tenants/15", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
	
		var body struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Tenant updated", body.Message)
	
		svc.AssertExpectations(t)
	})
}

func TestDeleteTenantController(t *testing.T) {
    ctxMatcher := mock.Anything

    t.Run("InvalidID_Format", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodDelete, "/tenants/abc", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidID_Zero", func(t *testing.T) {
        svc := new(MockTenantService)
        app := setupTenantApp(svc)

        req := httptest.NewRequest(http.MethodDelete, "/tenants/0", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ServiceError_InvalidID", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("DeleteTenant", ctxMatcher, uint(5)).
            Return(coreServices.ErrInvalidTenantID).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/5", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, coreServices.ErrInvalidTenantID.Error(), body["message"])

        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_NotFound", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("DeleteTenant", ctxMatcher, uint(7)).
            Return(coreServices.ErrTenantNotFound).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/7", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "Tenant not found", body["message"])

        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_InUse", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("DeleteTenant", ctxMatcher, uint(9)).
            Return(coreServices.ErrTenantInUse).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/9", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusConflict, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "Cannot delete tenant: still in use", body["message"])

        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_Internal", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("DeleteTenant", ctxMatcher, uint(11)).
            Return(errors.New("db down")).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/11", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "Failed to delete tenant", body["message"])

        svc.AssertExpectations(t)
    })

    t.Run("Success", func(t *testing.T) {
        svc := new(MockTenantService)
        svc.
            On("DeleteTenant", ctxMatcher, uint(13)).
            Return(nil).
            Once()

        app := setupTenantApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/13", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "success", body["status"])
        assert.Equal(t, "Tenant deleted", body["message"])

        svc.AssertExpectations(t)
    })
}

