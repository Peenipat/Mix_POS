package barberBookingControllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	// barberBookingPort "myapp/modules/barberbooking/port"
	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
)

type MockAppointmentService struct {
	mock.Mock
}

// GetAppointmentsByBarber implements barberBookingPort.IAppointment.
func (m *MockAppointmentService) GetAppointmentsByBarber(ctx context.Context, barberID uint, filter barberBookingPort.AppointmentFilter) ([]barberBookingPort.AppointmentBrief, error) {
	panic("unimplemented")
}

// GetAppointmentsByBranch implements barberBookingPort.IAppointment.
func (m *MockAppointmentService) GetAppointmentsByBranch(ctx context.Context, branchID uint, start *time.Time, end *time.Time) ([]barberBookingPort.AppointmentBrief, error) {
	panic("unimplemented")
}

// ListAppointmentsResponse implements barberBookingPort.IAppointment.
func (m *MockAppointmentService) ListAppointmentsResponse(ctx context.Context, filter barberBookingDto.AppointmentFilter) ([]barberBookingPort.AppointmentResponse, error) {
	panic("unimplemented")
}

func (m *MockAppointmentService) CheckBarberAvailability(
	ctx context.Context,
	tenantID, barberID uint,
	start, end time.Time,
) (bool, error) {
	args := m.Called(ctx, tenantID, barberID, start, end)
	return args.Bool(0), args.Error(1)
}

func (m *MockAppointmentService) CreateAppointment(
	ctx context.Context,
	input *barberBookingModels.Appointment,
) (*barberBookingDto.AppointmentResponseDTO, error) {
	args := m.Called(ctx, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingDto.AppointmentResponseDTO), args.Error(1)
}

func (m *MockAppointmentService) GetAvailableBarbers(
	ctx context.Context,
	tenantID, branchID uint,
	start, end time.Time,
) ([]barberBookingModels.Barber, error) {
	args := m.Called(ctx, tenantID, branchID, start, end)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]barberBookingModels.Barber), args.Error(1)
}

func (m *MockAppointmentService) UpdateAppointment(
	ctx context.Context,
	id uint,
	tenantID uint,
	input *barberBookingModels.Appointment,
) (*barberBookingModels.Appointment, error) {
	args := m.Called(ctx, id, tenantID, input)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Appointment), args.Error(1)
}

func (m *MockAppointmentService) GetAppointmentByID(
	ctx context.Context,
	id uint,
) (*barberBookingModels.Appointment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.Appointment), args.Error(1)
}

func (m *MockAppointmentService) ListAppointments(
	ctx context.Context,
	filter barberBookingDto.AppointmentFilter,
) ([]barberBookingModels.Appointment, error) {
	args := m.Called(ctx, filter)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]barberBookingModels.Appointment), args.Error(1)
}

func (m *MockAppointmentService) CancelAppointment(
	ctx context.Context,
	appointmentID uint,
	actorUserID *uint,
	actorCustomerID *uint,
) error {
	args := m.Called(ctx, appointmentID, actorUserID, actorCustomerID)
	return args.Error(0)
}

func (m *MockAppointmentService) RescheduleAppointment(
	ctx context.Context,
	appointmentID uint,
	newStartTime time.Time,
	actorUserID *uint,
	actorCustomerID *uint,
) error {
	args := m.Called(ctx, appointmentID, newStartTime, actorUserID, actorCustomerID)
	return args.Error(0)
}

func (m *MockAppointmentService) CalculateAppointmentEndTime(
	ctx context.Context,
	serviceID uint,
	startTime time.Time,
) (time.Time, error) {
	args := m.Called(ctx, serviceID, startTime)
	// TODO: Return real time in tests
	return args.Get(0).(time.Time), args.Error(1)
}

func (m *MockAppointmentService) DeleteAppointment(
	ctx context.Context,
	appointmentID uint,
) error {
	args := m.Called(ctx, appointmentID)
	return args.Error(0)
}

