package services

import (
	"errors"
	"os"
	"time"

	"myapp/database"
	"myapp/dto/auth"
	"myapp/models"
	

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
	
)

func Login(input dto.LoginRequest) (*dto.AuthResponse, error) {
	var user models.User
	database.DB.Where("email = ?", input.Email).First(&user)

	if user.ID == 0 {
		return nil, errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return nil, errors.New("invalid credentials")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(),
	})

	secret := os.Getenv("JWT_SECRET")
	t, err := token.SignedString([]byte(secret))
	if err != nil {
		return nil, errors.New("could not generate token")
	}

	return &dto.AuthResponse{
		Token: t,
		User: struct {
			ID       uint   `json:"id"`
			Username string `json:"username"`
			Email    string `json:"email"`
			Role     string `json:"role"`
		}{
			ID:       user.ID,
			Username: user.Username,
			Email:    user.Email,
			Role:     string(user.Role),
		},
	}, nil
	
}

