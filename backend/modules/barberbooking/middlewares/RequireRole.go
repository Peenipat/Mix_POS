package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func RequireTenant() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tenantID := c.Locals("tenant_id")
		if tenantID == nil {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "Permission denied: tenant access required",
			})
		}
		return c.Next()
	}
}
