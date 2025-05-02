package coreServiceTest

import (
    "context"
    "os"
    "testing"

    "github.com/golang-jwt/jwt/v4"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "golang.org/x/crypto/bcrypt"

    Core_authDto "myapp/modules/core/dto/auth"
    coreModels "myapp/modules/core/models"
    coreTests "myapp/modules/core/tests"
	"myapp/modules/core/services"
)

// fakeLogSvc สแต็บ SystemLogService
type fakeLogSvc struct{}

func (f *fakeLogSvc) Create(ctx context.Context, entry *coreModels.SystemLog) error {
    return nil
}
func (f *fakeLogSvc) Query(ctx context.Context, filter services.LogFilter) ([]coreModels.SystemLog, int64, error) {
    return nil, 0, nil
}
func (f *fakeLogSvc) GetByID(ctx context.Context, id uint) (coreModels.SystemLog, error) {
    return coreModels.SystemLog{}, nil
}

// ช่วยสร้าง AuthService พร้อม DB in‑memory และ seed data
func setupAuthService(t *testing.T) *services.AuthService {
    // เชื่อม in‑memory DB แล้ว migrate Role + User
    db := coreTests.SetupTestDB()
    err := db.AutoMigrate(&coreModels.Role{}, &coreModels.User{})
    require.NoError(t, err)

    // seed Role
    r := coreModels.Role{Name: coreModels.RoleNameUser}
    require.NoError(t, db.Create(&r).Error)

    // seed User (hash password)
    rawPw := "secret123"
    hash, err := bcrypt.GenerateFromPassword([]byte(rawPw), bcrypt.DefaultCost)
    require.NoError(t, err)
    u := coreModels.User{
        Username: "alice",
        Email:    "alice@example.com",
        Password: string(hash),
        RoleID:   r.ID,
    }
    require.NoError(t, db.Create(&u).Error)

    // ตั้ง JWT_SECRET
    os.Setenv("JWT_SECRET", "testsecret")

    return services.NewAuthService(db, &fakeLogSvc{})
}

func TestLogin_Success(t *testing.T) {
    svc := setupAuthService(t)

    req := Core_authDto.LoginRequest{
        Email:    "alice@example.com",
        Password: "secret123",
    }
    resp, err := svc.Login(context.Background(), req)
    require.NoError(t, err)
    require.NotEmpty(t, resp.Token)
    assert.Equal(t, uint(1), resp.User.ID)
    assert.Equal(t, "alice", resp.User.Username)
    assert.Equal(t, "alice@example.com", resp.User.Email)
    assert.Equal(t, uint(1), resp.User.RoleID)
    assert.Equal(t, string(coreModels.RoleNameUser), resp.User.Role)

    // ยืนยันว่า token มี claim ถูกต้อง
    token, err := jwt.Parse(resp.Token, func(tok *jwt.Token) (interface{}, error) {
        return []byte(os.Getenv("JWT_SECRET")), nil
    })
    require.NoError(t, err)
    claims := token.Claims.(jwt.MapClaims)
    assert.Equal(t, float64(1), claims["user_id"])
    assert.Equal(t, float64(1), claims["role_id"])
    assert.Equal(t, string(coreModels.RoleNameUser), claims["role"])
}

func TestLogin_UserNotFound(t *testing.T) {
    svc := setupAuthService(t)

    req := Core_authDto.LoginRequest{
        Email:    "noone@example.com",
        Password: "whatever",
    }
    _, err := svc.Login(context.Background(), req)
    require.EqualError(t, err, "invalid credentials")
}

func TestLogin_WrongPassword(t *testing.T) {
    svc := setupAuthService(t)

    req := Core_authDto.LoginRequest{
        Email:    "alice@example.com",
        Password: "wrongpass",
    }
    _, err := svc.Login(context.Background(), req)
    require.EqualError(t, err, "invalid credentials")
}
