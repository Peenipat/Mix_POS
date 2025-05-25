package Core_Controllers_test

import (
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
	coreServices "myapp/modules/core/services"
)

// MockTenantUserService mocks AddUserToTenant
type MockTenantUserService struct {
	mock.Mock
}

// IsUserInTenant implements corePort.ITenantUser.
func (m *MockTenantUserService) IsUserInTenant(ctx context.Context, tenantID uint, userID uint) (bool, error) {
	panic("unimplemented")
}

func (m *MockTenantUserService) ListTenantsByUser(ctx context.Context, userID uint) ([]coreModels.Tenant, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]coreModels.Tenant), args.Error(1)
}

// ListUsersByTenant implements corePort.ITenantUser.
func (m *MockTenantUserService) ListUsersByTenant(ctx context.Context, tenantID uint) ([]coreModels.User, error) {
	panic("unimplemented")
}

func (m *MockTenantUserService) RemoveUserFromTenant(ctx context.Context, tenantID, userID uint) error {
    args := m.Called(ctx, tenantID, userID)
    return args.Error(0)
}

func (m *MockTenantUserService) AddUserToTenant(ctx context.Context, tenantID, userID uint) error {
	args := m.Called(ctx, tenantID, userID)
	return args.Error(0)
}

// setupApp registers the POST route
func setupApp(svc *MockTenantUserService) *fiber.App {
	app := fiber.New()
	ctrl := coreControllers.NewTenantUserController(svc)
	app.Post("/tenants/:tenant_id/users/:user_id", ctrl.AddUserToTenant)
	app.Delete("/tenants/:tenant_id/users/:user_id", ctrl.RemoveUserFromTenant)
	app.Get("/users/:user_id/tenants", ctrl.ListTenantsByUser)
	return app
}

