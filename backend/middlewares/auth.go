package middlewares
import (
		"github.com/gofiber/fiber/v2"
)
func RequireSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		if userRole != "SUPER_ADMIN" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "permission denied",
			})
		}
		return c.Next()
	}
}
