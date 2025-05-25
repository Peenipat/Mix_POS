package Core_controllers

import (
	"errors"
	"strconv"

	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"
    coreModels "myapp/modules/core/models"

	"github.com/gofiber/fiber/v2"
)

// TenantUserController handles M2M endpoints between tenants and users.
type TenantUserController struct {
    Service corePort.ITenantUser
}

// NewTenantUserController constructs a new controller.
func NewTenantUserController(svc corePort.ITenantUser) *TenantUserController {
    return &TenantUserController{Service: svc}
}

var (
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrUserAlreadyAssigned = errors.New("user already assigned to tenant")
	ErrUserNotAssigned = errors.New("user not assigned to tenant")
    ErrNoTenantsAssigned  = errors.New("no tenants assigned to user")
    ErrUserNotFound       = errors.New("user not found")
)

// AddUserToTenant handles POST /tenants/:tenant_id/users/:user_id
func (ctrl *TenantUserController) AddUserToTenant(c *fiber.Ctx) error {
    tidParam := c.Params("tenant_id")
    tid64, err := strconv.ParseUint(tidParam, 10, 64)
    if err != nil || tid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }

    // 3. Parse user_id
    uidParam := c.Params("user_id")
    uid64, err := strconv.ParseUint(uidParam, 10, 64)
    if err != nil || uid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid user ID",
        })
    }

    // 4. Call service
    err = ctrl.Service.AddUserToTenant(c.Context(), uint(tid64), uint(uid64))
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID),
             errors.Is(err, coreServices.ErrInvalidUserID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        case errors.Is(err, coreServices.ErrUserNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "User not found",
            })
        case errors.Is(err, coreServices.ErrUserAlreadyAssigned):
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "status":  "error",
                "message": "User already assigned to this tenant",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to assign user to tenant",
            })
        }
    }

    // 5. Success
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status":  "success",
        "message": "User assigned to tenant",
    })
}

func (ctrl *TenantUserController) RemoveUserFromTenant(c *fiber.Ctx) error {
    // 1. Parse tenant_id
    tidParam := c.Params("tenant_id")
    tid64, err := strconv.ParseUint(tidParam, 10, 64)
    if err != nil || tid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant ID",
        })
    }
    tenantID := uint(tid64)

    // 2. Parse user_id
    uidParam := c.Params("user_id")
    uid64, err := strconv.ParseUint(uidParam, 10, 64)
    if err != nil || uid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid user ID",
        })
    }
    userID := uint(uid64)

    // 3. Call service
    err = ctrl.Service.RemoveUserFromTenant(c.Context(), tenantID, userID)
    if err != nil {
        switch {
        case errors.Is(err, coreServices.ErrInvalidTenantID),
             errors.Is(err, coreServices.ErrInvalidUserID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrTenantNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "Tenant not found",
            })
        case errors.Is(err, coreServices.ErrUserNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "User not found",
            })
        case errors.Is(err, coreServices.ErrUserNotAssigned):
            return c.Status(fiber.StatusConflict).JSON(fiber.Map{
                "status":  "error",
                "message": "User is not assigned to this tenant",
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to remove user from tenant",
            })
        }
    }

    // 4. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "User removed from tenant",
    })
}

func (ctrl *TenantUserController) ListTenantsByUser(c *fiber.Ctx) error {
    // 1. Parse and validate user_id
    uidParam := c.Params("user_id")
    uid64, err := strconv.ParseUint(uidParam, 10, 64)
    if err != nil || uid64 == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid user ID",
        })
    }
    userID := uint(uid64)

    // 2. Call service
    tenants, err := ctrl.Service.ListTenantsByUser(c.Context(), userID)
    if err != nil {
        // 3. Handle service errors
        switch {
        case errors.Is(err, coreServices.ErrInvalidUserID):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": err.Error(),
            })
        case errors.Is(err, coreServices.ErrUserNotFound):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": "User not found",
            })
        case errors.Is(err, coreServices.ErrNoTenantsAssigned):
            // เปลี่ยนจาก 404+error เป็น 200+success พร้อม data ว่าง
            return c.JSON(fiber.Map{
                "status": "success",
                "data":   []coreModels.Tenant{},
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to fetch tenants",
            })
        }
    }

    // 4. Success
    return c.JSON(fiber.Map{
        "status": "success",
        "data":   tenants,
    })
}
