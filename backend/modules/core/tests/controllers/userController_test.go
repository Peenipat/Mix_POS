package Core_Controllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
    "strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
	"github.com/stretchr/testify/mock"

	coreControllers "myapp/modules/core/controllers"
	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"
)

// Mock implementation of IUserService
type MockUserService struct {
	mock.Mock
    MeFunc func(ctx context.Context, userID uint) (*corePort.MeDTO, error)
}

// Me implements corePort.IUser.
func (m *MockUserService) Me(ctx context.Context, userID uint) (*corePort.MeDTO, error) {
    return m.MeFunc(ctx, userID)
}

func (m *MockUserService) CreateUserFromRegister(input corePort.RegisterInput) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserService) CreateUserFromAdmin(input corePort.CreateUserInput) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserService) ChangeRoleFromAdmin(input corePort.ChangeRoleInput) error {
	args := m.Called(input)
	return args.Error(0)
}

func (m *MockUserService) GetAllUsers(limit int, offset int) ([]corePort.UserInfoResponse, error) {
	args := m.Called(limit, offset)
	// Get(0) คือ slice, Error(1) คือ error
	return args.Get(0).([]corePort.UserInfoResponse), args.Error(1)
}

func (m *MockUserService) FilterUsersByRole(role string) ([]corePort.UserInfoResponse, error) {
	args := m.Called(role)
	return args.Get(0).([]corePort.UserInfoResponse), args.Error(1)
}

func (m *MockUserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	args := m.Called(ctx, userID, oldPassword, newPassword)
	return args.Error(0)
}

func setupUserApp(svc corePort.IUser) *fiber.App {
    app := fiber.New()
    ctrl := coreControllers.NewUserController(svc)

    // inject user_id และ role ลงใน Locals
    app.Use(func(c *fiber.Ctx) error {
        if h := c.Get("X-User-ID"); h != "" {
            if uid, err := strconv.ParseUint(h, 10, 32); err == nil {
                c.Locals("user_id", uint(uid))
            }
        }
        c.Locals("role", c.Get("X-Role"))
        return c.Next()
    })

    app.Put("/users/:id/password", ctrl.ChangePassword)
    app.Get("/auth/me", ctrl.Me)
    return app
}


