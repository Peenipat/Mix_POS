// services/auth_service.go
package services

import (
    "context"
    "errors"
    "os"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

    authDto "myapp/dto/auth"
    "myapp/models"
)

type AuthService struct {
    db     *gorm.DB
    logSvc SystemLogService
}

func NewAuthService(db *gorm.DB, logSvc SystemLogService) *AuthService {
    return &AuthService{db: db, logSvc: logSvc}
}

func (s *AuthService) Login(ctx context.Context, input authDto.LoginRequest) (*authDto.LoginResponse, error) {
    var user models.User
    err := s.db.WithContext(ctx).
        Where("email = ?", input.Email).
        First(&user).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid credentials")
        }
        return nil, err
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
        return nil, errors.New("invalid credentials")
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(72 * time.Hour).Unix(),
    })
    signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return nil, errors.New("could not generate token")
    }

    return &authDto.LoginResponse{
        Token: signed,
        User: authDto.UserInfoResponse{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
            Role:     string(user.Role),
        },
    }, nil
}
