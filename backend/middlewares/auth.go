package middlewares
import (
		"github.com/gofiber/fiber/v2"
		"myapp/models"
)
// check สิทธิ SUPER_ADMIN
func RequireSuperAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง role จาก local
		userRole := c.Locals("role")
		if userRole != models.RoleSuperAdmin{
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"error": "permission denied",
			})
		}
		// ผ่านการตรวจสอบ ไป middleware ตัวถัดไป
		return c.Next()
	}
}
