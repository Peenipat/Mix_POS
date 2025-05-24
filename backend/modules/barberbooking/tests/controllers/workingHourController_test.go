package barberBookingControllers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"strings"
)

type MockWorkingHourService struct {
	mock.Mock
}

func (m *MockWorkingHourService) GetWorkingHours(ctx context.Context, branchID uint) ([]barberBookingModels.WorkingHour, error) {
	args := m.Called(ctx, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]barberBookingModels.WorkingHour), args.Error(1)
}

func (m *MockWorkingHourService) UpdateWorkingHours(ctx context.Context, branchID uint, input []barberBookingDto.WorkingHourInput) error {
	args := m.Called(ctx, branchID, input)
	return args.Error(0)
}
func (m *MockWorkingHourService) CreateWorkingHours(ctx context.Context, branchID uint, input barberBookingDto.WorkingHourInput) error {
	args := m.Called(ctx, branchID, input)
	return args.Error(0)
}

func setupWorkingHourTestApp(mockSvc barberBookingPort.IWorkingHourService) *fiber.App {
	app := fiber.New()
	controller := barberBookingControllers.NewWorkingHourController(mockSvc)

	// Mock Role Middleware
	app.Use(func(c *fiber.Ctx) error {
		role := c.Get("X-Mock-Role")
		if role != "" {
			c.Locals("role", role)
		}
		return c.Next()
	})

	// Register routes AFTER middleware
	app.Get("/branches/:branch_id/working-hours", controller.GetWorkingHours)
	app.Put("/branches/:branch_id/working-hours", controller.UpdateWorkingHours)
	app.Post("/branches/:branch_id/working-hours", controller.CreateWorkingHours)


	return app
}

func TestGetWorkingHours(t *testing.T) {
	mockSvc := new(MockWorkingHourService)
	app := setupWorkingHourTestApp(mockSvc)

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/branches/abc/working-hours", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc.On("GetWorkingHours", mock.Anything, uint(1)).Return([]barberBookingModels.WorkingHour{}, nil).Once()
		req := httptest.NewRequest(http.MethodGet, "/branches/1/working-hours", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
	})

	t.Run("Success_ShouldReturnList", func(t *testing.T) {
		hours := []barberBookingModels.WorkingHour{
			{ID: 1, BranchID: 1, Weekday: 1},
			{ID: 2, BranchID: 1, Weekday: 2},
		}
		mockSvc.On("GetWorkingHours", mock.Anything, uint(1)).Return(hours, nil).Once()

		req := httptest.NewRequest(http.MethodGet, "/branches/1/working-hours", nil)
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status  string                            `json:"status"`
			Message string                            `json:"message"`
			Data    []barberBookingModels.WorkingHour `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, "Working hours retrieved", body.Message)
		assert.Len(t, body.Data, 2)
	})
}

func TestUpdateWorkingHoursAuth(t *testing.T) {
	mockSvc := new(MockWorkingHourService)
	app := setupWorkingHourTestApp(mockSvc)

	t.Run("InvalidID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/branches/abc/working-hours", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
}

func TestUpdateWorkingHours(t *testing.T) {
	mockSvc := new(MockWorkingHourService)
	app := setupWorkingHourTestApp(mockSvc)

	app.Use(func(c *fiber.Ctx) error {
		role := c.Get("X-Mock-Role")
		if role != "" {
			c.Locals("role", role)
		}
		return c.Next()
	})
	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/branches/1/working-hours", strings.NewReader("not-json"))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("EmptyInput_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPut, "/branches/1/working-hours", strings.NewReader("[]"))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	// t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
	// 	payload := `[{"weekday":1,"start_time":"09:00","end_time":"17:00"}]`
	// 	req := httptest.NewRequest(http.MethodPut, "/branches/1/working-hours", strings.NewReader(payload))
	// 	req.Header.Set("X-Mock-Role", "USER") // ไม่มีสิทธิ์
	// 	req.Header.Set("Content-Type", "application/json")
	// 	resp, err := app.Test(req)

	// 	assert.NoError(t, err)
	// 	assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	// })

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		layout := time.RFC3339

		// เวลาแบบ full-format
		startTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)

		// JSON ต้องใช้ format ให้ match กับ time.Time
		payload := `[{
			"weekday": 1,
			"start_time": "` + startTime.Format(layout) + `",
			"end_time": "` + endTime.Format(layout) + `"
		}]`

		inputs := []barberBookingDto.WorkingHourInput{{
			Weekday:   1,
			StartTime: startTime,
			EndTime:   endTime,
		}}

		mockSvc.On("UpdateWorkingHours", mock.Anything, uint(1), inputs).Return(nil).Once()

		req := httptest.NewRequest(http.MethodPut, "/branches/1/working-hours", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")
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
		assert.Equal(t, "Working hours updated", body.Message)
	})

}

func TestCreateWorkingHours(t *testing.T) {
	mockSvc := new(MockWorkingHourService)
	app := setupWorkingHourTestApp(mockSvc)

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		payload := `{"weekday":1,"start_time":"0000-01-01T09:00:00Z","end_time":"0000-01-01T17:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/branches/1/working-hours", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "USER") // ไม่มีสิทธิ์
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("InvalidBranchID_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/branches/abc/working-hours", nil)
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBody_ShouldReturn400", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/branches/1/working-hours", strings.NewReader("not-json"))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingTime_ShouldReturn400", func(t *testing.T) {
		payload := `{"weekday":1}`
		req := httptest.NewRequest(http.MethodPost, "/branches/1/working-hours", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidWeekday_ShouldReturn400", func(t *testing.T) {
		payload := `{"weekday":8,"start_time":"0000-01-01T09:00:00Z","end_time":"0000-01-01T17:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/branches/1/working-hours", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("Success_ShouldReturn201", func(t *testing.T) {
		startTime := time.Date(0, 1, 1, 9, 0, 0, 0, time.UTC)
		endTime := time.Date(0, 1, 1, 17, 0, 0, 0, time.UTC)

		input := barberBookingDto.WorkingHourInput{
			Weekday:   1,
			StartTime: startTime,
			EndTime:   endTime,
		}

		mockSvc.On("CreateWorkingHours", mock.Anything, uint(1), input).Return(nil).Once()

		payload := `{"weekday":1,"start_time":"0000-01-01T09:00:00Z","end_time":"0000-01-01T17:00:00Z"}`
		req := httptest.NewRequest(http.MethodPost, "/branches/1/working-hours", strings.NewReader(payload))
		req.Header.Set("X-Mock-Role", "SAAS_SUPER_ADMIN")
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
		assert.Equal(t, "Working hour created", body.Message)
	})
}

