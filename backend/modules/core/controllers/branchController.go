package Core_controllers

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	helperFunc "myapp/modules/core"
	coreModels "myapp/modules/core/models"
	coreServices "myapp/modules/core/services"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// BranchController handles HTTP requests for branches
type BranchController struct {
	BranchService coreServices.BranchPort
}
func NewBranchController(svc coreServices.BranchPort) *BranchController {
	return &BranchController{BranchService: svc}
}
var RolesCanManageBranch = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameTenant,
}




// CreateBranch godoc
// @Summary      สร้างสาขาใหม่ เช็คตัวเชื่อมกับ Tenant
// @Description  ใช้สร้างสาขาย่อยใหม่โดยระบุชื่อสาขาและที่อยู่เท่านั้น ยังไม่มีการเชื่อมโยงกับ Tenant
// @Tags         Branch
// @Accept       json
// @Produce      json
// @Param        body  body      corePort.CreateBranchInput  true  "ข้อมูลสำหรับสร้างสาขา (ชื่อ, ที่อยู่) — tenant_id จะถูกดึงจาก context ของผู้ใช้"
// @Success      201   {object}  map[string]interface{}  "คืนค่า status, message และข้อมูลสาขาที่สร้าง"
// @Failure      400   {object}  map[string]string       "ข้อมูลส่งมาไม่ถูกต้อง, ขาด tenant ID หรือสาขานี้มีอยู่แล้ว"
// @Failure      403   {object}  map[string]string       "ไม่มีสิทธิ์เข้าถึง"
// @Failure      500   {object}  map[string]string       "สร้างสาขาไม่สำเร็จ"
// @Router       /core/tenants/:tenant_id/branches [post]
// @Security     ApiKeyAuth
func (ctrl *BranchController) CreateBranch(c *fiber.Ctx) error {
	// Authorization: ตรวจสอบ role
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr,RolesCanManageBranch) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// Parse and bind JSON body
	var payload coreModels.Branch
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// Inject tenant_id from locals
	tenantID, ok := c.Locals("tenant_id").(uint)
	if !ok || tenantID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing or invalid tenant ID",
		})
	}
	payload.TenantID = tenantID

	// Call service to create branch
	if err := ctrl.BranchService.CreateBranch(&payload); err != nil {
		if strings.Contains(err.Error(), "already exists") {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create branch",
		})
	}

	// Success
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Branch created",
		"data":    payload,
	})
}

// GetBranches godoc
// @Summary      ดึงรายการสาขาทั้งหมด
// @Description  ดึงข้อมูลสาขาทั้งหมดโดยไม่กรองตาม Tenant (ต้องมีสิทธิ์ RolesCanManageBranch)
// @Tags         Branch
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "คืนค่า status และ array ของสาขาใน key `data`"
// @Failure      403  {object}  map[string]string       "ไม่มีสิทธิ์เข้าถึง"
// @Failure      500  {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึงข้อมูลสาขา"
// @Router       /core/branches/all [get]
// @Security     ApiKeyAuth
func (ctrl *BranchController) GetBranches(c *fiber.Ctx) error {
    // 1. Authorization
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBranch) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // 2. Call global service (no tenant filter)
    branches, err := ctrl.BranchService.GetAllBranches()
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to fetch branches",
        })
    }

    // 3. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   branches,
    })
}


// GetBranchByID godoc
// @Summary      ดึงข้อมูลสาขาตาม ID
// @Description  ดึงข้อมูลสาขาเดียวตามรหัสสาขาที่ระบุ (ต้องมีสิทธิ์ RolesCanManageBranch)
// @Tags         Branch
// @Produce      json
// @Param        id   path      int  true  "รหัสสาขา"
// @Success      200  {object}  map[string]interface{}  "คืนค่า status และข้อมูลสาขาใน key `data`"
// @Failure      400  {object}  map[string]string       "Invalid branch ID หรือ รหัสไม่ถูกต้อง"
// @Failure      403  {object}  map[string]string       "ไม่มีสิทธิ์เข้าถึง"
// @Failure      404  {object}  map[string]string       "ไม่พบสาขาตามรหัสที่ระบุ"
// @Failure      500  {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึงข้อมูลสาขา"
// @Router       /core/branch/:id [get]
// @Security     ApiKeyAuth
func (ctrl *BranchController) GetBranchByID(c *fiber.Ctx) error {
    // 1. Authorization
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBranch) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // 2. Parse ID from path
    idParam := c.Params("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil || id == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid branch ID",
        })
    }

    // 3. Call service
    branch, err := ctrl.BranchService.GetBranchByID(uint(id))
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrBranchNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Branch not found",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to retrieve branch",
            })
        }
    }

    // 4. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   branch,
    })
}