func TestAddUserToTenantController(t *testing.T) {
	ctxMatcher := mock.Anything

	t.Run("InvalidTenantID_Format", func(t *testing.T) {
		svc := new(MockTenantUserService)
		app := setupApp(svc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/abc/users/1", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidUserID_Format", func(t *testing.T) {
		svc := new(MockTenantUserService)
		app := setupApp(svc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/users/xyz", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_InvalidTenantID", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", ctxMatcher, uint(5), uint(2)).
			Return(coreServices.ErrInvalidTenantID).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/5/users/2", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, coreServices.ErrInvalidTenantID.Error(), body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_InvalidUserID", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", mock.Anything, uint(3), uint(7)).
			Return(coreServices.ErrInvalidUserID).
			Once()
	
		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/3/users/7", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	
		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, coreServices.ErrInvalidUserID.Error(), body["message"])
	
		svc.AssertExpectations(t)
	})
	

	t.Run("ServiceError_TenantNotFound", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", ctxMatcher, uint(7), uint(2)).
			Return(coreServices.ErrTenantNotFound).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/7/users/2", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "Tenant not found", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_UserNotFound", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", ctxMatcher, uint(1), uint(99)).
			Return(coreServices.ErrUserNotFound).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/users/99", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "User not found", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_AlreadyAssigned", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", ctxMatcher, uint(4), uint(5)).
			Return(coreServices.ErrUserAlreadyAssigned).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/4/users/5", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusConflict, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "User already assigned to this tenant", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Internal", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", ctxMatcher, uint(8), uint(9)).
			Return(errors.New("db down")).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/8/users/9", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "Failed to assign user to tenant", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("AddUserToTenant", ctxMatcher, uint(2), uint(3)).
			Return(nil).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodPost, "/tenants/2/users/3", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User assigned to tenant", body["message"])

		svc.AssertExpectations(t)
	})
}

func TestRemoveUserFromTenantController(t *testing.T) {
    ctxMatcher := mock.Anything

    t.Run("InvalidTenantID", func(t *testing.T) {
        app := setupApp(new(MockTenantUserService))
        req := httptest.NewRequest(http.MethodDelete, "/tenants/abc/users/1", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidUserID", func(t *testing.T) {
        app := setupApp(new(MockTenantUserService))
        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/users/xyz", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ServiceError_InvalidIDs", func(t *testing.T) {
        svc := new(MockTenantUserService)
        svc.
            On("RemoveUserFromTenant", ctxMatcher, uint(5), uint(7)).
            Return(coreServices.ErrInvalidTenantID).
            Once()

        app := setupApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/5/users/7", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, coreServices.ErrInvalidTenantID.Error(), body["message"])
        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_TenantNotFound", func(t *testing.T) {
        svc := new(MockTenantUserService)
        svc.
            On("RemoveUserFromTenant", ctxMatcher, uint(2), uint(3)).
            Return(coreServices.ErrTenantNotFound).
            Once()

        app := setupApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/2/users/3", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "Tenant not found", body["message"])
        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_UserNotFound", func(t *testing.T) {
        svc := new(MockTenantUserService)
        svc.
            On("RemoveUserFromTenant", ctxMatcher, uint(1), uint(9)).
            Return(coreServices.ErrUserNotFound).
            Once()

        app := setupApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/1/users/9", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "User not found", body["message"])
        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_NotAssigned", func(t *testing.T) {
        svc := new(MockTenantUserService)
        svc.
            On("RemoveUserFromTenant", ctxMatcher, uint(4), uint(5)).
            Return(coreServices.ErrUserNotAssigned).
            Once()

        app := setupApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/4/users/5", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusConflict, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "User is not assigned to this tenant", body["message"])
        svc.AssertExpectations(t)
    })

    t.Run("ServiceError_Internal", func(t *testing.T) {
        svc := new(MockTenantUserService)
        svc.
            On("RemoveUserFromTenant", ctxMatcher, uint(8), uint(9)).
            Return(errors.New("db down")).
            Once()

        app := setupApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/8/users/9", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "error", body["status"])
        assert.Equal(t, "Failed to remove user from tenant", body["message"])
        svc.AssertExpectations(t)
    })

    t.Run("Success", func(t *testing.T) {
        svc := new(MockTenantUserService)
        svc.
            On("RemoveUserFromTenant", ctxMatcher, uint(7), uint(2)).
            Return(nil).
            Once()

        app := setupApp(svc)
        req := httptest.NewRequest(http.MethodDelete, "/tenants/7/users/2", nil)
        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body map[string]string
        require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
        assert.Equal(t, "success", body["status"])
        assert.Equal(t, "User removed from tenant", body["message"])
        svc.AssertExpectations(t)
    })
}


func TestListTenantsByUserController(t *testing.T) {
	ctxMatcher := mock.Anything

	t.Run("InvalidUserID", func(t *testing.T) {
		app := setupApp(new(MockTenantUserService))
		req := httptest.NewRequest(http.MethodGet, "/users/abc/tenants", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_InvalidID", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("ListTenantsByUser", ctxMatcher, uint(5)).
			Return([]coreModels.Tenant{}, coreServices.ErrInvalidUserID).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodGet, "/users/5/tenants", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, coreServices.ErrInvalidUserID.Error(), body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_UserNotFound", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("ListTenantsByUser", ctxMatcher, uint(2)).
			Return([]coreModels.Tenant{}, coreServices.ErrUserNotFound).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodGet, "/users/2/tenants", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "User not found", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_NoTenants", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("ListTenantsByUser", ctxMatcher, uint(3)).
			Return([]coreModels.Tenant{}, coreServices.ErrNoTenantsAssigned).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodGet, "/users/3/tenants", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "No tenants assigned to this user", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Internal", func(t *testing.T) {
		svc := new(MockTenantUserService)
		svc.
			On("ListTenantsByUser", ctxMatcher, uint(4)).
			Return([]coreModels.Tenant{}, errors.New("db error")).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodGet, "/users/4/tenants", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
		assert.Equal(t, "error", body["status"])
		assert.Equal(t, "Failed to fetch tenants", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		tenants := []coreModels.Tenant{
			{ID: 10, Name: "A", Domain: "a"},
			{ID: 20, Name: "B", Domain: "b"},
		}
		svc := new(MockTenantUserService)
		svc.
			On("ListTenantsByUser", ctxMatcher, uint(7)).
			Return(tenants, nil).
			Once()

		app := setupApp(svc)
		req := httptest.NewRequest(http.MethodGet, "/users/7/tenants", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var result struct {
			Status string              `json:"status"`
			Data   []coreModels.Tenant `json:"data"`
		}
		require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))
		assert.Equal(t, "success", result.Status)
		assert.Len(t, result.Data, 2)
		assert.Equal(t, uint(10), result.Data[0].ID)
		assert.Equal(t, uint(20), result.Data[1].ID)

		svc.AssertExpectations(t)
	})
}








