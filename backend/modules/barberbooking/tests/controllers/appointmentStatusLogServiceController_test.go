package barberBookingControllers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	controller "myapp/modules/barberbooking/controllers"
	models "myapp/modules/barberbooking/models"
)

// MockLogService now implements both methods of IAppointmentStatusLogService
type MockLogService struct {
	mock.Mock
}

// LogStatusChange implements barberBookingPort.IAppointmentStatusLogService.
func (m *MockLogService) LogStatusChange(ctx context.Context, appointmentID uint, oldStatus string, newStatus string, userID *uint, customerID *uint, notes string) error {
	panic("unimplemented")
}

func (m *MockLogService) GetLogsForAppointment(ctx context.Context, appointmentID uint) ([]models.AppointmentStatusLog, error) {
    args := m.Called(ctx, appointmentID)
    var logs []models.AppointmentStatusLog
    if v := args.Get(0); v != nil {
        logs = v.([]models.AppointmentStatusLog)
    }
    return logs, args.Error(1)
}

func (m *MockLogService) DeleteLogsByAppointmentID(ctx context.Context, appointmentID uint) error {
	args := m.Called(ctx, appointmentID)
	return args.Error(0)
}

func setupApp(mockSvc *MockLogService) *fiber.App {
	app := fiber.New()
	ctrl := controller.NewAppointmentStatusLogController(mockSvc)
	api := app.Group("/tenants/:tenant_id")
	api.Get("/appointments/:appointment_id/logs", ctrl.GetAppointmentLogs)
	return app
}

func TestGetAppointmentLogs(t *testing.T) {
	t.Run("InvalidTenantID", func(t *testing.T) {
		mockSvc := new(MockLogService)
		app := setupApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/foo/appointments/1/logs", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidAppointmentID", func(t *testing.T) {
		mockSvc := new(MockLogService)
		app := setupApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/bar/logs", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError", func(t *testing.T) {
		mockSvc := new(MockLogService)
		app := setupApp(mockSvc)

		mockSvc.
			On("GetLogsForAppointment", mock.Anything, uint(42)).
			Return(nil, errors.New("db down"))

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/42/logs", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("NoLogs", func(t *testing.T) {
		mockSvc := new(MockLogService)
		app := setupApp(mockSvc)

		mockSvc.
			On("GetLogsForAppointment", mock.Anything, uint(123)).
			Return([]models.AppointmentStatusLog{}, nil)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/123/logs", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                        `json:"status"`
			Data   []models.AppointmentStatusLog `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Empty(t, body.Data)
		mockSvc.AssertExpectations(t)
	})

	t.Run("WithLogs", func(t *testing.T) {
		mockSvc := new(MockLogService)
		app := setupApp(mockSvc)

		now := time.Now().UTC()
		exampleLogs := []models.AppointmentStatusLog{
			{ID: 1, AppointmentID: 123, OldStatus: "", NewStatus: "PENDING", Notes: "init", ChangedAt: now},
			{ID: 2, AppointmentID: 123, OldStatus: "PENDING", NewStatus: "COMPLETED", Notes: "done", ChangedAt: now},
		}
		mockSvc.
			On("GetLogsForAppointment", mock.Anything, uint(123)).
			Return(exampleLogs, nil)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/123/logs", nil)
		resp, _ := app.Test(req, -1)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                        `json:"status"`
			Data   []models.AppointmentStatusLog `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, exampleLogs, body.Data)
		mockSvc.AssertExpectations(t)
	})
}
