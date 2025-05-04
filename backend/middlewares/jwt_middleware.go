package middlewares

import (
	"os"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
)

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		tokenStr := c.Cookies("token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing or invalid token",
			})
		}

		secret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// ดึง user_id
		userIDFloat, ok := claims["user_id"].(float64)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid user_id",
			})
		}
		userID := uint(userIDFloat)

		// ดึง role
		role, ok := claims["role"].(string)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid role",
			})
		}

		// ดึง tenant_id (optional สำหรับ SUPER_ADMIN)
		var tenantID uint
		if tid, ok := claims["tenant_id"].(float64); ok {
			tenantID = uint(tid)
			c.Locals("tenant_id", tenantID)
		}

		c.Locals("user_id", userID)
		c.Locals("role", role)

		return c.Next()
	}
}
