package barberBookingControllers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	// "strconv"
	"testing"

	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Mock Service ----------
type MockService struct {
	mock.Mock
}

// UpdateService implements barberBookingPort.IServiceService.
func (m *MockService) UpdateService(ctx context.Context, serviceID uint, payload *barberBookingPort.UpdateServiceRequest, file *multipart.FileHeader) (*barberBookingModels.Service, error) {
	panic("unimplemented")
}

// CreateService implements barberBookingPort.IServiceService.
func (m *MockService) CreateService(ctx context.Context, tenantID uint, branchID uint, payload *barberBookingPort.CreateServiceRequest, file *multipart.FileHeader) (*barberBookingModels.Service, error) {
	panic("unimplemented")
}

func (m *MockService) GetAllServices(tenantID uint, branchID uint) ([]barberBookingModels.Service, error) {
	args := m.Called()
	return args.Get(0).([]barberBookingModels.Service), args.Error(1)
}

func (m *MockService) GetServiceByID(id uint) (*barberBookingModels.Service, error) {
	args := m.Called(id)
	return args.Get(0).(*barberBookingModels.Service), args.Error(1)
}

func (m *MockService) DeleteService(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ---------- Test ----------
func setupTestApp(mockSvc barberBookingPort.IServiceService) *fiber.App {
	app := fiber.New()
	controller := barberBookingControllers.NewServiceController(mockSvc)

	app.Get("/services", controller.GetAllServices)
	app.Get("/services/:id", controller.GetServiceByID)
	app.Post("/services", func(c *fiber.Ctx) error {
		// mock auth ด้วย header
		role := c.Get("X-Mock-Role")
		c.Locals("role", role)
		return controller.CreateService(c)
	})

	app.Put("/services/:id", func(c *fiber.Ctx) error {
		// mock auth ด้วย header
		role := c.Get("X-Mock-Role")
		c.Locals("role", role)
		return controller.UpdateService(c)
	})

	return app
}

func generateExpiredToken() string {
	claims := jwt.MapClaims{
		"user_id":   1,
		"role":      "TENANT_ADMIN",
		"tenant_id": 123,
		"exp":       time.Now().Add(-1 * time.Hour).Unix(), // ❌ หมดอายุไปแล้ว
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := os.Getenv("JWT_SECRET")
	signed, _ := token.SignedString([]byte(secret))
	return signed
}

func TestGetAllServices_Success(t *testing.T) {
	mockSvc := new(MockService)
	mockSvc.On("GetAllServices").Return([]barberBookingModels.Service{
		{Name: "Haircut"}, {Name: "Shampoo"},
	}, nil)

	app := setupTestApp(mockSvc)
	req := httptest.NewRequest("GET", "/services", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetAllServices_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	mockSvc.On("GetAllServices").Return([]barberBookingModels.Service{}, errors.New("db error"))

	app := setupTestApp(mockSvc)
	req := httptest.NewRequest("GET", "/services", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestGetServiceByID_Success(t *testing.T) {
	mockSvc := new(MockService)
	mockSvc.On("GetServiceByID", uint(1)).Return(&barberBookingModels.Service{ID: 1, Name: "Beard Trim"}, nil)

	app := setupTestApp(mockSvc)
	req := httptest.NewRequest("GET", "/services/1", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}

func TestGetServiceByID_InvalidID(t *testing.T) {
	mockSvc := new(MockService) // won't be called
	app := setupTestApp(mockSvc)

	req := httptest.NewRequest("GET", "/services/abc", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)
}

func TestGetServiceByID_NotFound(t *testing.T) {
	mockSvc := new(MockService)
	mockSvc.On("GetServiceByID", uint(999)).Return((*barberBookingModels.Service)(nil), nil)

	app := setupTestApp(mockSvc)
	req := httptest.NewRequest("GET", "/services/999", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 404, resp.StatusCode)
}

func TestGetServiceByID_ServiceError(t *testing.T) {
	mockSvc := new(MockService)
	mockSvc.On("GetServiceByID", uint(2)).Return((*barberBookingModels.Service)(nil), errors.New("db error"))

	app := setupTestApp(mockSvc)
	req := httptest.NewRequest("GET", "/services/2", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 500, resp.StatusCode)
}

func TestCreateService(t *testing.T) {
	type request struct {
		Name     string  `json:"name"`
		Duration int     `json:"duration"`
		Price    float64 `json:"price"`
	}
	mockSvc := new(MockService)
	mockSvc.On("CreateService", mock.MatchedBy(func(s *barberBookingModels.Service) bool {
		return s.Name == "Haircut" && s.Duration == 30
	})).Return(nil).Once()

	app := fiber.New()
	ctrl := barberBookingControllers.NewServiceController(mockSvc)
	app.Post("/services", func(c *fiber.Ctx) error {
		// middleware จำลอง auth
		c.Locals("role", c.Get("X-Mock-Role"))
		c.Locals("tenant_id", uint(123))
		return ctrl.CreateService(c)
	})

	app.Use(func(c *fiber.Ctx) error {
		c.Locals("role", c.Get("X-Mock-Role"))
		c.Locals("tenant_id", uint(1)) // ← ตรงนี้เพิ่ม tenant_id
		return c.Next()
	})

	// CASE 1: Success (TENANT_ADMIN)
	t.Run("CreateService_Success_TenantAdmin", func(t *testing.T) {
		body := request{Name: "Haircut", Duration: 30, Price: 200}
		mockSvc.On("CreateService", mock.AnythingOfType("*models.Service")).Return(nil).Once()

		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 201, resp.StatusCode)
	})

	// CASE 2: Forbidden (role USER) ไม่มีสิทธิ
	t.Run("CreateService_Forbidden_RoleMismatch", func(t *testing.T) {
		body := request{Name: "Beard", Duration: 15, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "USER") //  ไม่ใช่ TENANT

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	// CASE 3: Unauthorized (no JWT / role) // ไม่ได้ login หรือไม่ก็ token หมดอายุ
	t.Run("CreateService_Unauthorized_NoToken", func(t *testing.T) {
		body := request{Name: "Beard", Duration: 15, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		//  ไม่มี role

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode) // หรือ 401 ขึ้นกับ middleware จริงของคุณ
	})

	// CASE 4: Invalid Input (duration = 0)
	t.Run("CreateService_InvalidInput", func(t *testing.T) {
		body := request{Name: "Invalid", Duration: 0, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode) // หรือ 422 แล้วแต่คุณจัดการ
	})

	// CASE 5: Service Layer Error DB ล่ม หรือไม่ก็ service api พัง
	t.Run("CreateService_InternalError", func(t *testing.T) {
		body := request{Name: "Massage", Duration: 30, Price: 500}
		mockSvc.On("CreateService", mock.MatchedBy(func(arg interface{}) bool {
			svc, ok := arg.(*barberBookingModels.Service)
			return ok && svc.Name != ""
		})).Return(errors.New("mock internal error")).Once()

		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})

	// CASE 6: ไม่มีชื่อ ชื่องว่างเปล่า
	t.Run("CreateService_EmptyName", func(t *testing.T) {
		body := request{Name: " ", Duration: 30, Price: 200}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
	//CASE 7: ราคาติดลบ
	t.Run("CreateService_NegativePrice", func(t *testing.T) {
		body := request{Name: "Haircut", Duration: 30, Price: -50}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
	//CASE 8:ราคาเป็น 0 บริการฟรี
	t.Run("CreateService_ZeroPrice", func(t *testing.T) {
		body := request{Name: "FreeService", Duration: 10, Price: 0}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})
	//CASE 9:ชื่อยาวเกิน
	t.Run("CreateService_TooLongName", func(t *testing.T) {
		longName := strings.Repeat("A", 101)
		body := request{Name: longName, Duration: 20, Price: 150}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	//CASE 10: เวลาติดลบ
	t.Run("CreateService_NegativeDuration", func(t *testing.T) {
		body := request{Name: "Deep Clean", Duration: -10, Price: 300}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	//CASE 11: token หมดอายุ
	t.Run("CreateService_ExpiredToken", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/services", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "token=expired.jwt.token.here") // จำลอง token ที่หมดอายุ

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode) // หรือ 403 ขึ้นกับ middleware ของคุณ
	})

	//CASE 12: token แปลก ๆ
	t.Run("CreateService_InvalidToken", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/services", bytes.NewReader([]byte(`{}`)))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Cookie", "token=invalid.token") // จำลอง token ผิดโครงสร้าง

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})
}

func TestUpdateService(t *testing.T) {
	mockSvc := new(MockService)
	app := setupTestApp(mockSvc)

	ctrl := barberBookingControllers.NewServiceController(mockSvc)
	app.Put("/services/:id", ctrl.UpdateService)

	type request struct {
		Name     string `json:"name"`
		Duration int    `json:"duration"`
		Price    int    `json:"price"`
	}

	t.Run("UpdateService_Success_TenantAdmin", func(t *testing.T) {
		svcID := uint(1)
		body := request{Name: "New Name", Duration: 40, Price: 300}

		// STEP 1: mock GetServiceByID
		mockSvc.On("GetServiceByID", svcID).Return(&barberBookingModels.Service{
			ID:       svcID,
			Name:     "Old Name",
			Duration: 30,
			Price:    200,
		}, nil).Once()

		// STEP 2: mock UpdateService
		mockSvc.On("UpdateService", svcID, mock.MatchedBy(func(arg interface{}) bool {
			svc, ok := arg.(*barberBookingModels.Service)
			return ok && svc.Name != ""
		})).Return(&barberBookingModels.Service{
			ID:       svcID,
			Name:     "New Name",
			Duration: 40,
			Price:    300,
		}, nil).Once()

		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", fmt.Sprintf("/services/%d", svcID), bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, resp.StatusCode)
	})

	t.Run("UpdateService_Forbidden_RoleMismatch", func(t *testing.T) {
		body := request{Name: "Haircut", Duration: 20, Price: 150}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "USER")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	t.Run("UpdateService_Unauthorized_NoToken", func(t *testing.T) {
		body := request{Name: "Haircut", Duration: 20, Price: 150}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	t.Run("UpdateService_InvalidInput", func(t *testing.T) {
		body := request{Name: "", Duration: 0, Price: -10}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_InternalError", func(t *testing.T) {
		body := request{Name: "Massage", Duration: 30, Price: 500}

		// 🔧 เพิ่ม mock สำหรับ GetServiceByID ด้วย
		mockSvc.On("GetServiceByID", uint(1)).Return(&barberBookingModels.Service{
			ID:       1,
			Name:     "Old Name",
			Duration: 20,
			Price:    100,
		}, nil).Maybe()

		mockSvc.On("UpdateService", uint(1), mock.Anything).Return((*barberBookingModels.Service)(nil), errors.New("mock error")).Maybe()

		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 500, resp.StatusCode)
	})

	t.Run("UpdateService_EmptyName", func(t *testing.T) {
		body := request{Name: " ", Duration: 20, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_ZeroPrice", func(t *testing.T) {
		body := request{Name: "Cut", Duration: 20, Price: 0}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_NegativePrice", func(t *testing.T) {
		body := request{Name: "Cut", Duration: 20, Price: -100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_NegativeDuration", func(t *testing.T) {
		body := request{Name: "Cut", Duration: -5, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_TooLongName", func(t *testing.T) {
		body := request{Name: strings.Repeat("A", 200), Duration: 15, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_InvalidIDFormat", func(t *testing.T) {
		body := request{Name: "Cut", Duration: 15, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/abc", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 400, resp.StatusCode)
	})

	t.Run("UpdateService_NotFound", func(t *testing.T) {
		body := request{Name: "Cut", Duration: 15, Price: 100}
		mockSvc.On("GetServiceByID", mock.AnythingOfType("uint")).Return((*barberBookingModels.Service)(nil), errors.New("not found")).Maybe()

		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/999", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 404, resp.StatusCode)
	})

	t.Run("UpdateService_ExpiredToken", func(t *testing.T) {
		body := request{Name: "Cut", Duration: 15, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "token", Value: generateExpiredToken()})

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	t.Run("UpdateService_InvalidToken", func(t *testing.T) {
		body := request{Name: "Cut", Duration: 15, Price: 100}
		reqBody, _ := json.Marshal(body)
		req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
		req.Header.Set("Content-Type", "application/json")
		req.AddCookie(&http.Cookie{Name: "token", Value: "bad_token"})

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, resp.StatusCode)
	})

	t.Run("UpdateService_ConcurrentRequests", func(t *testing.T) {
		var wg sync.WaitGroup

		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()

				// 🔁 Create isolated mock service per goroutine
				mockSvc := new(MockService)

				mockSvc.On("GetServiceByID", mock.AnythingOfType("uint")).Return(&barberBookingModels.Service{
					ID:       1,
					Name:     "Original",
					Duration: 20,
					Price:    100,
				}, nil).Maybe()

				mockSvc.On("UpdateService", mock.AnythingOfType("uint"), mock.Anything).
					Return(&barberBookingModels.Service{
						ID:       1,
						Name:     fmt.Sprintf("Service %d", index),
						Duration: 30,
						Price:    200,
					}, nil).Maybe()

				app := setupTestApp(mockSvc)

				body := request{Name: fmt.Sprintf("Service %d", index), Duration: 30, Price: 200}
				reqBody, _ := json.Marshal(body)

				req := httptest.NewRequest("PUT", "/services/1", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Mock-Role", "TENANT_ADMIN")

				resp, err := app.Test(req)
				assert.NoError(t, err)
				assert.Equal(t, 200, resp.StatusCode)
			}(i)
		}
		wg.Wait()
	})

}
