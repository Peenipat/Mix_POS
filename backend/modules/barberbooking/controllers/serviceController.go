// barberbooking/controller/service_controller.go
package barberBookingController

import (
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
	"strings"

	helperFunc "myapp/modules/barberbooking"
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

// GetAllServices godoc
// @Summary      ดึงรายการบริการทั้งหมด
// @Description  คืนรายการ Service ทั้งหมดในระบบ
// @Tags         Service
// @Accept       json
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "คืนค่า status success, message และ array ของ Service ใน key `data`"
// @Failure      500  {object}  map[string]string       "Failed to fetch services"
// @Router       /tenants/:tenant_id/branch/:branch_id/services [get]
func (ctrl *ServiceController) GetAllServices(c *fiber.Ctx) error {
	// 1) อ่าน tenant_id จาก URL
	tenantParam := c.Params("tenant_id")
	tenantID, err := strconv.ParseUint(tenantParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant_id",
			"error":   err.Error(),
		})
	}

	// 2) อ่าน branch_id จาก URL
	branchParam := c.Params("branch_id")
	branchID, err := strconv.ParseUint(branchParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch_id",
			"error":   err.Error(),
		})
	}

	// 3) เรียก service layer เพื่อดึงข้อมูลเฉพาะ tenant & branch นี้
	services, err := ctrl.ServiceService.GetAllServices(uint(tenantID), uint(branchID))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch services",
			"error":   err.Error(),
		})
	}

	// 4) คืนผลลัพธ์กลับไป
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
// @Router       /tenants/:tenant_id/branch/:branch_id/services/:service_id [get]
func (ctrl *ServiceController) GetServiceByID(c *fiber.Ctx) error {
	idParam := c.Params("service_id")
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
	coreModels.RoleNameBranchAdmin,
}

// CreateService godoc
// @Summary      สร้างบริการใหม่
// @Description  เพิ่ม Service ใหม่ภายใต้ Tenant ของผู้ใช้ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant หรือ TenantAdmin)
// @Tags         Service
// @Accept       multipart/form-data
// @Produce      json
// @Param        name         formData  string  true   "ชื่อบริการ"
// @Param        description  formData  string  false  "คำอธิบาย"
// @Param        duration     formData  int     true   "ระยะเวลา (นาที)"
// @Param        price        formData  number  true   "ราคา"
// @Param        file         formData  file    false  "รูปภาพประกอบบริการ (optional)"
// @Success      201  {object}  map[string]interface{}  "คืนค่า status success พร้อมข้อมูลบริการที่ถูกสร้าง"
// @Failure      400  {object}  map[string]string       "Invalid request body หรือ tenant ID หรือข้อมูลไม่ถูกต้อง"
// @Failure      401  {object}  map[string]string       "Unauthorized"
// @Failure      403  {object}  map[string]string       "Permission denied"
// @Failure      500  {object}  map[string]string       "Failed to create service"
// @Router       /tenants/{tenant_id}/branch/{branch_id}/services [post]
// @Security     ApiKeyAuth
func (ctrl *ServiceController) CreateService(c *fiber.Ctx) error {
	// 1) ตรวจสอบ role
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

	// 2) ดึง tenant_id และ branch_id จาก path param
	tidStr := c.Params("tenant_id")
	tid64, err := strconv.ParseUint(tidStr, 10, 64)
	if err != nil || tid64 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant ID",
		})
	}
	tenantID := uint(tid64)

	bidStr := c.Params("branch_id")
	bid64, err := strconv.ParseUint(bidStr, 10, 64)
	if err != nil || bid64 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}
	branchID := uint(bid64)

	// 3) รับค่าจาก form-data
	name := strings.TrimSpace(c.FormValue("name"))
	description := strings.TrimSpace(c.FormValue("description"))
	durationStr := c.FormValue("duration")
	priceStr := c.FormValue("price")

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid duration",
		})
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid price",
		})
	}

	if name == "" || len(name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service name",
		})
	}

	// 4) รับไฟล์รูป (optional)
	file, _ := c.FormFile("file")

	// 5) Prepare payload
	payload := &barberBookingPort.CreateServiceRequest{
		Name:        name,
		Description: description,
		Duration:    duration,
		Price:       price,
	}

	// 6) เรียก service layer
	created, err := ctrl.ServiceService.CreateService(c.Context(), tenantID, branchID, payload, file)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create service",
			"error":   err.Error(),
		})
	}

	// 7) Success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Service created",
		"data":    created,
	})
}

// UpdateService godoc
// @Summary      แก้ไขบริการ
// @Description  แก้ไขข้อมูลบริการและรูปภาพ (เฉพาะผู้มีสิทธิ์)
// @Tags         Service
// @Accept       multipart/form-data
// @Produce      json
// @Param        service_id   path      int     true  "Service ID"
// @Param        name         formData  string  true  "ชื่อบริการ"
// @Param        description  formData  string  false "คำอธิบาย"
// @Param        duration     formData  int     true  "ระยะเวลา (นาที)"
// @Param        price        formData  number  true  "ราคา"
// @Param        file         formData  file    false "รูปภาพใหม่ (optional)"
// @Success      200 {object} map[string]interface{}
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      404 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /services/{service_id} [put]
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

	// 2. อ่าน service_id จาก path param
	idParam := c.Params("service_id")
	serviceID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || serviceID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid service ID",
		})
	}

	// 3. อ่าน form values
	name := strings.TrimSpace(c.FormValue("name"))
	description := strings.TrimSpace(c.FormValue("description"))
	durationStr := c.FormValue("duration")
	priceStr := c.FormValue("price")

	duration, err := strconv.Atoi(durationStr)
	if err != nil || duration <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid duration",
		})
	}

	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil || price <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid price",
		})
	}

	if name == "" || len(name) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid name",
		})
	}

	// 4. รับไฟล์ (optional)
	file, _ := c.FormFile("file")

	// 5. สร้าง payload DTO
	payload := &barberBookingPort.UpdateServiceRequest{
		Name:        name,
		Description: description,
		Duration:    duration,
		Price:       price,
	}

	// 6. เรียก service layer
	updated, err := ctrl.ServiceService.UpdateService(c.Context(), uint(serviceID), payload, file)
	if err != nil {
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Service not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update service",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Service updated",
		"data":    updated,
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
// @Router       /services/:service_id [delete]
// @Security     ApiKeyAuth
func (ctrl *ServiceController) DeleteService(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageService) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}
	idParam := c.Params("service_id")
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
