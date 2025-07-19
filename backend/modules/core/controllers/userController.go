package Core_controllers

import (
	"errors"
	"log"
	aws "myapp/cmd/worker"
	corePort "myapp/modules/core/port"
	coreServices "myapp/modules/core/services"
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type UserController struct {
	UserService corePort.IUser
}

func NewUserController(scv corePort.IUser) *UserController {
	return &UserController{
		UserService: scv,
	}
}

// @Summary        สร้าง Account Role USER
// @Description    ลงทะเบียนเพื่อ สร้าง Account โดย User เป็นคนสร้างเอง
// @Tags           Auth
// @Accept         json
// @Produce        json
// @Param          body body corePort.RegisterInput true "ข้อมูลผู้ใช้"
// @Success        200 {object} map[string]string "ลงทะเบียนสำเร็จ"
// @Failure        400 {object} map[string]string "ข้อมูลไม่ถูกต้องหรือลงทะเบียนล้มเหลว"
// @Router         /auth/register [post]
func (ctrl *UserController) CreateUserFromRegister(c *fiber.Ctx) error {
	var input corePort.RegisterInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	err := ctrl.UserService.CreateUserFromRegister(input)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(fiber.Map{"message": "User registered successfully"})
}

// @Summary      สร้างผู้ใช้โดย Super Admin
// @Description  ใช้สำหรับ SUPER_ADMIN สร้าง User role อื่น ๆ แต่ไม่สามารถใช้สร้าง SUPER_ADMIN ได้
// @Tags         User
// @Accept       multipart/form-data
// @Produce      json
// @Param        username   formData  string  true  "ชื่อผู้ใช้"
// @Param        email      formData  string  true  "อีเมล"
// @Param        password   formData  string  true  "รหัสผ่าน"
// @Param        role       formData  string  true  "บทบาทของผู้ใช้ (เช่น ADMIN, USER)"
// @Param        branch_id  formData  uint    false "รหัสสาขา (ถ้ามี)"
// @Param        avatar     formData  file    false "รูปภาพโปรไฟล์"
// @Success      201        {object}  map[string]interface{}
// @Failure      400        {object}  map[string]string
// @Failure      500        {object}  map[string]string
// @Router       /admin/create_users [post]
// @Security     ApiKeyAuth
func (ctrl *UserController) CreateUserFromAdmin(c *fiber.Ctx) error {
	input := new(corePort.CreateUserInput)

	input.Username = c.FormValue("username")
	input.Email = c.FormValue("email")
	input.Password = c.FormValue("password")
	input.Role = c.FormValue("role")

	// แปลงค่า branch_id จาก string → uint
	if b := c.FormValue("branch_id"); b != "" {
		if bid, err := strconv.ParseUint(b, 10, 64); err == nil {
			u := uint(bid)
			input.BranchID = &u
		}
	}

	keyPrefix := c.FormValue("keyprefix")
	if keyPrefix == "" {
		keyPrefix = "avatars/"
	}
    log.Printf("keyPrefix : %s",keyPrefix)
	// ตรวจสอบและอัปโหลด avatar (ถ้ามี)
	fileHeader, err := c.FormFile("file")
	if err == nil && fileHeader != nil {
		url, filename, err := aws.UploadToS3(fileHeader,keyPrefix)
		if err != nil {
			log.Printf("UploadToS3 error: %v", err)
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "failed to upload avatar"})
		}
		input.Img_path = url
		input.Img_name = filename
	}

	// เรียกใช้ service เพื่อสร้าง user
	if err := ctrl.UserService.CreateUserFromAdmin(*input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message":  "User created successfully",
		"img_path": input.Img_path,
		"img_name": input.Img_name,
	})
}

// ChangeUserRole godoc
// @Summary      เปลี่ยน Role ของผู้ใช้
// @Description  สำหรับ Super Admin เพื่อเปลี่ยน Role ของผู้ใช้
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        body body corePort.ChangeRoleInput true "ข้อมูลผู้ใช้งาน"
// @Success      200 {object} map[string]string
// @Failure      400 {object} map[string]string
// @Router       /admin/change_role [put]
// @Security     ApiKeyAuth
func (ctrl *UserController) ChangeUserRole(c *fiber.Ctx) error {
	var input corePort.ChangeRoleInput

	// รับข้อมูลจาก body
	if err := c.BodyParser(&input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "invalid input",
		})
	}

	// เรียกใช้ Service
	if err := ctrl.UserService.ChangeRoleFromAdmin(input); err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(200).JSON(fiber.Map{
		"message": "user role changed successfully",
	})
}

