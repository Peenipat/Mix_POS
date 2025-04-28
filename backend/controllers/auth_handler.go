package controllers

import (
	"encoding/json"
	"net"
	"time"

	"github.com/gofiber/fiber/v2"
	authDto "myapp/dto/auth"
	"myapp/models"
	"myapp/services"
)

var (
	authSvc *services.AuthService
	logSvc  services.SystemLogService
)

// InitAuthHandler injects AuthService and SystemLogService into controller
func InitAuthHandler(a *services.AuthService, l services.SystemLogService) {
	authSvc = a
	logSvc = l
}

// LoginHandler handles POST /auth/login, logs success/failure
func LoginHandler(c *fiber.Ctx) error {
	var req authDto.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	resp, err := authSvc.Login(c.Context(), req)

	// Prepare common log fields
	entry := &models.SystemLog{
		CreatedAt:  time.Now(),
		HTTPMethod: c.Method(),
		Endpoint:   c.Path(),
		Resource:   "User",
		Action:     "LOGIN",
	}

	// Parse and assign IP address
	if ip := net.ParseIP(c.IP()); ip != nil {
		entry.IPAddress = &ip
	}

	if err != nil {
		// LOGIN_FAILURE
		entry.Status = "failure"
		// Include attempted email in details
		if b, jerr := json.Marshal(map[string]string{"email": req.Email}); jerr == nil {
			entry.Details = b
		}
		logSvc.Create(c.Context(), entry)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	// LOGIN_SUCCESS
	entry.Status = "success"
	entry.UserID = &resp.User.ID
	role := resp.User.Role
	entry.UserRole = &role
	logSvc.Create(c.Context(), entry)

	return c.JSON(resp)
}
