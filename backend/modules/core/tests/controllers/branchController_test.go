package Core_Controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	 "gorm.io/gorm"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	coreControllers "myapp/modules/core/controllers"
	coreModels "myapp/modules/core/models"
	coreServices "myapp/modules/core/services"
)

// Mock ของ BranchService interface
type MockBranchService struct {
	mock.Mock
}

func (m *MockBranchService) DeleteBranch(id uint) error {
    args := m.Called(id)
    return args.Error(0)
}

func (m *MockBranchService) UpdateBranch(b *coreModels.Branch) error {
    args := m.Called(b)
    return args.Error(0)
}

func (m *MockBranchService) GetBranchByID(id uint) (*coreModels.Branch, error) {
	args := m.Called(id)
	branch, _ := args.Get(0).(*coreModels.Branch)
	return branch, args.Error(1)
}

func (m *MockBranchService) CreateBranch(b *coreModels.Branch) error {
	args := m.Called(b)
	return args.Error(0)
}

func (m *MockBranchService) GetAllBranches(tenantID uint) ([]coreModels.Branch, error) {
	args := m.Called(tenantID)
	return args.Get(0).([]coreModels.Branch), args.Error(1)
}

func (m *MockBranchService) GetBranchesByTenantID(tenantID uint) ([]coreModels.Branch, error) {
    args := m.Called(tenantID)
    return args.Get(0).([]coreModels.Branch), args.Error(1)
}

func setupBranchApp(svc *MockBranchService) *fiber.App {
	app := fiber.New()

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", c.Get("X-Role"))
		if tid := c.Get("X-Tenant-ID"); tid != "" {
			c.Locals("tenant_id", uint(123))
		}
		return c.Next()
	})
	ctrl := coreControllers.NewBranchController(svc)

	// 3. ผูก route
	app.Post("/branches", ctrl.CreateBranch)
	app.Get("/branches", ctrl.GetBranches)
	app.Get("/branches/by-tenant", ctrl.GetBranchesByTenantID)
	app.Get("/branches/:id", ctrl.GetBranchByID)
	app.Put("/branches/:id", ctrl.UpdateBranch)
	app.Delete("/branches/:id", ctrl.DeleteBranch)
	
	return app
}

var RolesCanManageBranchTest = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameTenant,
}

