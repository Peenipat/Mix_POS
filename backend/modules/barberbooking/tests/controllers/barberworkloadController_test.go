package barberBookingControllers_test

import (
	"context"
	"encoding/json"
	"github.com/stretchr/testify/require"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockWorkloadService struct {
	mock.Mock
}

func (m *MockWorkloadService) GetWorkloadByBarber(ctx context.Context, barberID uint, date time.Time) (*barberBookingModels.BarberWorkload, error) {
	args := m.Called(ctx, barberID, date)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*barberBookingModels.BarberWorkload), args.Error(1)
}

// MockWorkloadService already has other stubs …
func (m *MockWorkloadService) GetWorkloadSummaryByBranch(
	ctx context.Context,
	date time.Time,
	tenantID uint,
	branchID uint,
) ([]barberBookingDto.BranchWorkloadSummary, error) {
	args := m.Called(ctx, date, tenantID, branchID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]barberBookingDto.BranchWorkloadSummary), args.Error(1)
}

func (m *MockWorkloadService) UpsertBarberWorkload(
	ctx context.Context,
	barberId uint,
	date time.Time,
	appointments int,
	hours int,
) error {
	args := m.Called(ctx, barberId, date, appointments, hours)
	return args.Error(0)
}

func SetupWorkloadTestApp(mockSvc *MockWorkloadService) *fiber.App {
	app := fiber.New()
	app.Use(func(c *fiber.Ctx) error {
		role := c.Get("X-Mock-Role")
		if role != "" {
			c.Locals("role", role)
		}
		return c.Next()
	})
	ctrl := barberBookingController.NewBarberWorkloadController(mockSvc)

	app.Get("/barbers/:barber_id/workload", ctrl.GetWorkloadByBarber)
	app.Get("/workloads/summary", ctrl.GetWorkloadSummaryByBranch)
	app.Put("/barbers/:barber_id/workload", ctrl.UpsertBarberWorkload)
	return app
}

