package controllers

import (
	"myapp/dto/auth"
	"myapp/services"

	"github.com/gofiber/fiber/v2"
)

// @Summary Login
// @Description Auth login เพื่อรับ JWT Token
// @Tags Auth
// @Accept json
// @Produce json
// @Param body body authDto.LoginRequest true "ข้อมูลเข้าสู่ระบบ"
// @Success 200 {object} authDto.LoginResponse
// @Failure 401 {object} map[string]string
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	var input authDto.LoginRequest
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// เรียก service ที่ return *AuthResponse
	response, err := services.Login(input)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// ตอบกลับด้วย token + user info
	return c.JSON(response)
}
