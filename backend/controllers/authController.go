package controllers

import (
	"myapp/database"
	"myapp/models"
	"os"
	"time"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)



func Login(c *fiber.Ctx) error {
	type LoingInput struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	var input LoingInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)
	if user.ID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid credentials"})
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("JWT_SECRET is not set")
	}

	t, err := token.SignedString([]byte(secret))
	if err != nil{
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error":"Could not login"})
	}
	return c.JSON(fiber.Map{"token":t})
}
