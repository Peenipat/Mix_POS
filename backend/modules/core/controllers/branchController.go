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




// CreateBranch handles POST /branches
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

// Line SAAS_Admin
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


