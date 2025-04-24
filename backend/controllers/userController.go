package controllers

import (
	"myapp/database"
	"myapp/models"
	
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

func CreateUser(c *fiber.Ctx) error {
	type RegisterInput struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var existingUser models.User
	database.DB.Where("email = ?", input.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Email already in use"})
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), 14)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Cannot hash password"})
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Create user failed"})
	}
	return c.JSON(fiber.Map{"message": "User registered successfully"})
}