package Core_controllers
import (
    "errors"
	"strconv"
	"strings"

    "github.com/gofiber/fiber/v2"
	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"
)

type TenantController struct {
    TenantService corePort.ITenant
}

func NewTenantController(svc corePort.ITenant) *TenantController {
    return &TenantController{TenantService: svc}
}

// CreateTenant godoc
// @Summary      สร้าง Tenant ใหม่
// @Description  สร้าง Tenant ใหม่โดยระบุชื่อ (name) และโดเมน (domain)
// @Tags         Tenant
// @Accept       json
// @Produce      json
// @Param        body  body      corePort.CreateTenantInput  true  "ข้อมูลสำหรับสร้าง Tenant (name, domain)"
// @Success      201   {object}  map[string]interface{}      "คืนค่า status และข้อมูล tenant ที่สร้าง"
// @Failure      400   {object}  map[string]string           "Malformed JSON หรือข้อมูลไม่ถูกต้อง"
// @Failure      500   {object}  map[string]string           "เกิดข้อผิดพลาดระหว่างสร้าง tenant"
// @Router       /core/tenant-route/create [post]
// @Security     ApiKeyAuth
func (ctrl *TenantController) CreateTenant(c *fiber.Ctx) error {
    var req corePort.CreateTenantInput
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Malformed JSON"})
    }

    tenant, err := ctrl.TenantService.CreateTenant(c.Context(), corePort.CreateTenantInput{
        Name:   req.Name,
        Domain: req.Domain,
    })
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantInput):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":err.Error()})
        case errors.Is(err, coreServices.ErrDomainTaken):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":err.Error()})
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status":"error","message":"Failed to create tenant"})
        }
    }

    return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status":"success","data":tenant})
}


// GetTenantByID godoc
// @Summary      ดึงข้อมูล Tenant ตาม ID
// @Description  ดึงรายละเอียดของ Tenant หนึ่งรายการตามรหัสที่ระบุ
// @Tags         Tenant
// @Produce      json
// @Param        id   path      int  true  "รหัส Tenant"
// @Success      200  {object}  map[string]interface{}  "คืนค่า status และข้อมูล Tenant ใน key `data`"
// @Failure      400  {object}  map[string]string       "Invalid tenant ID หรือ รหัสไม่ถูกต้อง"
// @Failure      404  {object}  map[string]string       "ไม่พบ Tenant ตามรหัสที่ระบุ"
// @Failure      500  {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึงข้อมูล Tenant"
// @Router       /core/tenant-route/:id [get]
// @Security     ApiKeyAuth
func (ctrl *TenantController) GetTenantByID(c *fiber.Ctx) error {
    // 1. Parse and validate tenant ID from path
    idParam := c.Params("id")
    id64, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil || id64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }
    tenantID := uint(id64)

    // 2. Call service
    tenant, err := ctrl.TenantService.GetTenantByID(c.Context(), tenantID)
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to retrieve tenant",
            })
        }
    }

    // 3. Success
    return c.JSON(fiber.Map{
        "status": "success",
        "data":   tenant,
    })
}


// ListTenants godoc
// @Summary      ดึงรายการ Tenant ทั้งหมด
// @Description  ดึงรายการ Tenant โดยสามารถกรองเฉพาะที่ active ได้ผ่าน query parameter `active` (default = true)
// @Tags         Tenant
// @Produce      json
// @Param        active  query     bool  false  "กรองเฉพาะ Tenant ที่ active (true/false), default = true"
// @Success      200     {object}  map[string]interface{}  "คืนค่า status และ array ของ Tenant ใน key `data`"
// @Failure      400     {object}  map[string]string       "Invalid `active` query parameter"
// @Failure      500     {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึงรายการ Tenant"
// @Router       /core/tenant-route [get]
// @Security     ApiKeyAuth
func (ctrl *TenantController) ListTenants(c *fiber.Ctx) error {
    // 1. Parse optional ?active query (default = true)
    onlyActive := true
    if v := c.Query("active"); v != "" {
        b, err := strconv.ParseBool(v)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": "Invalid `active` query parameter",
            })
        }
        onlyActive = b
    }

    // 2. Call service
    tenants, err := ctrl.TenantService.ListTenants(c.Context(), onlyActive)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to fetch tenants",
        })
    }

    // 3. Return result
    return c.JSON(fiber.Map{
        "status": "success",
        "data":   tenants,
    })
}


