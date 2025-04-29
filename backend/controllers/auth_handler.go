package controllers

import (
    "context"
    "encoding/json"
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

func InitAuthHandler(a *services.AuthService, l services.SystemLogService) {
    authSvc = a
    logSvc = l
}

func LoginHandler(c *fiber.Ctx) error {
    // 1) Bind request
    var req authDto.LoginRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // 2) Call AuthService.Login
    resp, err := authSvc.Login(context.Background(), req)

    // 3) Prepare common log entry fields
    entry := &models.SystemLog{
        CreatedAt:  time.Now(),
        HTTPMethod: c.Method(),
        Endpoint:   c.Path(),
        Resource:   "Auth",
        Action:     "LOGIN",
    }
    // Store IP as string
    ip := c.IP()
    entry.IPAddress = &ip

    if err != nil {
        // LOGIN_FAILURE
        entry.Status = "failure"
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

    // 4) Return response
    return c.JSON(resp)
}