// UpdateBranch godoc
// @Summary      แก้ไขชื่อสาขาตาม ID
// @Description  อัปเดตชื่อสาขาที่ระบุด้วยรหัสสาขา (ต้องมีสิทธิ์ RolesCanManageBranch)
// @Tags         Branch
// @Accept       json
// @Produce      json
// @Param        id    path      int                         true  "รหัสสาขา"
// @Param        body  body      corePort.UpdateBranchInput  true  "ข้อมูลที่ต้องการอัปเดต (name)"
// @Success      200   {object}  map[string]interface{}      "คืนค่า status และข้อมูลสาขาที่อัปเดต"
// @Failure      400   {object}  map[string]string           "Invalid branch ID, malformed JSON, หรือ validation error"
// @Failure      403   {object}  map[string]string           "ไม่มีสิทธิ์เข้าถึง"
// @Failure      404   {object}  map[string]string           "ไม่พบสาขาตามรหัสที่ระบุ"
// @Failure      500   {object}  map[string]string           "เกิดข้อผิดพลาดระหว่างอัปเดตสาขา"
// @Router       /core/tenants/:tenant_id/branches/:id [put]
// @Security     ApiKeyAuth
func (ctrl *BranchController) UpdateBranch(c *fiber.Ctx) error {
    // 1. Authorization
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBranch) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // 2. Parse ID from path
    idParam := c.Params("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil || id == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid branch ID",
        })
    }

    // 3. Parse request body
    var dto struct {
        Name string `json:"name"`
    }
    if err := c.BodyParser(&dto); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Malformed JSON",
        })
    }

    // 4. Prepare model for update
    branch := &coreModels.Branch{
        ID:   uint(id),
        Name: dto.Name,
    }

    // 5. Call service method (not ctrl.UpdateBranch)
    if err := ctrl.BranchService.UpdateBranch(branch); err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case strings.Contains(err.Error(), "branch name is required"):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, gorm.ErrRecordNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Branch not found",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to update branch",
            })
        }
    }

    // 6. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   branch,
    })
}

// DeleteBranch godoc
// @Summary      ลบสาขาตาม ID
// @Description  ลบสาขาที่ระบุด้วยรหัสสาขา (ต้องมีสิทธิ์ RolesCanManageBranch)  
// @Tags         Branch
// @Produce      json
// @Param        id   path      int  true  "รหัสสาขา"
// @Success      200  {object}  map[string]string  "คืนค่า status และข้อความยืนยันการลบ"
// @Failure      400  {object}  map[string]string  "Invalid branch ID หรือ รหัสไม่ถูกต้อง"
// @Failure      403  {object}  map[string]string  "ไม่มีสิทธิ์เข้าถึง"
// @Failure      404  {object}  map[string]string  "ไม่พบสาขาตามรหัสที่ระบุ"
// @Failure      409  {object}  map[string]string  "ไม่สามารถลบสาขาได้ เนื่องจากมีการใช้งานอยู่"
// @Failure      500  {object}  map[string]string  "เกิดข้อผิดพลาดระหว่างการลบสาขา"
// @Router       /core/tenants/:tenant_id/branches/:id [delete]
// @Security     ApiKeyAuth
func (ctrl *BranchController) DeleteBranch(c *fiber.Ctx) error {
    // 1. Authorization
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBranch) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // 2. Parse ID from path
    idParam := c.Params("id")
    id, err := strconv.ParseUint(idParam, 10, 64)
    if err != nil || id == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid branch ID",
        })
    }

    // 3. Call service
    if err := ctrl.BranchService.DeleteBranch(uint(id)); err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrBranchNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Branch not found",
            })
        case errors.Is(err, coreServices.ErrBranchInUse):
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "status":  "error",
                "message": "Branch cannot be deleted: in use",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to delete branch",
            })
        }
    }

    // 4. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Branch deleted",
    })
}

// GetBranchesByTenantID godoc
// @Summary      ดึงรายการสาขาของ Tenant ปัจจุบัน
// @Description  ดึงข้อมูลสาขาทั้งหมดเฉพาะสำหรับ Tenant ที่ผู้ใช้ล็อกอินอยู่ (tenant_id มาจาก context หลังตรวจสอบ token)
// @Tags         Branch
// @Produce      json
// @Success      200  {object}  map[string]interface{}  "คืนค่า status และ array ของสาขาใน key `data`"
// @Failure      400  {object}  map[string]string       "Missing or invalid tenant ID หรือ InvalidTenantID"
// @Failure      403  {object}  map[string]string       "ไม่มีสิทธิ์เข้าถึง"
// @Failure      404  {object}  map[string]string       "ไม่พบ Tenant ตามรหัสที่ระบุ"
// @Failure      500  {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึงข้อมูลสาขา"
// @Router       /core/tenants/:tenant_id/branches [get]
// @Security     ApiKeyAuth
func (ctrl *BranchController) GetBranchesByTenantID(c *fiber.Ctx) error {
    // 1. Authorization
    roleStr, ok := c.Locals("role").(string)
    if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageBranch) {
        return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
            "status":  "error",
            "message": "Permission denied",
        })
    }

    // 2. Tenant ID from context
    tidAny := c.Locals("tenant_id")
    tenantID, ok := tidAny.(uint)
    if !ok || tenantID == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Missing or invalid tenant ID",
        })
    }

    // 3. Call service
    branches, err := ctrl.BranchService.GetBranchesByTenantID(tenantID)
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
                "message": "Failed to fetch branches",
            })
        }
    }

    // 4. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   branches,
    })
}


