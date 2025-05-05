package barberBookingControllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"

	// "strconv"
	"testing"

	barberBookingControllers "myapp/modules/barberbooking/controllers"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingServices "myapp/modules/barberbooking/services"

	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// ---------- Mock Service ----------
type MockService struct {
	mock.Mock
}

func (m *MockService) GetAllServices() ([]barberBookingModels.Service, error) {
	args := m.Called()
	return args.Get(0).([]barberBookingModels.Service), args.Error(1)
}

func (m *MockService) GetServiceByID(id uint) (*barberBookingModels.Service, error) {
	args := m.Called(id)
	return args.Get(0).(*barberBookingModels.Service), args.Error(1)
}

func (m *MockService) CreateService(service *barberBookingModels.Service) error {
	args := m.Called(service)
	return args.Error(0)
}

func (m *MockService) UpdateService(id uint, service *barberBookingModels.Service) (*barberBookingModels.Service, error) {
	args := m.Called(id, service)
	return args.Get(0).(*barberBookingModels.Service), args.Error(1)
}

func (m *MockService) DeleteService(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

// ---------- Test ----------
func setupTestApp(mockSvc barberBookingServices.IServiceService) *fiber.App {
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

	return app
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
		return ctrl.CreateService(c)
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
		req.Header.Set("X-Mock-Role", "USER") // ❌ ไม่ใช่ TENANT

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
		// ❌ ไม่มี role

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

	//CASE 13:จำลอง 20 คนสร้างพร้อมกัน
	t.Run("CreateService_ConcurrentRequests", func(t *testing.T) {
		var wg sync.WaitGroup
	
		for i := 0; i < 20; i++ {
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
	
				// ✅ mock แยกในแต่ละ thread
				mockSvc := new(MockService)
				mockSvc.On("CreateService", mock.AnythingOfType("*barberBookingModels.Service")).Return(nil).Maybe()

	
				app := setupTestApp(mockSvc)
	
				body := request{Name: fmt.Sprintf("Service%d", index), Duration: 15, Price: 100}
				reqBody, _ := json.Marshal(body)
				req := httptest.NewRequest("POST", "/services", bytes.NewReader(reqBody))
				req.Header.Set("Content-Type", "application/json")
				req.Header.Set("X-Mock-Role", "TENANT_ADMIN")
	
				resp, err := app.Test(req)
				assert.NoError(t, err)
				assert.Equal(t, 201, resp.StatusCode)
			}(i)
		}
		wg.Wait()
	})
	
	
	
}