// UpdateTenant godoc
// @Summary      แก้ไขข้อมูล Tenant ตาม ID
// @Description  อัปเดตชื่อ โดเมน หรือสถานะ active ของ Tenant หน่วยตามรหัสที่ระบุ
// @Tags         Tenant
// @Accept       json
// @Produce      json
// @Param        id    path      int                          true  "รหัส Tenant"
// @Param        body  body      corePort.UpdateTenantInput  true  "ข้อมูลที่ต้องการอัปเดต (name, domain, isActive)"
// @Success      200   {object}  map[string]string            "คืนค่า status และข้อความยืนยันการอัปเดต"
// @Failure      400   {object}  map[string]string            "Invalid tenant ID หรือ malformed JSON หรือ validation error"
// @Failure      404   {object}  map[string]string            "ไม่พบ Tenant ตามรหัสที่ระบุ"
// @Failure      500   {object}  map[string]string            "เกิดข้อผิดพลาดระหว่างอัปเดต Tenant"
// @Router       /core/tenant-route/:id [put]
// @Security     ApiKeyAuth
func (ctrl *TenantController) UpdateTenant(c *fiber.Ctx) error {
    // 1. Parse and validate ID
    idParam := c.Params("id")
    id64, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil || id64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }
    tenantID := uint(id64)

    // 2. Parse body
    var req corePort.UpdateTenantInput
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Malformed JSON",
        })
    }

    // 3. Trim whitespace if provided
    if req.Name != nil {
        *req.Name = strings.TrimSpace(*req.Name)
    }
    if req.Domain != nil {
        *req.Domain = strings.TrimSpace(*req.Domain)
    }

    // 4. Build port input
    input := corePort.UpdateTenantInput{
        ID:       tenantID,
        Name:     req.Name,
        Domain:   req.Domain,
        IsActive: req.IsActive,
    }

    // 5. Call service
    if err := ctrl.TenantService.UpdateTenant(c.Context(), input); err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrInvalidTenantInput),
             errors.Is(err, coreServices.ErrDomainTaken):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to update tenant",
            })
        }
    }

    // 6. Success
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Tenant updated",
    })
}

// DeleteTenant godoc
// @Summary      ลบ Tenant ตาม ID
// @Description  ลบ Tenant ที่ระบุด้วยรหัส (ต้องไม่มีการใช้งานอยู่ หรือจะเกิด conflict error หากยังมีทรัพยากรภายในใช้ Tenant นี้อยู่)
// @Tags         Tenant
// @Produce      json
// @Param        id   path      int  true  "รหัส Tenant"
// @Success      200  {object}  map[string]string  "คืนค่า status และข้อความยืนยันการลบ"
// @Failure      400  {object}  map[string]string  "Invalid tenant ID หรือ รหัสไม่ถูกต้อง"
// @Failure      404  {object}  map[string]string  "ไม่พบ Tenant ตามรหัสที่ระบุ"
// @Failure      409  {object}  map[string]string  "ไม่สามารถลบ Tenant ได้ เนื่องจากยังมีการใช้งานภายในระบบ"
// @Failure      500  {object}  map[string]string  "เกิดข้อผิดพลาดระหว่างการลบ Tenant"
// @Router       /core/tenant-route/:id [delete]
// @Security     ApiKeyAuth
func (ctrl *TenantController) DeleteTenant(c *fiber.Ctx) error {
    // 1. Parse and validate ID
    idParam := c.Params("id")
    id64, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil || id64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }

    // 2. Call service
    err = ctrl.TenantService.DeleteTenant(c.Context(), uint(id64))
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        case errors.Is(err, coreServices.ErrTenantInUse):
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "status":  "error",
                "message": "Cannot delete tenant: still in use",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to delete tenant",
            })
        }
    }

    // 3. Success
    return c.JSON(fiber.Map{
        "status":  "success",
        "message": "Tenant deleted",
    })
}