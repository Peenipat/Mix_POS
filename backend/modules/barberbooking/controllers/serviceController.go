// barberbooking/controller/service_controller.go
package controller

import (
	"net/http"
	"strconv"
	"strings"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingService "myapp/modules/barberbooking/services"
	coreModels "myapp/modules/core/models"

	"github.com/gofiber/fiber/v2"
)

type ServiceController struct {
	ServiceService barberBookingService.IServiceService
}

func NewServiceController(svc barberBookingService.IServiceService) *ServiceController {
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
	// ✅ ตรวจสิทธิ์ผู้ใช้งาน
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
	// ✅ เช็คสิทธิ์ก่อน
	role := c.Locals("role")
	roleStr, ok := role.(string)
	if !ok || (roleStr != string(coreModels.RoleNameTenant) && roleStr != string(coreModels.RoleNameTenantAdmin)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// ✅ แปลง id
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	// ✅ แปลง body
	var payload barberBookingModels.Service
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// ✅ เรียก service
	updated, err := ctrl.ServiceService.UpdateService(uint(id), &payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update service",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Service updated",
		"data":    updated,
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
