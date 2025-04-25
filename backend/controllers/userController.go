package controllers

import (
	"github.com/gofiber/fiber/v2"
	"myapp/dto/auth"
	"myapp/services"
)

func Resgister(c *fiber.Ctx) error {
	var input dto.RegisterInput 
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.CreateUserFromRegister(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}
