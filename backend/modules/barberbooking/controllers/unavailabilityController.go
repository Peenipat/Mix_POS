package barberBookingController

import (
	// "context"
	"errors"
	"gorm.io/gorm"
	"time"

	"github.com/gofiber/fiber/v2"
	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
)

type UnavailabilityController struct {
	Service barberBookingPort.IUnavailabilitySerivce
}

func NewUnavailabilityController(service barberBookingPort.IUnavailabilitySerivce) *UnavailabilityController {
	return &UnavailabilityController{
		Service: service,
	}
}

var RolesCanManageUnavailability = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameBranchAdmin,
}

// CreateUnavailability godoc
// @Summary      สร้างวันที่ไม่ว่าง
// @Description  เพิ่ม Unavailability ระบุวันที่และเลือกได้ว่าจะปิดช่างหรือสาขา (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Unavailability
// @Accept       json
// @Produce      json
// @Param        body  body      barberBookingModels.Unavailability  true  "Payload สำหรับสร้าง Unavailability (Date, BarberID หรือ BranchID)"
// @Success      201   {object}  barberBookingModels.Unavailability  "คืนค่า status success, message และข้อมูล Unavailability ใน key `data`"
// @Failure      400   {object}  map[string]string                  "Invalid request body หรือ Date และ BarberID/BranchID จำเป็นต้องระบุ"
// @Failure      403   {object}  map[string]string                  "Permission denied"
// @Failure      409   {object}  map[string]string                  "Unavailability already exists for this date"
// @Failure      500   {object}  map[string]string                  "Failed to create unavailability"
// @Router       /tenants/:tenant_id/unavailability [post]
// @Security     ApiKeyAuth
func (ctrl *UnavailabilityController) CreateUnavailability(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageUnavailability) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	var input barberBookingModels.Unavailability
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Validate required fields
	if input.Date.IsZero() || (input.BarberID == nil && input.BranchID == nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Date and either BarberID or BranchID are required",
		})
	}

	created, err := ctrl.Service.CreateUnavailability(c.Context(), &input)
	if err != nil {
		if errors.Is(err, errors.New("unavailability already exists for this date")) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create unavailability",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Unavailability created",
		"data":    created,
	})
}

// GetUnavailabilitiesByBranch godoc
// @Summary      ดึงรายการวันไม่ว่างตามสาขา
// @Description  คืน Unavailability ทั้งหมดสำหรับสาขาที่ระบุ ในช่วงวันที่กำหนด (ต้องระบุ query params `from` และ `to` ในรูปแบบ YYYY-MM-DD)
// @Tags         Unavailability
// @Accept       json
// @Produce      json
// @Param        branch_id  path      uint     true  "รหัส Branch"
// @Param        from       query     string   true  "วันที่เริ่มต้น (YYYY-MM-DD)"
// @Param        to         query     string   true  "วันสิ้นสุด (YYYY-MM-DD)"
// @Success      200        {object}  map[string]interface{}  "คืนค่า status success, message และ array ของ Unavailability ใน key `data`"
// @Failure      400        {object}  map[string]string       "Invalid branch ID หรือ Missing/invalid `from` or `to` date format"
// @Failure      404        {object}  map[string]string       "No unavailabilities found in the given date range"
// @Failure      500        {object}  map[string]string       "Failed to fetch unavailabilities"
// @Router       /tenants/:tenant_id/unavailability/branches/:branch_id [get]
func (ctrl *UnavailabilityController) GetUnavailabilitiesByBranch(c *fiber.Ctx) error {

	//  Optional: ตรวจสอบ Role (ถ้ามีข้อกำหนดเฉพาะ)
	// roleStr, ok := c.Locals("role").(string)
	// if !ok || roleStr == "" {
	// 	return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
	// 		"status":  "error",
	// 		"message": "Permission denied",
	// 	})
	// }
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
		})
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Query parameters 'from' and 'to' are required (format: YYYY-MM-DD)",
		})
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid 'from' date format. Use YYYY-MM-DD",
		})
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid 'to' date format. Use YYYY-MM-DD",
		})
	}

	data, err := ctrl.Service.GetUnavailabilitiesByBranch(c.Context(), branchID, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch unavailabilities",
			"error":   err.Error(),
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No unavailabilities found in the given date range",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Unavailabilities retrieved",
		"data":    data,
	})
}