func (m *MockAppointmentService) GetUpcomingAppointmentsByCustomer(
	ctx context.Context,
	customerID uint,
) (*barberBookingModels.Appointment, error) {
	args := m.Called(ctx, customerID)
	// TODO: Return real *Appointment in tests
	return args.Get(0).(*barberBookingModels.Appointment), args.Error(1)
}

func setupAppointmentApp(mockSvc *MockAppointmentService) *fiber.App {
	app := fiber.New()
	ctrl := barberBookingControllers.NewAppointmentController(mockSvc)
	app.Get("/tenants/:tenant_id/barbers/:barber_id/availability", ctrl.CheckBarberAvailability)
	app.Post("/tenants/:tenant_id/appointments", ctrl.CreateAppointment)
	app.Get("/tenants/:tenant_id/branches/:branch_id/available-barbers", ctrl.GetAvailableBarbers)
	app.Put("/tenants/:tenant_id/appointments/:appointment_id", ctrl.UpdateAppointment)
	app.Get("/tenants/:tenant_id/appointments/:appointment_id", ctrl.GetAppointmentByID)
	app.Get("/tenants/:tenant_id/appointments", ctrl.ListAppointments)
	app.Post("/tenants/:tenant_id/appointments/:appointment_id/cancel", ctrl.CancelAppointment)
	app.Post("/tenants/:tenant_id/appointments/:appointment_id/reschedule", ctrl.RescheduleAppointment)
	return app
}

func setupAppWithRole(svc *MockAppointmentService, role coreModels.RoleName) *fiber.App {
	app := fiber.New()
	// inject role into locals
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", string(role))
		return c.Next()
	})
	ctrl := barberBookingControllers.NewAppointmentController(svc)
	app.Delete("/tenants/:tenant_id/appointments/:appointment_id", ctrl.DeleteAppointment)
	return app
}

