package middlewares
import (
		"github.com/gofiber/fiber/v2"
		"myapp/models/core"

)
func RequireSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole := c.Locals("role")
		roleStr, ok := userRole.(string)
		if !ok || roleStr != string(coreModels.RoleNameSaaSSuperAdmin) {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "permission denied",
			})
		}
		return c.Next()
	}
}