// var getAllUsersFunc = coreServices.GetAllUsers
// func InitGetAllUsers(fn func(limit, offset int) ([]corePort.UserInfoResponse, error)) {
//   getAllUsersFunc = fn
// }
// GetAllUsers godoc
// @Summary ดึงข้อมูลผู้ใช้งานทั้งหมด
// @Description สำหรับ Super Admin ดึง Users ทั้งหมด พร้อม Pagination
// @Tags Admin
// @Accept json
// @Produce json
// @Param page query int false "หน้าที่ต้องการ (default 1)"
// @Param limit query int false "จำนวนรายการต่อหน้า (default 10)"
// @Success 200 {array} corePort.UserInfoResponse
// @Failure 500 {object} map[string]string
// @Router /admin/users [get]
// @Security ApiKeyAuth
func (ctrl *UserController) GetAllUsers(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	offset := (page - 1) * limit

	users, err := ctrl.UserService.GetAllUsers(limit, offset)
	if err != nil {
	  return c.Status(500).JSON(fiber.Map{"error": "failed to fetch users"})
	}
	return c.Status(200).JSON(users)
  }

//   var filterUsersByRoleFunc = coreServices.FilterUsersByRole

// func InitFilterUsersByRole(fn func(string) ([]corePort.UserInfoResponse, error)) {
//     filterUsersByRoleFunc = fn
// }
// FilterUsersByRole godoc
// @Summary      ดึงข้อมูลผู้ใช้งานโดยเลือกเฉพาะ role ที่ต้องการ
// @Description  ใช้สำหรับ Super Admin เพื่อดึง Users เฉพาะ role ที่ระบุ เช่น STAFF, USER, BRANCH_ADMIN
// @Tags         Admin
// @Accept       json
// @Produce      json
// @Param        role query string true "Role ที่ต้องการกรอง เช่น STAFF, USER, BRANCH_ADMIN, SUPER_ADMIN"
// @Success      200 {array} corePort.UserInfoResponse
// @Failure      400 {object} map[string]string
// @Failure      403 {object} map[string]string
// @Failure      500 {object} map[string]string
// @Router       /admin/user-by-role [get]
// @Security     ApiKeyAuth
func (ctrl *UserController) FilterUsersByRole(c *fiber.Ctx) error {
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
	users, err := ctrl.UserService.FilterUsersByRole(role)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(users)
}

type ChangePasswordRequest struct {
	OldPassword string `json:"old_password"`
	NewPassword string `json:"new_password"`
}

// ChangePassword handles PUT /users/:id/password
func (ctrl *UserController) ChangePassword(c *fiber.Ctx) error {
	// 1. Authorization (optional): check role / ownership
	// e.g. userIDFromToken := c.Locals("user_id").(uint)
	// if userIDFromToken != targetID && !IsAdmin(...) { return 403 }

	// 2. Parse user ID from path
	idParam := c.Params("id")
	id64, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil || id64 == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user ID",
		})
	}
	userID := uint(id64)

	// 3. Parse JSON body
	var body ChangePasswordRequest
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Malformed JSON",
		})
	}

	// 4. Call service
	err = ctrl.UserService.ChangePassword(c.Context(), userID, body.OldPassword, body.NewPassword)
	if err != nil {
		switch {
		case errors.Is(err, coreServices.ErrUserNotFound):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": "User not found",
			})
		case errors.Is(err, coreServices.ErrInvalidOldPassword):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Old password is incorrect",
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to change password",
			})
		}
	}

	// 5. Success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Password changed successfully",
	})
}

func (ac *UserController) Me(c *fiber.Ctx) error {
	// 1. ดึง user_id จาก middleware RequireAuth
	uidVal := c.Locals("user_id")
	if uidVal == nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"status":  "error",
			"message": "User not authenticated",
		})
	}
	userID, ok := uidVal.(uint)
	if !ok {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid user_id type in context",
		})
	}

	// 2. เรียก service Me
	meDTO, err := ac.UserService.Me(c.Context(), userID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch user info",
			"error":   err.Error(),
		})
	}
	if meDTO == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "User not found",
		})
	}

	// 3. ตอบกลับ
	return c.Status(http.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "User profile retrieved",
		"data":    meDTO,
	})
}