func TestCheckBarberAvailability_Controller(t *testing.T) {
	validStart := time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC).Format(time.RFC3339)
	validEnd := time.Date(2025, 5, 14, 10, 0, 0, 0, time.UTC).Format(time.RFC3339)

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/foo/barbers/123/availability?start="+validStart+"&end="+validEnd, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Invalid tenant_id", body["message"])

		mockSvc.AssertNotCalled(t, "CheckBarberAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("InvalidBarberID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/xyz/availability?start="+validStart+"&end="+validEnd, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Invalid barber_id", body["message"])

		mockSvc.AssertNotCalled(t, "CheckBarberAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("MissingStart_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/123/availability?end="+validEnd, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Missing start time", body["message"])

		mockSvc.AssertNotCalled(t, "CheckBarberAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("InvalidStartFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/123/availability?start=14-05-2025&end="+validEnd, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Invalid start time format. Expect RFC3339", body["message"])

		mockSvc.AssertNotCalled(t, "CheckBarberAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("MissingEnd_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/123/availability?start="+validStart, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Missing end time", body["message"])

		mockSvc.AssertNotCalled(t, "CheckBarberAvailability", mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("InvalidStartFormat_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		// ใช้ตัวอย่างง่ายๆ ที่ไม่มี space แต่ format invalid
		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/barbers/1/appointments?start=invalid-time-format",
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["message"], "Invalid start format")
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		mockSvc.
			On("CheckBarberAvailability",
				mock.Anything,
				uint(1),
				uint(123),
				mock.AnythingOfType("time.Time"),
				mock.AnythingOfType("time.Time"),
			).
			Return(false, assert.AnError).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/123/availability?start="+validStart+"&end="+validEnd, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnAvailability", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		mockSvc.
			On("CheckBarberAvailability",
				mock.Anything,
				uint(1),
				uint(123),
				mock.AnythingOfType("time.Time"),
				mock.AnythingOfType("time.Time"),
			).
			Return(true, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/123/availability?start="+validStart+"&end="+validEnd, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status    string `json:"status"`
			Available bool   `json:"available"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "success", body.Status)
		assert.True(t, body.Available)

		mockSvc.AssertExpectations(t)
	})
}

func TestCreateAppointment_Controller(t *testing.T) {
	validStart := time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC).Format(time.RFC3339)

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		body := `{"branch_id":1,"service_id":1,"customer_id":1,"start_time":"` + validStart + `"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/abc/appointments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidJSON_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments", strings.NewReader("not-json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingRequiredFields_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		// branch_id missing
		body := `{"service_id":1,"customer_id":1,"start_time":"` + validStart + `"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var j map[string]string
		json.NewDecoder(resp.Body).Decode(&j)
		assert.Contains(t, j["message"], "Missing required fields")
	})

	t.Run("InvalidStartTimeFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		body := `{"branch_id":1,"service_id":1,"customer_id":1,"start_time":"14-05-2025"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var j map[string]string
		json.NewDecoder(resp.Body).Decode(&j)
		assert.Contains(t, j["message"], "Invalid start_time format")
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)
		// stub CreateAppointment to error
		mockSvc.
			On("CreateAppointment",
				mock.Anything,
				mock.MatchedBy(func(a *barberBookingModels.Appointment) bool {
					return a.TenantID == 1 && a.BranchID == 1
				}),
			).
			Return(nil, assert.AnError).
			Once()

		body := `{"branch_id":1,"service_id":1,"customer_id":1,"start_time":"` + validStart + `"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnCreated", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		created := &barberBookingModels.Appointment{
			ID:         10,
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC),
		}
		mockSvc.
			On("CreateAppointment",
				mock.Anything,
				mock.MatchedBy(func(a *barberBookingModels.Appointment) bool {
					return a.TenantID == 1 && a.BranchID == 1
				}),
			).
			Return(created, nil).
			Once()

		body := `{"branch_id":1,"service_id":1,"customer_id":1,"start_time":"` + validStart + `"}`
		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusCreated, resp.StatusCode)

		var respBody struct {
			Status string                          `json:"status"`
			Data   barberBookingModels.Appointment `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&respBody)

		assert.Equal(t, "success", respBody.Status)
		assert.Equal(t, uint(10), respBody.Data.ID)
		mockSvc.AssertExpectations(t)
	})
}

func TestGetAvailableBarbers_Controller(t *testing.T) {
	start := time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC).Format(time.RFC3339)
	end := time.Date(2025, 5, 14, 10, 0, 0, 0, time.UTC).Format(time.RFC3339)

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/foo/branches/1/available-barbers?start="+start+"&end="+end,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBranchID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/abc/available-barbers?start="+start+"&end="+end,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingStart_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/1/available-barbers?end="+end,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidStartFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/1/available-barbers?start=14-05-2025&end="+end,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingEnd_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/1/available-barbers?start="+start,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidEndFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/1/available-barbers?start="+start+"&end=10:00",
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		mockSvc.
			On("GetAvailableBarbers",
				mock.Anything,
				uint(1), uint(1),
				mock.AnythingOfType("time.Time"),
				mock.AnythingOfType("time.Time"),
			).
			Return(nil, assert.AnError).
			Once()

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/1/available-barbers?start="+start+"&end="+end,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnBarbers", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		expected := []barberBookingModels.Barber{
			{ID: 10, BranchID: 1, TenantID: 1},
			{ID: 11, BranchID: 1, TenantID: 1},
		}
		mockSvc.
			On("GetAvailableBarbers",
				mock.Anything,
				uint(1), uint(1),
				mock.AnythingOfType("time.Time"),
				mock.AnythingOfType("time.Time"),
			).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/branches/1/available-barbers?start="+start+"&end="+end,
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                       `json:"status"`
			Data   []barberBookingModels.Barber `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected, body.Data)
		mockSvc.AssertExpectations(t)
	})
}

func TestUpdateAppointment_Controller(t *testing.T) {
	validBody := `{"start_time":"2025-05-14T09:00:00Z","service_id":2,"barber_id":3,"status":"CONFIRMED"}`

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/tenants/foo/appointments/10", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidAppointmentID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/abc", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidJSON_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/10", strings.NewReader("not-json"))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		mockSvc.
			On("UpdateAppointment",
				mock.Anything,
				uint(10),
				uint(1),
				mock.Anything,
			).
			Return(nil, assert.AnError).
			Once()

		req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/10", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnUpdated", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		updated := &barberBookingModels.Appointment{
			ID:         10,
			TenantID:   1,
			BranchID:   5,
			ServiceID:  2,
			BarberID:   3,
			CustomerID: 7,
			StartTime:  time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC),
			EndTime:    time.Date(2025, 5, 14, 9, 30, 0, 0, time.UTC),
			Status:     barberBookingModels.StatusConfirmed,
		}
		mockSvc.
			On("UpdateAppointment",
				mock.Anything,
				uint(10),
				uint(1),
				mock.Anything,
			).
			Return(updated, nil).
			Once()

		req := httptest.NewRequest(http.MethodPut, "/tenants/1/appointments/10", strings.NewReader(validBody))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                          `json:"status"`
			Data   barberBookingModels.Appointment `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, updated.ID, body.Data.ID)
		assert.Equal(t, updated.ServiceID, body.Data.ServiceID)
		assert.Equal(t, updated.Status, body.Data.Status)

		mockSvc.AssertExpectations(t)
	})
}

