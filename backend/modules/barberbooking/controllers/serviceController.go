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
)

type ServiceController struct {
	ServiceService barberBookingPort.IServiceService
}

func NewServiceController(svc barberBookingPort.IServiceService) *ServiceController {
	return &ServiceController{
		ServiceService: svc,
	}
}

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
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
}

// loop check role function
func IsAuthorizedRole(role string, allowed []coreModels.RoleName) bool {
	for _, r := range allowed {
		if string(r) == role {
			return true
		}
	}
	return false
}

func (ctrl *ServiceController) CreateService(c *fiber.Ctx) error {
	//  ตรวจสิทธิ์ผู้ใช้งาน
	roleStr, ok := c.Locals("role").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "Unauthorized",
		})
	}
	if !IsAuthorizedRole(roleStr, RolesCanManageService) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// ✅ แปลง request body
	var payload barberBookingModels.Service
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// ✅ ตรวจสอบความถูกต้องของข้อมูล
	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" || len(payload.Name) > 100 || payload.Duration <= 0 || payload.Price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service input",
		})
	}

	// ✅ เรียก service
	if err := ctrl.ServiceService.CreateService(&payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create service",
			"error":   err.Error(),
		})
	}

	// ✅ ส่ง response สำเร็จ
	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Service created",
	})
}

func (ctrl *ServiceController) UpdateService(c *fiber.Ctx) error {
	// 1. ตรวจสอบ role
	roleStr, ok := c.Locals("role").(string)
	if !ok || !IsAuthorizedRole(roleStr, RolesCanManageService) {
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

func (ctrl *ServiceController) DeleteService(c *fiber.Ctx) error {
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
