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

    Core_authDto "myapp/modules/core/dto/auth"
    "myapp/modules/core/models" 
)

// ตัว strct มีไว้เก็บ type ที่ต้องการสำหรับการทำ Dependency Injection
type AuthService struct {
    db     *gorm.DB //ตัวจัดการ database
    logSvc SystemLogService //ตัวแปรสำหรับเก็บ log
}
// function เริ่มสร้าง AuthService โดยการรับ Database และ service เข้ามา
func NewAuthService(db *gorm.DB, logSvc SystemLogService) *AuthService {
    return &AuthService{db: db, logSvc: logSvc}
}

// Login ตรวจสอบข้อมูลล็อกอินและสร้าง JWT
// คืนค่า DTO ที่ประกอบด้วย token และข้อมูล user หรือ error
// services/authService.go
func (s *AuthService) Login(ctx context.Context, input Core_authDto.LoginRequest) (*Core_authDto.LoginResponse, error) {
    // 1. ดึง user และ preload Role
    var user coreModels.User
    err := s.db.WithContext(ctx).
        Preload("Role").                             // ← โหลด Role struct มาให้
        Where("email = ?", input.Email).
        First(&user).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, errors.New("invalid credentials")
        }
        return nil, err
    }

    // 2. ตรวจรหัสผ่าน
    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
        return nil, errors.New("invalid credentials")
    }

    // 3. สร้าง JWT claims
    claims := jwt.MapClaims{
        "user_id": user.ID,
        "role_id": user.RoleID,
        "role":    user.Role.Name,
        "exp":     time.Now().Add(72 * time.Hour).Unix(),
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    secret := os.Getenv("JWT_SECRET")
    signed, err := token.SignedString([]byte(secret))
    if err != nil {
        return nil, errors.New("could not generate token")
    }

    // 4. คืนค่า LoginResponse
    return &Core_authDto.LoginResponse{
        Token: signed,
        User: Core_authDto.UserInfoResponse{
            ID:       user.ID,
            Username: user.Username,
            Email:    user.Email,
            RoleID:   user.RoleID,
            Role:     string(user.Role.Name),
        },
    }, nil
}