func TestCreateBranch(t *testing.T) {
	type payload struct {
		Name string `json:"name"`
	}

	valid := payload{Name: "New Branch"}
	validBody, _ := json.Marshal(valid)

	t.Run("Success", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		mockSvc.
			On("CreateBranch", mock.MatchedBy(func(b *coreModels.Branch) bool {
				return b.Name == valid.Name && b.TenantID == 123
			})).
			Return(nil).
			Once()

		app := setupBranchApp(mockSvc)
		req := httptest.NewRequest("POST", "/branches", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", string(RolesCanManageBranchTest[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Unauthorized", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		app := setupBranchApp(mockSvc)

		req := httptest.NewRequest("POST", "/branches", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", "STAFF") // assume not in RolesCanManageBranch
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		app := setupBranchApp(mockSvc)

		req := httptest.NewRequest("POST", "/branches", strings.NewReader(`{"name":`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", string(RolesCanManageBranchTest[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Missing TenantID", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		app := setupBranchApp(mockSvc)

		req := httptest.NewRequest("POST", "/branches", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", string(RolesCanManageBranchTest[0]))

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Duplicate Name", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		mockSvc.
			On("CreateBranch", mock.Anything).
			Return(errors.New("branch name already exists for this tenant")).
			Once()

		app := setupBranchApp(mockSvc)
		req := httptest.NewRequest("POST", "/branches", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", string(RolesCanManageBranchTest[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		mockSvc.
			On("CreateBranch", mock.Anything).
			Return(errors.New("db down")).
			Once()

		app := setupBranchApp(mockSvc)
		req := httptest.NewRequest("POST", "/branches", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", string(RolesCanManageBranchTest[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
	})
}

func TestGetBranches(t *testing.T) {
	mockSvc := new(MockBranchService)
	app := setupBranchApp(mockSvc)

	type payload struct {
		Name string `json:"name"`
	}

	valid := payload{Name: "New Branch"}
	validBody, _ := json.Marshal(valid)

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/branches", nil)
		req.Header.Set("X-Role", "STAFF") // not allowed
		req.Header.Set("X-Tenant-ID", "123")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("MissingTenantID", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/branches", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("TenantNotFound", func(t *testing.T) {
		mockSvc.
			On("GetAllBranches", uint(123)).
			Return(make([]coreModels.Branch, 0), coreServices.ErrTenantNotFound).
			Once()
		req := httptest.NewRequest("GET", "/branches", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Service Error", func(t *testing.T) {
		mockSvc := new(MockBranchService)
		mockSvc.
			On("CreateBranch", mock.Anything).
			Return(errors.New("db down")).
			Once()

		app := setupBranchApp(mockSvc)
		req := httptest.NewRequest("POST", "/branches", bytes.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", string(RolesCanManageBranchTest[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_Empty", func(t *testing.T) {
		mockSvc.
			On("GetAllBranches", uint(123)).
			Return([]coreModels.Branch{}, nil).
			Once()
		req := httptest.NewRequest("GET", "/branches", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string              `json:"status"`
			Data   []coreModels.Branch `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "success", body.Status)
		assert.Len(t, body.Data, 0)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_NonEmpty", func(t *testing.T) {
		expected := []coreModels.Branch{
			{ID: 1, TenantID: 123, Name: "A"},
			{ID: 2, TenantID: 123, Name: "B"},
		}
		mockSvc.
			On("GetAllBranches", uint(123)).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest("GET", "/branches", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string              `json:"status"`
			Data   []coreModels.Branch `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected, body.Data)

		mockSvc.AssertExpectations(t)
	})
}

func TestGetBranchController(t *testing.T) {
	app := setupBranchApp(new(MockBranchService))

	t.Run("Unauthorized", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/branches/1", nil)
		req.Header.Set("X-Role", "STAFF")
		req.Header.Set("X-Tenant-ID", "123")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidID_Format", func(t *testing.T) {
		svc := new(MockBranchService)
		app := setupBranchApp(svc)
		req := httptest.NewRequest("GET", "/branches/abc", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidID_Zero", func(t *testing.T) {
		svc := new(MockBranchService)
		app := setupBranchApp(svc)
		req := httptest.NewRequest("GET", "/branches/0", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_InvalidID", func(t *testing.T) {
		svc := new(MockBranchService)
		svc.On("GetBranchByID", uint(5)).Return(nil, coreServices.ErrInvalidID).Once()
		app := setupBranchApp(svc)
		req := httptest.NewRequest("GET", "/branches/5", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_NotFound", func(t *testing.T) {
		svc := new(MockBranchService)
		svc.On("GetBranchByID", uint(7)).Return(nil, coreServices.ErrBranchNotFound).Once()
		app := setupBranchApp(svc)
		req := httptest.NewRequest("GET", "/branches/7", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Other", func(t *testing.T) {
		svc := new(MockBranchService)
		svc.On("GetBranchByID", uint(9)).Return(nil, errors.New("db down")).Once()
		app := setupBranchApp(svc)
		req := httptest.NewRequest("GET", "/branches/9", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		expected := &coreModels.Branch{ID: 11, TenantID: 123, Name: "HQ"}
		svc := new(MockBranchService)
		svc.On("GetBranchByID", uint(11)).Return(expected, nil).Once()
		app := setupBranchApp(svc)

		req := httptest.NewRequest("GET", "/branches/11", nil)
		req.Header.Set("X-Role", string(coreControllers.RolesCanManageBranch[0]))
		req.Header.Set("X-Tenant-ID", "123")

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string            `json:"status"`
			Data   coreModels.Branch `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, *expected, body.Data)

		svc.AssertExpectations(t)
	})
}

func TestUpdateBranch(t *testing.T) {
    type payload struct {
        Name string `json:"name"`
    }
    valid := payload{Name: "Updated Name"}
    validBody, _ := json.Marshal(valid)

    roles := []string{
        string(coreModels.RoleNameSaaSSuperAdmin),
        string(coreModels.RoleNameTenantAdmin),
    }

    t.Run("Success", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        // expectation: will receive branch.ID=42, Name="Updated Name"
        mockSvc.
            On("UpdateBranch", mock.MatchedBy(func(b *coreModels.Branch) bool {
                return b.ID == 42 && b.Name == valid.Name
            })).
            Return(nil).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/branches/42", bytes.NewReader(validBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, err := app.Test(req)
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string            `json:"status"`
            Data   coreModels.Branch `json:"data"`
        }
        json.NewDecoder(resp.Body).Decode(&body)
        assert.Equal(t, "success", body.Status)
        assert.Equal(t, uint(42), body.Data.ID)
        assert.Equal(t, valid.Name, body.Data.Name)

        mockSvc.AssertExpectations(t)
    })

    t.Run("Unauthorized", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodPut, "/branches/42", bytes.NewReader(validBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", "STAFF")
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusForbidden, resp.StatusCode)
    })

    t.Run("InvalidID_Format", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodPut, "/branches/abc", bytes.NewReader(validBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidID_Zero", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodPut, "/branches/0", bytes.NewReader(validBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("Malformed JSON", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodPut, "/branches/42", strings.NewReader(`{"name":`))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("ValidationError_NameRequired", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        // service will reject empty name
        mockSvc.
            On("UpdateBranch", mock.MatchedBy(func(b *coreModels.Branch) bool {
                return b.ID == 42 && strings.TrimSpace(b.Name) == ""
            })).
            Return(errors.New("branch name is required")).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/branches/42", strings.NewReader(`{"name":""}`))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("NotFound", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("UpdateBranch", mock.Anything).
            Return(gorm.ErrRecordNotFound).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/branches/42", bytes.NewReader(validBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("ServiceError", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("UpdateBranch", mock.Anything).
            Return(errors.New("db down")).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodPut, "/branches/42", bytes.NewReader(validBody))
        req.Header.Set("Content-Type", "application/json")
        req.Header.Set("X-Role", roles[0])
        req.Header.Set("X-Tenant-ID", "123")

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })
}

func TestDeleteBranchController(t *testing.T) {
    validRole := string(RolesCanManageBranchTest[0])
    tenantHeader := "123"

    t.Run("Success", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("DeleteBranch", uint(42)).
            Return(nil).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodDelete, "/branches/42", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, err := app.Test(req)
        assert.NoError(t, err)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        mockSvc.AssertExpectations(t)
    })

    t.Run("Unauthorized", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/branches/42", nil)
        req.Header.Set("X-Role", "STAFF")
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusForbidden, resp.StatusCode)
    })

    t.Run("InvalidID_Format", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/branches/abc", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidID_Zero", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodDelete, "/branches/0", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("NotFound", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("DeleteBranch", uint(7)).
            Return(coreServices.ErrBranchNotFound).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodDelete, "/branches/7", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("InUse", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("DeleteBranch", uint(9)).
            Return(coreServices.ErrBranchInUse).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodDelete, "/branches/9", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusConflict, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("ServiceError", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("DeleteBranch", uint(11)).
            Return(errors.New("db down")).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodDelete, "/branches/11", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })
}

func TestGetBranchesByTenantID_Controller(t *testing.T) {
    validRole := string(RolesCanManageBranchTest[0])
    tenantHeader := "123"

    t.Run("Unauthorized", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", "STAFF")
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusForbidden, resp.StatusCode)
    })

    t.Run("MissingTenantID", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        app := setupBranchApp(mockSvc)

        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", validRole)
        // no X-Tenant-ID

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
    })

    t.Run("InvalidTenantID", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("GetBranchesByTenantID", uint(123)).
            Return([]coreModels.Branch{}, coreServices.ErrInvalidTenantID).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("TenantNotFound", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("GetBranchesByTenantID", uint(123)).
            Return([]coreModels.Branch{}, coreServices.ErrTenantNotFound).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusNotFound, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })

    t.Run("Success_Empty", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("GetBranchesByTenantID", uint(123)).
            Return([]coreModels.Branch{}, nil).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string              `json:"status"`
            Data   []coreModels.Branch `json:"data"`
        }
        json.NewDecoder(resp.Body).Decode(&body)
        assert.Equal(t, "success", body.Status)
        assert.Len(t, body.Data, 0)

        mockSvc.AssertExpectations(t)
    })

    t.Run("Success_NonEmpty", func(t *testing.T) {
        expected := []coreModels.Branch{
            {ID: 1, TenantID: 123, Name: "A"},
            {ID: 2, TenantID: 123, Name: "B"},
        }
        mockSvc := new(MockBranchService)
        mockSvc.
            On("GetBranchesByTenantID", uint(123)).
            Return(expected, nil).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusOK, resp.StatusCode)

        var body struct {
            Status string              `json:"status"`
            Data   []coreModels.Branch `json:"data"`
        }
        json.NewDecoder(resp.Body).Decode(&body)
        assert.Equal(t, "success", body.Status)
        assert.Equal(t, expected, body.Data)

        mockSvc.AssertExpectations(t)
    })

    t.Run("ServiceError", func(t *testing.T) {
        mockSvc := new(MockBranchService)
        mockSvc.
            On("GetBranchesByTenantID", uint(123)).
            Return([]coreModels.Branch{}, errors.New("db down")).
            Once()

        app := setupBranchApp(mockSvc)
        req := httptest.NewRequest(http.MethodGet, "/branches/by-tenant", nil)
        req.Header.Set("X-Role", validRole)
        req.Header.Set("X-Tenant-ID", tenantHeader)

        resp, _ := app.Test(req)
        assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
        mockSvc.AssertExpectations(t)
    })
}
