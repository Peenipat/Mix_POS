package barberBookingController

import (
	"github.com/gofiber/fiber/v2"
	helperFunc "myapp/modules/barberbooking"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
)

type WorkingHourController struct {
	Service barberBookingPort.IWorkingHourService
}

func NewWorkingHourController(service barberBookingPort.IWorkingHourService) *WorkingHourController {
	return &WorkingHourController{
		Service: service,
	}
}

// GetWorkingHours ดึงเวลาเปิด-ปิดของสาขา
func (ctrl *WorkingHourController) GetWorkingHours(c *fiber.Ctx) error {
	// 1. Parse branch_id จาก URL param
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}

	// 2. เรียก service
	workingHours, err := ctrl.Service.GetWorkingHours(c.Context(), branchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch working hours",
			"error":   err.Error(),
		})
	}

	if len(workingHours) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No working hours found for this branch",
		})
	}

	// 3. Success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Working hours retrieved",
		"data":    workingHours,
	})
}

var RolesCanManageWorkingHour = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameBranchAdmin,
}

// UpdateWorkingHours อัปเดตเวลาเปิด-ปิดของสาขา (หลายวัน)
func (ctrl *WorkingHourController) UpdateWorkingHours(c *fiber.Ctx) error {

	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}

	var input []barberBookingDto.WorkingHourInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if len(input) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No working hours provided",
		})
	}

	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageWorkingHour) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	err = ctrl.Service.UpdateWorkingHours(c.Context(), branchID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update working hours",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Working hours updated",
	})
}

// CreateWorkingHours สร้างวันทำการใหม่ (เฉพาะ 1 วัน)
func (ctrl *WorkingHourController) CreateWorkingHours(c *fiber.Ctx) error {
	// ตรวจสอบสิทธิ์
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageWorkingHour) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// ดึง branch_id จาก path
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}

	// Parse body เป็น WorkingHourInput
	var input barberBookingDto.WorkingHourInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate เบื้องต้น
	if input.Weekday < 0 || input.Weekday > 6 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid weekday (must be 0-6)",
		})
	}
	if input.StartTime.IsZero() || input.EndTime.IsZero() {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Start time and end time are required",
		})
	}

	// Call Service
	err = ctrl.Service.CreateWorkingHours(c.Context(), branchID, input)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create working hour",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Working hour created",
	})
}



