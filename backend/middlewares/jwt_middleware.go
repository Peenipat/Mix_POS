package middlewares

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/gofiber/fiber/v2"
)

func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		// ดึง token จาก cookie
		tokenStr := c.Cookies("token")
		if tokenStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Missing token",
			})
		}

		secret := os.Getenv("JWT_SECRET")

		// แกะ token
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		// อ่าน claims
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid token claims",
			})
		}

		// เช็ค expiry แบบปลอดภัย
		if expRaw, ok := claims["exp"].(float64); !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token missing expiry",
			})
		} else if int64(expRaw) < time.Now().Unix() {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Token has expired",
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
		roleStr, ok := claims["role"].(string)
		if !ok || roleStr == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid role",
			})
		}

		// ดึง tenant_id (optional)
		if tid, ok := claims["tenant_id"].(float64); ok {
			c.Locals("tenant_id", uint(tid))
		}

		// Set ลงใน context
		c.Locals("user_id", userID)
		c.Locals("role", roleStr)

		return c.Next()
	}
}


