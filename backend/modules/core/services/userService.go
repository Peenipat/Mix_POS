package services

import (
	"errors"
	"myapp/database"
	Core_authDto "myapp/modules/core/dto/auth"
	Core_userDto "myapp/modules/core/dto/user"
	"myapp/models"
	"myapp/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserFromRegister(input Core_authDto.RegisterInput) error {
	// ตรวจซ้ำ email
	var existingUser models.User
	database.DB.Where("email = ?", input.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return errors.New("email already in use")
	}

	// hash password ด้วย bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     models.RoleNameUser, // default เป็น User
	}

	// save user ลง database
	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
func CreateUserFromAdmin(input Core_userDto.CreateUserInput) error {
	
	// ตรวจ email ซ้ำ
	var existingUser models.User
	database.DB.Where("email = ?", input.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return errors.New("email already in use")
	}

	// ป้องกันการสร้าง SUPER_ADMIN โดยเด็ดขาด
	if input.Role == string(models.RoleNameSaaSSuperAdmin) {
		return errors.New("cannot create SUPER_ADMIN")
	}

	// ตรวจว่า role ที่ใส่มาเป็น role ที่ระบบอนุญาตให้ admin สร้างได้
	switch models.RoleName(input.Role) {
	case 
	models.RoleNameTenantAdmin,
	models.RoleNameBranchAdmin, 
	models.RoleNameAssistantManager, 
	models.RoleNameStaff,
	models.RoleNameUser:

	default:
		return errors.New("invalid role provided")
	}

	// hash password อย่างปลอดภัย
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	// เตรียมสร้าง user
	user := models.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		Role:     models.RoleName(input.Role),
	}
	// Save ลง database
	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func ChangeRoleFromAdmin(input Core_userDto.ChangeRoleInput) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var user models.User 

		// หา user
		if err := tx.First(&user, input.ID).Error; err != nil {
			return errors.New("user not found")
		}

		// ไม่อนุญาตเปลี่ยนเป็น SUPER_ADMIN
		if input.Role == string(models.RoleNameSaaSSuperAdmin) {
			return errors.New("cannot change role to SUPER_ADMIN")
		}

		// Validate role ใหม่
		validRoles := []models.RoleName{
			models.RoleNameTenantAdmin,
			models.RoleNameBranchAdmin, 
			models.RoleNameAssistantManager, 
			models.RoleNameStaff,
			models.RoleNameUser,
		}
		isValid := false
		for _, r := range validRoles {
			if models.RoleName(input.Role) == r {
				isValid = true
				break
			}
		}
		if !isValid {
			return errors.New("invalid role provided")
		}

		// อัปเดต role
		user.Role = models.RoleName(input.Role)
		if err := tx.Save(&user).Error; err != nil {
			return errors.New("failed to update user role")
		}

		return nil
	})
}

func GetAllUsers(limit int , offset int)([]Core_userDto.UserResponse, error){
	var users []models.User
	//ค้นหา user เช็ค limit และกำหนด offset 
	if err := database.DB.Order("id ASC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil,err
	}

	result := utils.MapSlice(users, func(u models.User) Core_userDto.UserResponse {
		return Core_userDto.UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
		}
	})

	return result, nil
}

func FilterUsersByRole(role string) ([]Core_userDto.UserResponse, error) {
	// Validate role ก่อน
	validRoles := []models.RoleName{
		models.RoleNameTenantAdmin,
		models.RoleNameBranchAdmin,
		models.RoleNameAssistantManager,
		models.RoleNameStaff,
		models.RoleNameUser, 
	}

	isValid := false
	for _, r := range validRoles {
		if models.RoleName(role) == r {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid role")
	}

	// ดึง users
	var users []models.User
	if err := database.DB.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users")
	}

	// Map เป็น UserResponse
	var result []Core_userDto.UserResponse
	for _, u := range users {
		result = append(result, Core_userDto.UserResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			Role:     string(u.Role),
		})
	}

	return result, nil
}



