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