func TestGetAppointmentByID_Controller(t *testing.T) {
	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/foo/appointments/10", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidAppointmentID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		mockSvc.
			On("GetAppointmentByID", mock.Anything, uint(10)).
			Return(
				(*barberBookingModels.Appointment)(nil),
				fmt.Errorf("appointment with ID %d not found", 10),
			).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/10", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["message"], "not found")

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		// Return a wrapped error not containing "not found"
		mockSvc.
			On("GetAppointmentByID", mock.Anything, uint(20)).
			Return((*barberBookingModels.Appointment)(nil), fmt.Errorf("db connection lost")).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/20", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Failed to fetch appointment", body["message"])

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		expected := &barberBookingModels.Appointment{
			ID:         10,
			TenantID:   1,
			BranchID:   5,
			ServiceID:  2,
			CustomerID: 3,
			StartTime:  time.Date(2025, 5, 14, 9, 0, 0, 0, time.UTC),
			EndTime:    time.Date(2025, 5, 14, 9, 30, 0, 0, time.UTC),
			Status:     barberBookingModels.StatusPending,
		}
		mockSvc.
			On("GetAppointmentByID", mock.Anything, uint(10)).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/appointments/10", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                          `json:"status"`
			Data   barberBookingModels.Appointment `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected.ID, body.Data.ID)

		mockSvc.AssertExpectations(t)
	})
}

func TestListAppointments_Controller(t *testing.T) {
	baseURL := "/tenants/1/appointments"

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/foo/appointments", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBranchID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?branch_id=abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidBarberID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?barber_id=xyz", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidCustomerID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?customer_id=!!!", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidStartDateFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?start_date=14-05-2025", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidEndDateFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?end_date=tomorrow", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidLimit_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?limit=abc", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidOffset_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet, baseURL+"?offset=-5", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		// stub service error for empty filter
		mockSvc.
			On("ListAppointments", mock.Anything, barberBookingDto.AppointmentFilter{TenantID: 1}).
			Return(nil, assert.AnError).
			Once()

		req := httptest.NewRequest(http.MethodGet, baseURL, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturnAppointments", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		// prepare filter and expected result
		filter := barberBookingDto.AppointmentFilter{TenantID: 1}
		expected := []barberBookingModels.Appointment{
			{ID: 1, TenantID: 1, BranchID: 2, ServiceID: 3},
			{ID: 2, TenantID: 1, BranchID: 2, ServiceID: 4},
		}

		mockSvc.
			On("ListAppointments", mock.Anything, filter).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, baseURL, nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                            `json:"status"`
			Data   []barberBookingModels.Appointment `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected, body.Data)

		mockSvc.AssertExpectations(t)
	})
}

