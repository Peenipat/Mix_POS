package controllers

import (
	"github.com/gofiber/fiber/v2"
	"myapp/dto/auth"
	"myapp/dto/user"
	"myapp/services"
)

func CreateUserFromRegister(c *fiber.Ctx) error {
	var input authDto.RegisterInput 
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := services.CreateUserFromRegister(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

func CreateUserFromAdmin(c *fiber.Ctx) error {
	var input userDto.CreateUserInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid input"})
	}

	if err := services.CreateUserFromAdmin(input); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User created successfully"})
}