// GetUnavailabilitiesByBarber godoc
// @Summary      ดึงรายการวันไม่ว่างตามช่างตัดผม
// @Description  คืน Unavailability ทั้งหมดสำหรับ Barber ที่ระบุ ในช่วงวันที่กำหนด (ต้องระบุ query params `from` และ `to` ในรูปแบบ YYYY-MM-DD)
// @Tags         Unavailability
// @Accept       json
// @Produce      json
// @Param        barber_id  path      uint     true  "รหัส Barber"
// @Param        from       query     string   true  "วันที่เริ่มต้น (YYYY-MM-DD)"
// @Param        to         query     string   true  "วันสิ้นสุด (YYYY-MM-DD)"
// @Success      200        {object}  map[string]interface{}  "คืนค่า status success, message และ array ของ Unavailability ใน key `data`"
// @Failure      400        {object}  map[string]string       "Invalid barber ID หรือ Missing/invalid `from` or `to` date format"
// @Failure      404        {object}  map[string]string       "No unavailabilities found"
// @Failure      500        {object}  map[string]string       "Failed to fetch unavailabilities"
// @Router      /tenants/:tenant_id/unavailabilitybarbers/:barber_id [get]
// @Security     ApiKeyAuth
func (ctrl *UnavailabilityController) GetUnavailabilitiesByBarber(c *fiber.Ctx) error {
	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid barber ID",
		})
	}

	fromStr := c.Query("from")
	toStr := c.Query("to")
	if fromStr == "" || toStr == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Query parameters 'from' and 'to' are required (format: YYYY-MM-DD)",
		})
	}

	from, err := time.Parse("2006-01-02", fromStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid 'from' date format. Use YYYY-MM-DD",
		})
	}

	to, err := time.Parse("2006-01-02", toStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid 'to' date format. Use YYYY-MM-DD",
		})
	}

	data, err := ctrl.Service.GetUnavailabilitiesByBarber(c.Context(), barberID, from, to)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch unavailabilities",
			"error":   err.Error(),
		})
	}

	if len(data) == 0 {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "No unavailabilities found",
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Unavailabilities retrieved",
		"data":    data,
	})
}

// UpdateUnavailability godoc
// @Summary      อัปเดตวันไม่ว่าง
// @Description  แก้ไขฟิลด์ของ Unavailability ตามที่ส่งมา (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Unavailability
// @Accept       json
// @Produce      json
// @Param        id     path      uint                        true  "รหัส Unavailability"
// @Param        body   body      map[string]interface{}      true  "ฟิลด์ที่ต้องการอัปเดต (เช่น date, barber_id, branch_id)"
// @Success      200    {object}  map[string]string           "คืนค่า status success และข้อความยืนยันการอัปเดต"
// @Failure      400    {object}  map[string]string           "Invalid unavailability ID, Invalid request body หรือ No update fields provided"
// @Failure      403    {object}  map[string]string           "Permission denied"
// @Failure      404    {object}  map[string]string           "Unavailability not found"
// @Failure      500    {object}  map[string]string           "Failed to update unavailability"
// @Router       /tenants/:tenant_id/unavailability/:id [put]
// @Security     ApiKeyAuth
func (ctrl *UnavailabilityController) UpdateUnavailability(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageUnavailability) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	unavailID, err := helperFunc.ParseUintParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid unavailability ID",
		})
	}

	var updates map[string]interface{}
	if err := c.BodyParser(&updates); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	if len(updates) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "No update fields provided",
		})
	}

	err = ctrl.Service.UpdateUnavailability(c.Context(), unavailID, updates)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Unavailability not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update unavailability",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Unavailability updated",
	})
}

// DeleteUnavailability godoc
// @Summary      ลบวันไม่ว่าง
// @Description  ลบ Unavailability ตามรหัสที่ระบุ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Unavailability
// @Accept       json
// @Produce      json
// @Param        id   path      uint   true  "รหัส Unavailability"
// @Success      200  {object}  map[string]string  "คืนค่า status success และข้อความยืนยันการลบ"
// @Failure      400  {object}  map[string]string  "Invalid unavailability ID"
// @Failure      403  {object}  map[string]string  "Permission denied"
// @Failure      404  {object}  map[string]string  "Unavailability not found"
// @Failure      500  {object}  map[string]string  "Failed to delete unavailability"
// @Router       /tenants/:tenant_id/unavailability/:id  [delete]
// @Security     ApiKeyAuth
func (ctrl *UnavailabilityController) DeleteUnavailability(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageUnavailability) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	unavailID, err := helperFunc.ParseUintParam(c, "id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid unavailability ID",
		})
	}

	err = ctrl.Service.DeleteUnavailability(c.Context(), unavailID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "Unavailability not found",
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete unavailability",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Unavailability deleted",
	})
}
