package controllers

import (
	"myapp/dto/auth"
	"myapp/services"

	"github.com/gofiber/fiber/v2"
)

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
