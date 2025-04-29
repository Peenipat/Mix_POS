package middlewares
import (
		"github.com/gofiber/fiber/v2"
		"myapp/models"

)
func RequireSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		roleStr, ok := userRole.(string)
		if !ok || roleStr != string(models.RoleSuperAdmin) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "permission denied",
			})
		}
		return c.Next()
	}
}


