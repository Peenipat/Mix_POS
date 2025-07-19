package barberBookingController

import (
	"context"
	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
)

type BarberWorkloadController struct {
	Service barberBookingPort.IbarberWorkload
}

func NewBarberWorkloadController(service barberBookingPort.IbarberWorkload) *BarberWorkloadController {
	return &BarberWorkloadController{
		Service: service,
	}
}

var RolesCanManageWorkload = []coreModels.RoleName{
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameTenant,
}

// GetWorkloadByBarber godoc
// @Summary      ดึงปริมาณงานของช่างตัดผม
// @Description  คืนข้อมูล BarberWorkload สำหรับช่างตัดผมที่ระบุ ภายใน Tenant ที่กำหนด โดยกรองตามวันที่ (YYYY-MM-DD) ถ้าไม่ส่ง date จะใช้วันที่ปัจจุบัน
// @Tags         Workload
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                              true   "รหัส Tenant"
// @Param        barber_id  path      uint                              true   "รหัส Barber"
// @Param        date       query     string                            false  "วันที่กรอง (YYYY-MM-DD), default = today"
// @Success      200        {object}  barberBookingModels.BarberWorkload  "คืนค่า status success และข้อมูล workload ใน key `data`"
// @Failure      400        {object}  map[string]string                 "Invalid Barber ID หรือรูปแบบวันที่ไม่ถูกต้อง"
// @Failure      403        {object}  map[string]string                 "Permission denied"
// @Failure      500        {object}  map[string]string                 "Failed to fetch workload"
// @Router      /tenants/:tenant_id/barberworkload/barbers/:barber_id [get]
// @Security     ApiKeyAuth
func (ctrl *BarberWorkloadController) GetWorkloadByBarber(c *fiber.Ctx) error {

    roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageWorkload) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}
    
	// 1. Parse barber id จาก URL param
	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Barber ID",
		})
	}

	dateStr := c.Query("date", "")
	var dateParsed time.Time
	if dateStr == "" {
		now := time.Now()
		dateParsed = time.Date(now.Year(), now.Month(), now.Day(),
			0, 0, 0, 0, now.Location())
	} else {
		dateParsed, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid date format. Expect YYYY-MM-DD",
			})
		}
	}

	

	workload, err := ctrl.Service.GetWorkloadByBarber(c.Context(), barberID, dateParsed)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch workload",
			"error":   err.Error(),
		})
	}

	if workload == nil {
		workload = &barberBookingModels.BarberWorkload{
			BarberID:          barberID,
			Date:              dateParsed,
			TotalAppointments: 0,
			TotalHours:        0,
		}
	}

	return c.JSON(fiber.Map{
		"status": "success",
		"data":   workload,
	})

}


// @Summary      สร้าง/อัปเดตปริมาณงานของช่างตัดผม
// @Description  Upsert (insert หรือ update) BarberWorkload ตามรหัสช่างและวันที่
// @Tags         Workload
// @Accept       json
// @Produce      json
// @Param        barber_id    path      uint                              true  "รหัส Barber"
// @Param        body         body      barberBookingPort.UpsertBarberWorkloadRequest  true  "Payload สำหรับ upsert ปริมาณงาน (date, appointments, hours)"
// @Success      200          {object}  map[string]string               "คืนค่า status success และข้อความยืนยันการ upsert"
// @Failure      400          {object}  map[string]string               "Invalid barber_id หรือ invalid JSON body หรือ invalid date format"
// @Failure      403          {object}  map[string]string               "Permission denied"
// @Failure      500          {object}  map[string]string               "Failed upsert workload"
// @Router       /tenants/:tenant_id/barberworkload/barbers/:barber_id [post]
// @Security     ApiKeyAuth
func (ctrl *BarberWorkloadController) UpsertBarberWorkload(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageWorkload) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}
	// 1. parse barber_id
	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid barber_id"})
	}

	// 2. parse body
	var payload struct {
		Date         string `json:"date"`
		Appointments int    `json:"appointments"`
		Hours        int    `json:"hours"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
	}

	// 3. parse date
	dateParsed, err := time.Parse("2006-01-02", payload.Date)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid date format. Expect YYYY-MM-DD"})
	}

	// 4. call service
	if err := ctrl.Service.UpsertBarberWorkload(
		context.Background(),
		barberID,
		dateParsed,
		payload.Appointments,
		payload.Hours,
	); err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "Failed upsert workload", "error": err.Error()})
	}

	// 5. success
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "message": "Workload upserted"})
}

var RolesCanGetSummaryBarber = []coreModels.RoleName{
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameBranchAdmin,
}

// GetWorkloadSummaryByBranch godoc
// @Summary      สรุปปริมาณงานรายสาขา
// @Description  คืนรายการสรุป BarberWorkload grouped by Branch ตามวันที่ และตัวกรอง Tenant/Branch (ถ้ามี)  
// @Tags         Workload
// @Accept       json
// @Produce      json
// @Param        date        query     string   true   "วันที่กรอง (YYYY-MM-DD)"
// @Param        tenant_id   query     uint     false  "กรองเฉพาะ Tenant (optional)"
// @Param        branch_id   query     uint     false  "กรองเฉพาะ Branch (optional)"
// @Success      200         {object}  map[string][]barberBookingModels.BarberWorkload  "คืนค่า status success และ array ของสรุปปริมาณงานใน key `data`"
// @Failure      400         {object}  map[string]string   "Missing date หรือ Invalid date/tenant_id/branch_id"
// @Failure      403         {object}  map[string]string   "Permission denied"
// @Failure      500         {object}  map[string]string   "Failed to fetch workload summary"
// @Router       /tenants/:tenant_id/barberworkload/branches/:branch_id/summary [get]
// @Security     ApiKeyAuth
func (ctrl *BarberWorkloadController) GetWorkloadSummaryByBranch(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanGetSummaryBarber) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}
	// Parse date
	dateStr := c.Query("date", "")
	if dateStr == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing date query param"})
	}
	dateParsed, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid date format. Expect YYYY-MM-DD"})
	}

	// Parse optional tenant_id
	var tenantID uint
	if tid := c.Query("tenant_id", ""); tid != "" {
		parsed, err := strconv.ParseUint(tid, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
		}
		tenantID = uint(parsed)
	}

	// Parse optional branch_id
	var branchID uint
	if bid := c.Query("branch_id", ""); bid != "" {
		parsed, err := strconv.ParseUint(bid, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).
				JSON(fiber.Map{"status": "error", "message": "Invalid branch_id"})
		}
		branchID = uint(parsed)
	}

	// Call service with filters
	summaries, err := ctrl.Service.GetWorkloadSummaryByBranch(c.Context(), dateParsed, tenantID, branchID)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "Failed to fetch workload summary", "error": err.Error()})
	}

	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": summaries})
}
