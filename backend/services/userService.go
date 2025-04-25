package services

import (
	"errors"
	"myapp/database"
	"myapp/dto/auth"
	"myapp/dto/user"
	"myapp/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserFromRegister(input authDto.RegisterInput) error {
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
		Role:     models.RoleUser, // default เป็น User
	}

	// save user ลง database
	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}

func CreateUserFromAdmin(input userDto.CreateUserInput) error {
	// ตรวจ email ซ้ำ
	var existingUser models.User
	database.DB.Where("email = ?", input.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return errors.New("email already in use")
	}

	// ป้องกันการสร้าง SUPER_ADMIN โดยเด็ดขาด
	if input.Role == string(models.RoleSuperAdmin) {
		return errors.New("cannot create SUPER_ADMIN")
	}

	// ตรวจว่า role ที่ใส่มาเป็น role ที่ระบบอนุญาตให้ admin สร้างได้
	switch models.Role(input.Role) {
	case models.RoleBranchAdmin, models.RoleUser, models.RoleStaff:
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
		Role:     models.Role(input.Role),
	}

	// Save ลง database
	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}


