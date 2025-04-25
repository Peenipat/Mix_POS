package services

import (
	"errors"
	"myapp/database"
	"myapp/dto/auth"
	"myapp/models"
	"golang.org/x/crypto/bcrypt"
)

func CreateUserFromRegister(input dto.RegisterInput) error {
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

