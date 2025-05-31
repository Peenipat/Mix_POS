package Core_controllers

import (
	"context"
	"encoding/json"
	"time"

	Core_authDto "myapp/modules/core/dto/auth"	
	coreServices "myapp/modules/core/services"
    "myapp/modules/core/models"
	"github.com/gofiber/fiber/v2"
)

var (
    authSvc *coreServices.AuthService // ตัว logic login 
    logSvc  coreServices.SystemLogService //ตัวสำหรับ save log login
)
// init Dependency Injection
func InitAuthHandler(a *coreServices.AuthService, l coreServices.SystemLogService) {
    authSvc = a
    logSvc = l
}


// LoginHandler godoc
// @Summary      User login
// @Description  Authenticates a user with email & password, issues a JWT cookie, and logs the attempt.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        body  body      Core_authDto.LoginRequest  true  "Login credentials (email & password)"
// @Success      200   {object}  map[string]interface{}     "Returns the authenticated user object"
// @Failure      400   {object}  map[string]string          "Invalid input"
// @Failure      401   {object}  map[string]string          "Unauthorized – wrong credentials"
// @Router       /login [post]
// @Security     ApiKeyAuth
func LoginHandler(c *fiber.Ctx) error {
    // 1) Bind request
    var req Core_authDto.LoginRequest 
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
    }

    // 2) Call AuthService.Login
    resp, err := authSvc.Login(context.Background(), req)

    // 3) Prepare common log entry fields
    entry := &coreModels.SystemLog{
        CreatedAt:  time.Now(),
        HTTPMethod: c.Method(),
        Endpoint:   c.Path(),
        Resource:   "Auth",
        Action:     "LOGIN",
    }
    // เก็บค่า Ip เป็น string
    ip := c.IP() //ดึง ip จาก request
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

    // สร้าง cookie 
    c.Cookie(&fiber.Cookie{
        Name:     "token",
        Value:    resp.Token, // <-- ใช้ resp.Token (ที่ service login สร้างไว้แล้ว)
        Expires:  time.Now().Add(72 * time.Hour),
        HTTPOnly: true, // อ่าน cookies จาก client
        Secure:   false,    // ต้องใช้ https ตอน production
        SameSite: "Lax",   
    })

    // 4) Return response
    return c.JSON(fiber.Map{
        "user": resp.User,
    })
}




