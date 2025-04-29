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
func (s *AuthService) Login(ctx context.Context, input authDto.LoginRequest) (*authDto.LoginResponse, error) {
    //ดึงข้อมูล user จาก Database
    var user models.User
    err := s.db.WithContext(ctx).
        Where("email = ?", input.Email).
        First(&user).Error

    if err != nil { 
        if errors.Is(err, gorm.ErrRecordNotFound) {
            //กรณี ไม่เจอ user
            return nil, errors.New("invalid credentials")
        }
        //กรณ๊อื่น ๆ 
        return nil, err
    }

    if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)) != nil {
        return nil, errors.New("invalid credentials")
    }

    //สร้าง JWT claims 
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
        "user_id": user.ID,
        "role":    user.Role,
        "exp":     time.Now().Add(72 * time.Hour).Unix(),
    })

    // ลงลายเซ็นจาก JWT_SECRET ที่อยู่ใน env
    signed, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
    if err != nil {
        return nil, errors.New("could not generate token")
    }

    //ส่งคืนของที่จำเป็นตามที่กำหนดไว้ใน LoginRespone
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