func TestCancelAppointment_Controller(t *testing.T) {
	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/foo/appointments/1/cancel", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidAppointmentID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/bar/cancel", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidJSONBody_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/1/cancel",
			bytes.NewBufferString("{invalid-json}"),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingActor_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		body := `{}` // no actor_user_id or actor_customer_id
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/1/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Contains(t, respBody["message"], "Either actor_user_id or actor_customer_id")
	})

	t.Run("BothActorsProvided_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		body := `{"actor_user_id":1,"actor_customer_id":2}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/1/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Contains(t, respBody["message"], "Either actor_user_id or actor_customer_id")
	})

	t.Run("ServiceError_NotFound_ShouldReturn404", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		userID := uint(1)
		mockSvc.
			On("CancelAppointment", mock.Anything, uint(42), &userID, (*uint)(nil)).
			Return(fmt.Errorf("appointment with ID %d not found", 42)).
			Once()

		body := `{"actor_user_id":1}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/42/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Contains(t, respBody["message"], "not found")

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError_CannotCancel_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		userID := uint(1)
		mockSvc.
			On("CancelAppointment", mock.Anything, uint(100), &userID, (*uint)(nil)).
			Return(errors.New("appointment cannot be cancelled in its current status")).
			Once()

		body := `{"actor_user_id":1}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/100/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Contains(t, respBody["message"], "cannot be cancelled")

		mockSvc.AssertExpectations(t)
	})

	t.Run("ServiceError_Other_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		userID := uint(1)
		mockSvc.
			On("CancelAppointment", mock.Anything, uint(200), &userID, (*uint)(nil)).
			Return(errors.New("database failure")).
			Once()

		body := `{"actor_user_id":1}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/200/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "Failed to cancel appointment", respBody["message"])

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_UserCancels_ShouldReturn200", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		userID := uint(123)
		mockSvc.
			On("CancelAppointment", mock.Anything, uint(300), &userID, (*uint)(nil)).
			Return(nil).
			Once()

		body := `{"actor_user_id":123}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/300/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "success", respBody["status"])
		assert.Equal(t, "appointment cancelled", respBody["message"])

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_CustomerCancels_ShouldReturn200", func(t *testing.T) {
		mockSvc := new(MockAppointmentService)
		app := setupAppointmentApp(mockSvc)

		custID := uint(77)
		mockSvc.
			On("CancelAppointment", mock.Anything, uint(400), (*uint)(nil), &custID).
			Return(nil).
			Once()

		body := `{"actor_customer_id":77}`
		req := httptest.NewRequest(
			http.MethodPost,
			"/tenants/1/appointments/400/cancel",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var respBody map[string]string
		json.NewDecoder(resp.Body).Decode(&respBody)
		assert.Equal(t, "success", respBody["status"])
		assert.Equal(t, "appointment cancelled", respBody["message"])

		mockSvc.AssertExpectations(t)
	})
}

