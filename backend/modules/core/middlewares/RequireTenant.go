package middlewares

import (
	"github.com/gofiber/fiber/v2"
	"strconv"
)

func RequireTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantIDParam := c.Params("tenant_id")
		if tenantIDParam == "" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Tenant ID is required in the path",
			})
		}

		tenantID, err := strconv.ParseUint(tenantIDParam, 10, 64)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid tenant ID format",
			})
		}

		c.Locals("tenant_id", uint(tenantID))
		return c.Next()
	}
}

