package controllers

import (
	"github.com/gofiber/fiber/v2"
	"myapp/dto/auth"
	"myapp/dto/user"
	"myapp/services"
)

// @Summary        สร้าง Account Role USER
// @Description    ลงทะเบียนเพื่อ สร้าง Account โดย User เป็นคนสร้างเอง
// @Tags           Auth
// @Accept         json
// @Produce        json
// @Param          body body authDto.RegisterInput true "ข้อมูลผู้ใช้"
// @Success        200 {object} map[string]string "ลงทะเบียนสำเร็จ"
// @Failure        400 {object} map[string]string "ข้อมูลไม่ถูกต้องหรือลงทะเบียนล้มเหลว"
// @Router         /auth/register [post]
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

// CreateUserFromAdmin godoc
// @Summary      สร้างผู้ใช้โดย Super Admin
// @Description  ใช้สำหรับ SUPER_ADMIN สร้าง User role อื่น ๆ แต่ไม่สามารถใช้สร้าง SUPER_ADMIN ได้
// @Tags         User
// @Accept       json
// @Produce      json
// @Param        body  body  userDto.CreateUserInput  true  "ข้อมูลผู้ใช้งาน"
// @Success      200  {object}  models.User
// @Failure 400 {object} map[string]string
// @Router       /admin/create_users [post]
// @Security     ApiKeyAuth
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

