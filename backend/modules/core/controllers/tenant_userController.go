package Core_controllers

import (
	"errors"
	"strconv"

	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"

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

// AddUserToTenant handles POST /tenants/:tenant_id/users/:user_id
func (ctrl *TenantUserController) AddUserToTenant(c *fiber.Ctx) error {
    // 1. (Optional) Authorization: check role from c.Locals("role") if needed

    // 2. Parse tenant_id
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