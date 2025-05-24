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

    coreControllers "myapp/modules/core/controllers"
    corePort "myapp/modules/core/port"
    coreServices "myapp/modules/core/services"
)

// Mock implementation of IUserService
type MockUserService struct {
    mock.Mock
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

// set up Fiber app with the ChangePassword route
func setupUserApp(svc *MockUserService) *fiber.App {
    app := fiber.New()
    ctrl := coreControllers.NewUserController(svc)

    // inject role into locals
    app.Use(func(c *fiber.Ctx) error {
        c.Locals("role", c.Get("X-Role"))
        return c.Next()
    })

    app.Put("/users/:id/password", ctrl.ChangePassword)
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