func TestRescheduleAppointment_Controller(t *testing.T) {
	const dateFmt = time.RFC3339
	validNew := "2025-05-21T10:00:00Z"
	newTime, _ := time.Parse(dateFmt, validNew)

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/abc/appointments/1/reschedule", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidAppointmentID_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		req := httptest.NewRequest(http.MethodPost, "/tenants/1/appointments/xyz/reschedule", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidJSONBody_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/1/reschedule",
			bytes.NewBufferString("{invalid-json}"),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("MissingNewStartTime_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		body := `{}` // no new_start_time
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/1/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Equal(t, "Missing new_start_time", b["message"])
	})

	t.Run("InvalidNewStartTimeFormat_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		body := `{"new_start_time":"2025-05-21 10:00:00"}`
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/1/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Contains(t, b["message"], "Invalid new_start_time format")
	})

	t.Run("MissingActor_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		body := fmt.Sprintf(`{"new_start_time":"%s"}`, validNew)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/1/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Contains(t, b["message"], "Either actor_user_id or actor_customer_id")
	})

	t.Run("BothActorsProvided_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		body := fmt.Sprintf(
			`{"new_start_time":"%s","actor_user_id":1,"actor_customer_id":2}`,
			validNew,
		)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/1/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Contains(t, b["message"], "Either actor_user_id or actor_customer_id")
	})

	t.Run("ServiceError_NotFound_ShouldReturn404", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		userID := uint(5)
		svc.
			On("RescheduleAppointment",
				mock.Anything,
				uint(42),
				newTime,
				&userID,
				(*uint)(nil),
			).
			Return(fmt.Errorf("appointment with ID %d not found", 42)).
			Once()

		body := fmt.Sprintf(`{"new_start_time":"%s","actor_user_id":5}`, validNew)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/42/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Contains(t, b["message"], "not found")

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_CannotReschedule_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		custID := uint(7)
		svc.
			On("RescheduleAppointment",
				mock.Anything,
				uint(100),
				newTime,
				(*uint)(nil),
				&custID,
			).
			Return(errors.New("cannot reschedule: time slot conflicts")).
			Once()

		body := fmt.Sprintf(`{"new_start_time":"%s","actor_customer_id":7}`, validNew)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/100/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Contains(t, b["message"], "cannot reschedule")

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Other_ShouldReturn500", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		userID := uint(9)
		svc.
			On("RescheduleAppointment",
				mock.Anything,
				uint(200),
				newTime,
				&userID,
				(*uint)(nil),
			).
			Return(errors.New("database failure")).
			Once()

		body := fmt.Sprintf(`{"new_start_time":"%s","actor_user_id":9}`, validNew)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/200/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Equal(t, "Failed to reschedule appointment", b["message"])

		svc.AssertExpectations(t)
	})

	t.Run("Success_UserReschedules_ShouldReturn200", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		userID := uint(123)
		svc.
			On("RescheduleAppointment",
				mock.Anything,
				uint(300),
				newTime,
				&userID,
				(*uint)(nil),
			).
			Return(nil).
			Once()

		body := fmt.Sprintf(`{"new_start_time":"%s","actor_user_id":123}`, validNew)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/300/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Equal(t, "success", b["status"])
		assert.Equal(t, "appointment rescheduled", b["message"])

		svc.AssertExpectations(t)
	})

	t.Run("Success_CustomerReschedules_ShouldReturn200", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		custID := uint(456)
		svc.
			On("RescheduleAppointment",
				mock.Anything,
				uint(400),
				newTime,
				(*uint)(nil),
				&custID,
			).
			Return(nil).
			Once()

		body := fmt.Sprintf(`{"new_start_time":"%s","actor_customer_id":456}`, validNew)
		req := httptest.NewRequest(http.MethodPost,
			"/tenants/1/appointments/400/reschedule",
			bytes.NewBufferString(body),
		)
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var b map[string]string
		json.NewDecoder(resp.Body).Decode(&b)
		assert.Equal(t, "success", b["status"])
		assert.Equal(t, "appointment rescheduled", b["message"])

		svc.AssertExpectations(t)
	})
}