func TestGetWorkloadByBarber(t *testing.T) {

	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/barbers/1/workload", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("PermissionDenied_ShouldReturn403# Not branch_admin access", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/barbers/1/workload", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[2]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})
	t.Run("BarberID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/barbers/abc/workload", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[0]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})
	t.Run("BadDateFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(
			http.MethodGet,
			"/barbers/123/workload?date=2025-13-01",
			nil,
		)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		mockSvc.On("GetWorkloadByBarber", mock.Anything, uint(123), mock.AnythingOfType("time.Time")).
			Return(nil, assert.AnError).Once()

		req := httptest.NewRequest(
			http.MethodGet,
			"/barbers/123/workload?date=2025-05-14",
			nil,
		)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)

		assert.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("NoRecord_ShouldReturnZeroValue", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		mockSvc.
			On("GetWorkloadByBarber",
				mock.Anything, // context ไหนก็ได้
				uint(123),
				mock.AnythingOfType("time.Time"), // ยอมให้เวลาเป็นตัวไหนก็ได้
			).
			Return(nil, nil).
			Once()

		req := httptest.NewRequest("GET", "/barbers/123/workload?date=2025-05-14", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, _ := app.Test(req)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                             `json:"status"`
			Data   barberBookingModels.BarberWorkload `json:"data"`
		}
		err := json.NewDecoder(resp.Body).Decode(&body)
		assert.NoError(t, err)

		// assertions on the zero‐value response
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, uint(123), body.Data.BarberID)
		// ไม่เปรียบเทียบ location ตรงนี้ เพราะเราใช้ mock.AnythingOfType
		assert.Equal(t, 0, body.Data.TotalAppointments)
		assert.Equal(t, 0, body.Data.TotalHours)

		mockSvc.AssertExpectations(t)

		mockSvc.AssertExpectations(t)
	})

	t.Run("FoundRecord_ShouldReturnWorkload", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		// เตรียมวันที่ให้ตรงกับ query
		testDate := time.Date(2025, time.May, 14, 0, 0, 0, 0, time.UTC)

		// สร้าง expected workload object
		expected := &barberBookingModels.BarberWorkload{
			ID:                42,
			BarberID:          123,
			Date:              testDate,
			TotalAppointments: 5,
			TotalHours:        8,
		}

		// ตั้ง mock ให้คืน expected object และ no error
		mockSvc.
			On("GetWorkloadByBarber",
				mock.Anything,
				uint(123),
				mock.MatchedBy(func(dt time.Time) bool {
					y, m, d := dt.Date()
					return y == 2025 && m == time.May && d == 14
				}),
			).
			Return(expected, nil).
			Once()

		// สร้าง request
		req := httptest.NewRequest(
			http.MethodGet,
			"/barbers/123/workload?date=2025-05-14",
			nil,
		)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))

		// ดึง response
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		// Decode body
		var body struct {
			Status string                             `json:"status"`
			Data   barberBookingModels.BarberWorkload `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		// Assertions
		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected.ID, body.Data.ID)
		assert.Equal(t, expected.BarberID, body.Data.BarberID)
		assert.Equal(t, expected.Date, body.Data.Date)
		assert.Equal(t, expected.TotalAppointments, body.Data.TotalAppointments)
		assert.Equal(t, expected.TotalHours, body.Data.TotalHours)

		mockSvc.AssertExpectations(t)
	})

}

func TestGetWorkloadSummaryByBranch_Controller(t *testing.T) {
	testDate := time.Date(2025, time.May, 14, 0, 0, 0, 0, time.UTC)
	dateStr := testDate.Format("2006-01-02")

    t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(http.MethodGet, "/workloads/summary?date="+dateStr, nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		// stub error when tenantID=0, branchID=0
		mockSvc.
			On("GetWorkloadSummaryByBranch",
				mock.Anything,
				mock.AnythingOfType("time.Time"),
				uint(0),
				uint(0),
			).
			Return(nil, assert.AnError).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/workloads/summary?date="+dateStr, nil)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("EmptyResult_ShouldReturnEmptyData", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		mockSvc.
			On("GetWorkloadSummaryByBranch",
				mock.Anything,
				mock.AnythingOfType("time.Time"),
				uint(0),
				uint(0),
			).
			Return([]barberBookingDto.BranchWorkloadSummary{}, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/workloads/summary?date="+dateStr, nil)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                                   `json:"status"`
			Data   []barberBookingDto.BranchWorkloadSummary `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Empty(t, body.Data)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_NoFilters_ShouldReturnSummaries", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		expected := []barberBookingDto.BranchWorkloadSummary{
			{TenantID: 1, BranchID: 1, NumWorked: 2, TotalBarbers: 3},
			{TenantID: 2, BranchID: 5, NumWorked: 0, TotalBarbers: 1},
		}
		mockSvc.
			On("GetWorkloadSummaryByBranch",
				mock.Anything,
				mock.AnythingOfType("time.Time"),
				uint(0),
				uint(0),
			).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/workloads/summary?date="+dateStr, nil)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                                   `json:"status"`
			Data   []barberBookingDto.BranchWorkloadSummary `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected, body.Data)

		mockSvc.AssertExpectations(t)
	})

	t.Run("FilterByTenant1_ShouldCallServiceWithTenant1", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		expected := []barberBookingDto.BranchWorkloadSummary{
			{TenantID: 1, BranchID: 1, NumWorked: 2, TotalBarbers: 2},
		}
		mockSvc.
			On("GetWorkloadSummaryByBranch",
				mock.Anything,
				mock.AnythingOfType("time.Time"),
				uint(1),
				uint(0),
			).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/workloads/summary?date="+dateStr+"&tenant_id=1", nil)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                                   `json:"status"`
			Data   []barberBookingDto.BranchWorkloadSummary `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected, body.Data)

		mockSvc.AssertExpectations(t)
	})

	t.Run("FilterByBranch5_ShouldCallServiceWithBranch5", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		expected := []barberBookingDto.BranchWorkloadSummary{
			{TenantID: 2, BranchID: 5, NumWorked: 0, TotalBarbers: 1},
		}
		mockSvc.
			On("GetWorkloadSummaryByBranch",
				mock.Anything,
				mock.AnythingOfType("time.Time"),
				uint(0),
				uint(5),
			).
			Return(expected, nil).
			Once()

		req := httptest.NewRequest(http.MethodGet, "/workloads/summary?date="+dateStr+"&branch_id=5", nil)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.StatusCode)

		var body struct {
			Status string                                   `json:"status"`
			Data   []barberBookingDto.BranchWorkloadSummary `json:"data"`
		}
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)

		assert.Equal(t, "success", body.Status)
		assert.Equal(t, expected, body.Data)

		mockSvc.AssertExpectations(t)
	})

	t.Run("InvalidTenantID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		// service ไม่ควรถูกเรียกเลย
		req := httptest.NewRequest(http.MethodGet,
			"/workloads/summary?date="+dateStr+"&tenant_id=foo",
			nil,
		)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertNotCalled(t, "GetWorkloadSummaryByBranch", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})

	t.Run("InvalidBranchID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodGet,
			"/workloads/summary?date="+dateStr+"&branch_id=bar",
			nil,
		)
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

		mockSvc.AssertNotCalled(t, "GetWorkloadSummaryByBranch", mock.Anything, mock.Anything, mock.Anything, mock.Anything)
	})
}

