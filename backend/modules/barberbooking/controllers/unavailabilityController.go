package barberBookingController

import (
	// "context"
	"time"
	"errors"
	"gorm.io/gorm"

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
		if errors.Is(err,errors.New("unavailability already exists for this date")) {
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



