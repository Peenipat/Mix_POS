// barberbooking/controller/service_controller.go
package barberBookingController

import (
	"net/http"
	"strconv"
	"strings"
	"github.com/gofiber/fiber/v2"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	helperFunc "myapp/modules/barberbooking"
)

type ServiceController struct {
	ServiceService barberBookingPort.IServiceService
}

func NewServiceController(svc barberBookingPort.IServiceService) *ServiceController {
	return &ServiceController{
		ServiceService: svc,
	}
}

// GetAllServices godoc
// @Summary      ดึงรายการบริการทั้งหมด
// @Description  คืนรายการ Service ทั้งหมดในระบบ
// @Tags         Service
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "คืนค่า status success, message และ array ของ Service ใน key `data`"
// @Failure      500  {object}  map[string]string       "Failed to fetch services"
// @Router       /tenants/:tenant_id/services [get]
func (ctrl *ServiceController) GetAllServices(c *fiber.Ctx) error {
	services, err := ctrl.ServiceService.GetAllServices()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch services",
			"error":   err.Error(),
		})
	}
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Services retrieved",
		"data":    services,
	})
}

// GetServiceByID godoc
// @Summary      ดึงข้อมูลบริการตาม ID
// @Description  คืนข้อมูล Service ตามรหัสที่ระบุ
// @Tags         Service
// @Accept       json
// @Produce      json
// @Param        id   path      int                           true  "รหัส Service"
// @Success      200  {object}  map[string]interface{}        "คืนค่า status success, message และข้อมูล Service ใน key `data`"
// @Failure      400  {object}  map[string]string             "Invalid service ID"
// @Failure      404  {object}  map[string]string             "Service not found"
// @Failure      500  {object}  map[string]string             "Failed to fetch service"
// @Router       /tenants/:tenant_id/services/:service_id [get]
func (ctrl *ServiceController) GetServiceByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	service, err := ctrl.ServiceService.GetServiceByID(uint(id))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch service",
			"error":   err.Error(),
		})
	}
	if service == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Service not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Service retrieved",
		"data":    service,
	})
}

// Role ที่มีสิทธิที่จะ create , update , delete service
var RolesCanManageService = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
}

// CreateService godoc
// @Summary      สร้างบริการใหม่
// @Description  เพิ่ม Service ใหม่ภายใต้ Tenant ของผู้ใช้ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant หรือ TenantAdmin)
// @Tags         Service
// @Accept       json
// @Produce      json
// @Param        body  body      barberBookingModels.Service  true  "Payload สำหรับสร้าง Service (Name, Duration, Price)"
// @Success      201   {object}  map[string]string            "คืนค่า status success และข้อความยืนยันการสร้าง"
// @Failure      400   {object}  map[string]string            "Invalid request body หรือ Invalid tenant ID หรือ Invalid service input"
// @Failure      401   {object}  map[string]string            "Unauthorized"
// @Failure      403   {object}  map[string]string            "Permission denied"
// @Failure      500   {object}  map[string]string            "Failed to create service"
// @Router       /tenants/:tenant_id/services [post]
// @Security     ApiKeyAuth
func (ctrl *ServiceController) CreateService(c *fiber.Ctx) error {
	// ตรวจสิทธิ์ผู้ใช้งาน
	roleStr, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	if !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageService) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// แปลง request body
	var payload barberBookingModels.Service
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// ดึง tenant_id จาก Locals แล้วกำหนดให้ payload
	tenantID, ok := c.Locals("tenant_id").(uint)
	if !ok || tenantID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant ID",
		})
	}
	payload.TenantID = tenantID

	// ตรวจสอบความถูกต้องของข้อมูล
	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" || len(payload.Name) > 100 || payload.Duration <= 0 || payload.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service input",
		})
	}

	// เรียก service
	if err := ctrl.ServiceService.CreateService(&payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create service",
			"error":   err.Error(),
		})
	}

	// ส่ง response สำเร็จ
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Service created",
	})
}


// UpdateService godoc
// @Summary      แก้ไขข้อมูลบริการ
// @Description  อัปเดต Service ตามรหัสที่ระบุ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant หรือ TenantAdmin)
// @Tags         Service
// @Accept       json
// @Produce      json
// @Param        id    path      uint                          true  "รหัส Service"
// @Param        body  body      barberBookingModels.Service   true  "Payload สำหรับอัปเดต Service (Name, Duration, Price)"
// @Success      200   {object}  barberBookingModels.Service   "คืนค่า status success, message และข้อมูล Service ที่อัปเดตใน key `data`"
// @Failure      400   {object}  map[string]string             "Invalid service ID หรือ Invalid request body หรือ Invalid service input"
// @Failure      403   {object}  map[string]string             "Permission denied"
// @Failure      404   {object}  map[string]string             "Service not found"
// @Failure      500   {object}  map[string]string             "Failed to update service"
// @Router       /tenants/:tenant_id/services/:service_id [put]
// @Security     ApiKeyAuth
func (ctrl *ServiceController) UpdateService(c *fiber.Ctx) error {
	// 1. ตรวจสอบ role
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageService) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// 2. อ่าน ID จาก path param
	idParam := c.Params("id")
	serviceID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || serviceID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	// 3. parse body
	var payload barberBookingModels.Service
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// 4. validate
	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" || len(payload.Name) > 100 || payload.Duration <= 0 || payload.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service input",
		})
	}

	// 5. หา service เดิมก่อน (เผื่อไม่เจอ)
	existingService, err := ctrl.ServiceService.GetServiceByID(uint(serviceID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Service not found",
		})
	}

	// 6. อัปเดตค่าจาก payload
	existingService.Name = payload.Name
	existingService.Duration = payload.Duration
	existingService.Price = payload.Price

	// 7. call service
	updatedService, err := ctrl.ServiceService.UpdateService(uint(serviceID), existingService)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update service",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Service updated",
		"data":    updatedService,
	})
}

// DeleteService godoc
// @Summary      ลบบริการ
// @Description  ลบ Service ตามรหัสที่ระบุ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant หรือ TenantAdmin)
// @Tags         Service
// @Accept       json
// @Produce      json
// @Param        id   path      uint               true  "รหัส Service"
// @Success      200  {object}  map[string]string  "คืนค่า status success และข้อความยืนยันการลบ"
// @Failure      400  {object}  map[string]string  "Invalid service ID"
// @Failure      403  {object}  map[string]string  "Permission denied"
// @Failure      500  {object}  map[string]string  "Failed to delete service"
// @Router       /tenants/:tenant_id/services/:service_id [delete]
// @Security     ApiKeyAuth
func (ctrl *ServiceController) DeleteService(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageService) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	if err := ctrl.ServiceService.DeleteService(uint(id)); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete service",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Service deleted successfully",
	})
}