func TestGetAppointmentsByBarber_Controller(t *testing.T) {
	const dateFmt = time.RFC3339
	now := time.Now().UTC().Truncate(time.Second)
	later := now.Add(1 * time.Hour)

	t.Run("InvalidStartFormat_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		// ใช้ตัวอย่างที่ไม่มี space แต่ format invalid
		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/barbers/1/appointments?start=2025-05-21%2010:00:00",
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["message"], "Invalid start format")
	})

	t.Run("InvalidBarberID_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		req := httptest.NewRequest(http.MethodGet, "/tenants/1/barbers/xyz/appointments", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_BarberNotFound_ShouldReturn404", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		svc.
			On("GetAppointmentsByBarber", mock.Anything, uint(42), (*time.Time)(nil), (*time.Time)(nil)).
			Return(nil, fmt.Errorf("barber with ID %d not found", 42)).
			Once()

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/barbers/42/appointments",
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["message"], "not found")

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Other_ShouldReturn500", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		svc.
			On("GetAppointmentsByBarber", mock.Anything, uint(1), (*time.Time)(nil), (*time.Time)(nil)).
			Return(nil, errors.New("db error")).
			Once()

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/barbers/1/appointments",
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Failed to fetch appointments", body["message"])

		svc.AssertExpectations(t)
	})

	t.Run("Success_NoFilters_ShouldReturnAppointments", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		sample := []barberBookingModels.Appointment{
			{ID: 101, StartTime: now, EndTime: later},
		}
		svc.
			On("GetAppointmentsByBarber", mock.Anything, uint(1), (*time.Time)(nil), (*time.Time)(nil)).
			Return(sample, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet,
			"/tenants/1/barbers/1/appointments",
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                            `json:"status"`
			Data   []barberBookingModels.Appointment `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, sample, body.Data)

		svc.AssertExpectations(t)
	})

	t.Run("Success_WithFilters_ShouldReturnAppointments", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppointmentApp(svc)

		sample := []barberBookingModels.Appointment{
			{ID: 202, StartTime: now.Add(10 * time.Minute), EndTime: later},
		}
		startPtr := &now
		endPtr := &later
		svc.
			On("GetAppointmentsByBarber", mock.Anything, uint(1), startPtr, endPtr).
			Return(sample, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet,
			fmt.Sprintf("/tenants/1/barbers/1/appointments?start=%s&end=%s", now.Format(dateFmt), later.Format(dateFmt)),
			nil,
		)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                            `json:"status"`
			Data   []barberBookingModels.Appointment `json:"data"`
		}
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, sample, body.Data)

		svc.AssertExpectations(t)
	})
}

func TestDeleteAppointment_Controller(t *testing.T) {
	validRole := coreModels.RoleNameTenantAdmin
	invalidRole := coreModels.RoleNameBranchAdmin // assume this is not in RolesCanManageAppointment

	t.Run("Forbidden_InvalidRole_ShouldReturn403", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppWithRole(svc, invalidRole)

		req := httptest.NewRequest(http.MethodDelete,
			"/tenants/1/appointments/10", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Permission denied", body["message"])
	})

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppWithRole(svc, validRole)

		req := httptest.NewRequest(http.MethodDelete,
			"/tenants/abc/appointments/10", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Invalid tenant_id", body["message"])
	})

	t.Run("InvalidAppointmentID_ShouldReturn400", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppWithRole(svc, validRole)

		req := httptest.NewRequest(http.MethodDelete,
			"/tenants/1/appointments/xyz", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Invalid appointment_id", body["message"])
	})

	t.Run("ServiceError_NotFound_ShouldReturn404", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppWithRole(svc, validRole)

		svc.
			On("DeleteAppointment", mock.Anything, uint(42)).
			Return(fmt.Errorf("appointment with ID %d not found", 42)).
			Once()

		req := httptest.NewRequest(http.MethodDelete,
			"/tenants/1/appointments/42", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusNotFound, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["message"], "not found")

		svc.AssertExpectations(t)
	})

	t.Run("ServiceError_Other_ShouldReturn500", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppWithRole(svc, validRole)

		svc.
			On("DeleteAppointment", mock.Anything, uint(100)).
			Return(errors.New("db failure")).
			Once()

		req := httptest.NewRequest(http.MethodDelete,
			"/tenants/1/appointments/100", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Contains(t, body["message"], "Failed to delete appointment")

		svc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		svc := new(MockAppointmentService)
		app := setupAppWithRole(svc, validRole)

		svc.
			On("DeleteAppointment", mock.Anything, uint(200)).
			Return(nil).
			Once()

		req := httptest.NewRequest(http.MethodDelete,
			"/tenants/1/appointments/200", nil)
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "appointment deleted", body["message"])

		svc.AssertExpectations(t)
	})
}
