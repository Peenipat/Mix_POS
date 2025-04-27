package controllers

import (
	authDto "myapp/dto/auth"
	userDto "myapp/dto/user"
	"myapp/services"
	"strconv"

	"github.com/gofiber/fiber/v2"
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

// ChangeUserRole godoc
// @Summary      เปลี่ยน Role ของผู้ใช้
// @Description  สำหรับ Super Admin เพื่อเปลี่ยน Role ของผู้ใช้
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        body body userDto.ChangeRoleInput true "ข้อมูลผู้ใช้งาน"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /admin/change_role [put]
// @Security     ApiKeyAuth
func ChangeUserRole(c *fiber.Ctx) error {
	var input userDto.ChangeRoleInput

	// รับข้อมูลจาก body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	// เรียกใช้ Service
	if err := services.ChangeRoleFromAdmin(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "user role changed successfully",
	})
}

// GetAllUsers godoc
// @Summary ดึงข้อมูลผู้ใช้งานทั้งหมด
// @Description สำหรับ Super Admin ดึง Users ทั้งหมด พร้อม Pagination
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "หน้าที่ต้องการ (default 1)"
// @Param limit query int false "จำนวนรายการต่อหน้า (default 10)"
// @Success 200 {array} userDto.UserResponse
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
// @Security ApiKeyAuth
func GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	offset := (page - 1) * limit

	users, err := services.GetAllUsers(limit, offset)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": "failed to fetch users"})
	}

	return c.Status(200).JSON(users)
}

// FilterUsersByRole godoc
// @Summary      ดึงข้อมูลผู้ใช้งานโดยเลือกเฉพาะ role ที่ต้องการ
// @Description  ใช้สำหรับ Super Admin เพื่อดึง Users เฉพาะ role ที่ระบุ เช่น STAFF, USER, BRANCH_ADMIN
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        role query string true "Role ที่ต้องการกรอง เช่น STAFF, USER, BRANCH_ADMIN, SUPER_ADMIN"
// @Success      200 {array} userDto.UserResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/user-by-role [get]
// @Security     ApiKeyAuth
func FilterUsersByRole(c *fiber.Ctx) error {
	// ดึง role ของคนที่ login (จาก JWT Middleware ที่ c.Locals())
	userRole := c.Locals("userRole")
	if userRole != "SUPER_ADMIN" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "permission denied"})
	}

	// อ่าน query string ?role=STAFF
	role := c.Query("role")
	if role == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "missing role parameter"})
	}

	// เรียก service ไปหาข้อมูล
	users, err := services.FilterUsersByRole(role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(users)
}
