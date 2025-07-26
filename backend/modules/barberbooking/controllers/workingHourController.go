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

// GetWorkingHours godoc
// @Summary      ดึงเวลาเปิด-ปิดของสาขา
// @Description  ดึงรายการ WorkingHour ทั้งหมดสำหรับสาขาที่ระบุ
// @Tags         WorkingHour
// @Accept       json
// @Produce      json
// @Param        branch_id   path      uint   true  "รหัสสาขา"
// @Param        tenant_id   path      uint   true  "รหัสสาขา"
// @Success      200         {object}  map[string]interface{}  "คืนค่า status, message และ array ของ WorkingHour"
// @Failure      400         {object}  map[string]string       "Invalid branch ID"
// @Failure      404         {object}  map[string]string       "No working hours found for this branch"
// @Failure      500         {object}  map[string]string       "Failed to fetch working hours"
// @Router       /tenants/{tenant_id}/workinghour/branches/{branch_id} [get]
// @Security     ApiKeyAuth
func (ctrl *WorkingHourController) GetWorkingHours(c *fiber.Ctx) error {
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}

	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant ID",
		})
	}

	// 2. เรียก service
	workingHours, err := ctrl.Service.GetWorkingHours(c.Context(), branchID, tenantID)
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

// UpdateWorkingHours godoc
// @Summary      อัปเดตเวลาเปิด-ปิดของสาขา (หลายวัน)
// @Description  รับรายการ WorkingHourInput หลายรายการเพื่อนำมาอัปเดตเวลาเปิด-ปิดของสาขาที่ระบุ (วันในสัปดาห์, เปิด/ปิด)
// @Tags         WorkingHour
// @Accept       json
// @Produce      json
// @Param 		branch_id 	path 	uint true "รหัสสาขา"
// @Param 		tenant_id 	path 	uint true "รหัสผู้เช่า"                           
// @Param        body        body      []barberBookingDto.WorkingHourInput  true  "Array ของ WorkingHourInput — ตัวอย่าง: [{\"weekday\":1,\"open_time\":\"09:00\",\"close_time\":\"17:00\"}, ...]"
// @Success      200         {object}  map[string]string  "คืนค่า status และข้อความยืนยันการอัปเดต"
// @Failure      400         {object}  map[string]string  "Invalid branch ID, ไม่มีข้อมูลใน body หรือ JSON ไม่ถูกต้อง"
// @Failure      403         {object}  map[string]string  "Permission denied"
// @Failure      500         {object}  map[string]string  "Failed to update working hours"
// @Router       /tenants/{tenant_id}/workinghour/branches/{branch_id}  [put]
// @Security     ApiKeyAuth
func (ctrl *WorkingHourController) UpdateWorkingHours(c *fiber.Ctx) error {
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}

	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant ID",
		})
	}

	var inputArray []barberBookingDto.WorkingHourInput

	// ลอง parse เป็น array ก่อน
	err = c.BodyParser(&inputArray)
	if err != nil || len(inputArray) == 0 {
		// ถ้าไม่สำเร็จ ลอง parse เป็น object เดียว
		var singleInput barberBookingDto.WorkingHourInput
		if err := c.BodyParser(&singleInput); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid request body",
			})
		}
		inputArray = []barberBookingDto.WorkingHourInput{singleInput}
	}

	if len(inputArray) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No working hours provided",
		})
	}

	// ตรวจสอบสิทธิ์
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageWorkingHour) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// ส่งไป Service
	if err := ctrl.Service.UpdateWorkingHours(c.Context(), branchID, tenantID, inputArray); err != nil {
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


// CreateWorkingHours godoc
// @Summary      สร้างวันทำการใหม่ (เฉพาะ 1 วัน)
// @Description  สร้าง WorkingHour ใหม่สำหรับสาขาที่ระบุ กำหนดวันในสัปดาห์ (0=อาทิตย์,...,6=เสาร์) และเวลาเปิด-ปิด
// @Tags         WorkingHour
// @Accept       json
// @Produce      json
// @Param        branch_id   path      uint                        true  "รหัสสาขา"
// @Param        body        body      barberBookingDto.WorkingHourInput  true  "Payload สำหรับสร้างวันทำการ (weekday, start_time, end_time)"
// @Success      201         {object}  map[string]string  "คืนค่า status และข้อความยืนยันการสร้าง"
// @Failure      400         {object}  map[string]string  "Invalid branch ID, weekday ไม่ถูกต้อง หรือ JSON ไม่ถูกต้อง"
// @Failure      403         {object}  map[string]string  "Permission denied"
// @Failure      500         {object}  map[string]string  "Failed to create working hour"
// @Router       /tenants/:tenant_id/workinghour/branches/:branch_id [post]
// @Security     ApiKeyAuth
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


// GetAvailableSlots godoc
// @Summary      ดึง slot เวลาที่ว่างของสาขา
// @Description  ดึงช่วงเวลาที่เปิดให้บริการในแต่ละวัน โดยสามารถกรองได้ตามวันนี้, สัปดาห์นี้ หรือเดือนนี้ พร้อมทั้งกรองช่วงเวลา (เช่น 09:00-12:30)
// @Tags         WorkingHour
// @Accept       json
// @Produce      json
// @Param        tenant_id    path      uint     true   "รหัสผู้เช่า (Tenant ID)"
// @Param        branch_id    path      uint     true   "รหัสสาขา (Branch ID)"
// @Param        filter       query     string   false  "ประเภทการกรองวันที่: week (สัปดาห์นี้), month (เดือนนี้), หรือเว้นไว้เพื่อดึงแค่วันนี้"
// @Param        from_time    query     string   false  "เวลาเริ่มต้นที่ต้องการกรอง (รูปแบบ: HH:mm เช่น 09:00)"
// @Param        to_time      query     string   false  "เวลาสิ้นสุดที่ต้องการกรอง (รูปแบบ: HH:mm เช่น 12:30)"
// @Success      200          {object}  map[string]interface{}  "คืนค่า status และ slot รายวัน เช่น { '2025-07-22': ['09:00', '09:30', ...] }"
// @Failure      400          {object}  map[string]string  "พารามิเตอร์ไม่ถูกต้อง เช่น tenant_id หรือ branch_id ไม่ถูกต้อง"
// @Failure      500          {object}  map[string]string  "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์"
// @Router       /tenants/{tenant_id}/workinghour/branches/{branch_id}/slots [get]
// @Security     ApiKeyAuth
func (ctrl *WorkingHourController) GetAvailableSlots(c *fiber.Ctx) error {
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil || branchID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid branch_id",
		})
	}

	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil || tenantID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "invalid tenant_id",
		})
	}

	filter := c.Query("filter") 

	// เวลาเริ่ม–สิ้นสุด
	fromTime := c.Query("from_time")
	toTime := c.Query("to_time")

	var fromPtr, toPtr *string
	if fromTime != "" {
		fromPtr = &fromTime
	}
	if toTime != "" {
		toPtr = &toTime
	}

	// เรียก service
	slots, err := ctrl.Service.GetAvailableSlots(c.Context(), branchID, tenantID, filter, fromPtr, toPtr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   slots,
	})
}