func TestChangePasswordController(t *testing.T) {
	validRole := "USER"

	t.Run("InvalidID_Format", func(t *testing.T) {
		svc := new(MockUserService)
		app := setupUserApp(svc)

		req := httptest.NewRequest(http.MethodPut, "/users/abc/password", nil)
		req.Header.Set("X-Role", validRole)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidID_Zero", func(t *testing.T) {
		svc := new(MockUserService)
		app := setupUserApp(svc)

		req := httptest.NewRequest(http.MethodPut, "/users/0/password", nil)
		req.Header.Set("X-Role", validRole)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		svc := new(MockUserService)
		svc.
			On("ChangePassword", mock.Anything, uint(42), "old", "new").
			Return(coreServices.ErrUserNotFound).
			Once()

		app := setupUserApp(svc)
		body := map[string]string{"old_password": "old", "new_password": "new"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/42/password", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", validRole)

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("InvalidOldPassword", func(t *testing.T) {
		svc := new(MockUserService)
		svc.
			On("ChangePassword", mock.Anything, uint(7), "wrong", "new").
			Return(coreServices.ErrInvalidOldPassword).
			Once()

		app := setupUserApp(svc)
		body := map[string]string{"old_password": "wrong", "new_password": "new"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/7/password", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", validRole)

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("ServiceError", func(t *testing.T) {
		svc := new(MockUserService)
		svc.
			On("ChangePassword", mock.Anything, uint(5), "old", "new").
			Return(errors.New("db down")).
			Once()

		app := setupUserApp(svc)
		body := map[string]string{"old_password": "old", "new_password": "new"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/5/password", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", validRole)

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		svc.AssertExpectations(t)
	})

	t.Run("Success", func(t *testing.T) {
		svc := new(MockUserService)
		svc.
			On("ChangePassword", mock.Anything, uint(3), "old", "new").
			Return(nil).
			Once()

		app := setupUserApp(svc)
		body := map[string]string{"old_password": "old", "new_password": "new"}
		buf, _ := json.Marshal(body)

		req := httptest.NewRequest(http.MethodPut, "/users/3/password", bytes.NewReader(buf))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Role", validRole)

		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody struct {
			Status  string `json:"status"`
			Message string `json:"message"`
		}
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "success", respBody.Status)
		assert.Equal(t, "Password changed successfully", respBody.Message)

		svc.AssertExpectations(t)
	})
}

func TestAuthController_Me(t *testing.T) {
    tests := []struct {
        name           string
        headerValue    string
        mockBehavior   func(ctx context.Context, userID uint) (*corePort.MeDTO, error)
        expectedStatus int
        expectedBody   fiber.Map
    }{
        {
            name:           "NoUserID",
            headerValue:    "",
            mockBehavior:   nil, // service should not be called
            expectedStatus: http.StatusUnauthorized,
            expectedBody: fiber.Map{
                "status":  "error",
                "message": "User not authenticated",
            },
        },
        {
            name:        "InvalidUserIDType",
            headerValue: "abc", // parse fail → no Locals, same as no user
            mockBehavior: nil,
            expectedStatus: http.StatusUnauthorized,
            expectedBody: fiber.Map{
                "status":  "error",
                "message": "User not authenticated",
            },
        },
        {
            name:        "ServiceError",
            headerValue: "42",
            mockBehavior: func(ctx context.Context, userID uint) (*corePort.MeDTO, error) {
                return nil, errors.New("db failure")
            },
            expectedStatus: http.StatusInternalServerError,
            expectedBody: fiber.Map{
                "status":  "error",
                "message": "Failed to fetch user info",
                "error":   "db failure",
            },
        },
        {
            name:        "UserNotFound",
            headerValue: "100",
            mockBehavior: func(ctx context.Context, userID uint) (*corePort.MeDTO, error) {
                return nil, nil
            },
            expectedStatus: http.StatusNotFound,
            expectedBody: fiber.Map{
                "status":  "error",
                "message": "User not found",
            },
        },
        {
            name:        "Success",
            headerValue: "7",
            mockBehavior: func(ctx context.Context, userID uint) (*corePort.MeDTO, error) {
                return &corePort.MeDTO{
                    ID:        7,
                    Username:  "jane",
                    Email:     "jane@x.com",
                    BranchID:  nil,
                    TenantIDs: []uint{1, 2},
                }, nil
            },
            expectedStatus: http.StatusOK,
            expectedBody: fiber.Map{
                "status":  "success",
                "message": "User profile retrieved",
                "data": map[string]interface{}{
                    "id":         float64(7), // json.Unmarshal turns numbers into float64
                    "username":   "jane",
                    "email":      "jane@x.com",
                    "branch_id":  nil,
                    "tenant_ids": []interface{}{float64(1), float64(2)},
                },
            },
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // เตรียม mock service
            var svc corePort.IUser
            if tt.mockBehavior != nil {
                svc = &MockUserService{MeFunc: tt.mockBehavior}
            } else {
                svc = &MockUserService{} // won't be called
            }

            app := setupUserApp(svc)
            req := httptest.NewRequest("GET", "/auth/me", nil)
            if tt.headerValue != "" {
                req.Header.Set("X-User-ID", tt.headerValue)
            }

            resp, err := app.Test(req, -1)
            require.NoError(t, err)
            defer resp.Body.Close()

            // เช็ค status code
            assert.Equal(t, tt.expectedStatus, resp.StatusCode)

            // อ่าน body
            buf := new(bytes.Buffer)
            _, err = buf.ReadFrom(resp.Body)
            require.NoError(t, err)

            var body fiber.Map
            require.NoError(t, json.Unmarshal(buf.Bytes(), &body))

            // เปรียบเทียบ map
            assert.Equal(t, tt.expectedBody, body)
        })
    }
}