func TestUpsertBarberWorkload_Controller(t *testing.T) {
	validBody := `{"date":"2025-05-14","appointments":3,"hours":5}`


	t.Run("PermissionDenied_ShouldReturn403", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(http.MethodPut, "/barbers/1/workload", nil)
		req.Header.Set("X-Mock-Role", "USER")
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("PermissionDenied_ShouldReturn403# Not branch_admin access", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)
		req := httptest.NewRequest(http.MethodPut, "/barbers/1/workload", nil)
		req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[2]))
		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusForbidden, resp.StatusCode)
	})

	t.Run("BadBarberID_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		// ไม่ต้อง stub service เพราะจะไม่ถึง
		req := httptest.NewRequest(http.MethodPut, "/barbers/abc/workload", strings.NewReader(validBody))
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidJSON_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		req := httptest.NewRequest(http.MethodPut, "/barbers/1/workload", strings.NewReader("not-json"))
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)
	})

	t.Run("InvalidDateFormat_ShouldReturn400", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		badDateBody := `{"date":"14-05-2025","appointments":3,"hours":5}`
		req := httptest.NewRequest(http.MethodPut, "/barbers/1/workload", strings.NewReader(badDateBody))
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

		var body map[string]string
		json.NewDecoder(resp.Body).Decode(&body)
		assert.Equal(t, "Invalid date format. Expect YYYY-MM-DD", body["message"])
	})

	t.Run("ServiceError_ShouldReturn500", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		// Stub service to return error
		mockSvc.
			On("UpsertBarberWorkload",
				mock.Anything,
				uint(1),
				mock.AnythingOfType("time.Time"),
				3,
				5,
			).
			Return(assert.AnError).
			Once()

		req := httptest.NewRequest(http.MethodPut, "/barbers/1/workload", strings.NewReader(validBody))
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

		mockSvc.AssertExpectations(t)
	})

	t.Run("Success_ShouldReturn200", func(t *testing.T) {
		mockSvc := new(MockWorkloadService)
		app := SetupWorkloadTestApp(mockSvc)

		// Stub service to succeed
		mockSvc.
			On("UpsertBarberWorkload",
				mock.Anything,
				uint(1),
				mock.AnythingOfType("time.Time"),
				3,
				5,
			).
			Return(nil).
			Once()

		req := httptest.NewRequest(http.MethodPut, "/barbers/1/workload", strings.NewReader(validBody))
        req.Header.Set("X-Mock-Role", string(barberBookingController.RolesCanGetSummaryBarber[1]))
		req.Header.Set("Content-Type", "application/json")
		resp, err := app.Test(req)
		require.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)

		var body map[string]string
		err = json.NewDecoder(resp.Body).Decode(&body)
		require.NoError(t, err)
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Workload upserted", body["message"])

		mockSvc.AssertExpectations(t)
	})
}
