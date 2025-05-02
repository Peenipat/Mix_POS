package coreControllersTest
import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"
	"errors"

    "github.com/gofiber/fiber/v2"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"

	 "myapp/modules/core/controllers"
    authDto "myapp/modules/core/dto/auth"
    coreModels "myapp/modules/core/models"
    coreTests "myapp/modules/core/tests"
    "myapp/database"
)

func setupApp(t *testing.T) (*fiber.App, *gorm.DB) {
    // 1) เปิด in-memory DB แล้ว migrate Role + User
    db := coreTests.SetupTestDB()
    require.NoError(t, db.AutoMigrate(&coreModels.Role{}, &coreModels.User{}))

    // 2) seed default Role.USER
    role := coreModels.Role{Name: coreModels.RoleNameUser}
    require.NoError(t, db.Create(&role).Error)

    // 3) override global DB ของ application ให้เป็นตัวนี้
    database.DB = db

    // 4) สร้าง Fiber app และผูก route
    app := fiber.New()
    app.Post("/register", Core_controllers.CreateUserFromRegister)

    return app, db
}

func setupAppGetuser() *fiber.App {
	app := fiber.New()
	app.Get("/admin/users", Core_controllers.GetAllUsers)
	return app
  }

  func setupAppFilterRole() *fiber.App {
    app := fiber.New()
    app.Get("/admin/users-by-role", func(c *fiber.Ctx) error {
        // จำลอง JWT middleware
        c.Locals("userRole", c.Query("as")) 
        return Core_controllers.FilterUsersByRole(c)
    })
    return app
}


func mustHash(t *testing.T, raw string) []byte {
    h, err := bcrypt.GenerateFromPassword([]byte(raw), bcrypt.DefaultCost)
    require.NoError(t, err)
    return h
}

func TestCreateUserFromRegister_InvalidPayload(t *testing.T) {
    app, _ := setupApp(t)

    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBufferString("not-json"))
    req.Header.Set("Content-Type", "application/json")
    resp, err := app.Test(req, -1)
    require.NoError(t, err)

    assert.Equal(t, 400, resp.StatusCode)
    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Equal(t, "Invalid input", body["error"])
}

func TestCreateUserFromRegister_EmailAlreadyUsed(t *testing.T) {
    app, db := setupApp(t)

    // สร้าง user ซ้ำใน *same* DB ที่ app ใช้
    existing := coreModels.User{
        Username: "foo",
        Email:    "foo@example.com",
        Password: string(mustHash(t, "password")),
        RoleID:   1, // Role.USER.ID
    }
    require.NoError(t, db.Create(&existing).Error)

    payload := authDto.RegisterInput{
        Username: "bar",
        Email:    existing.Email,
        Password: "newpass",
    }
    buf, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(buf))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req, -1)
    require.NoError(t, err)
    assert.Equal(t, 400, resp.StatusCode)

    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Equal(t, "email already in use", body["error"])
}

func TestCreateUserFromRegister_Success(t *testing.T) {
    app, db := setupApp(t)

    payload := authDto.RegisterInput{
        Username: "newuser",
        Email:    "new@example.com",
        Password: "strongpass",
    }
    buf, _ := json.Marshal(payload)
    req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(buf))
    req.Header.Set("Content-Type", "application/json")

    resp, err := app.Test(req, -1)
    require.NoError(t, err)
    assert.Equal(t, 200, resp.StatusCode)

    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Equal(t, "User registered successfully", body["message"])

    // ตรวจใน *same* DB ว่ามี user เกิดใหม่
    var u coreModels.User
    err = db.First(&u, "email = ?", payload.Email).Error
    require.NoError(t, err)
    assert.Equal(t, payload.Username, u.Username)
    assert.NotEqual(t, payload.Password, u.Password) // ต้องถูก hash
}

func TestGetAllUsers_DefaultPaging(t *testing.T) {
	// stub ให้คืน slice เดียวกันที่เราต้องการ
	Core_controllers.InitGetAllUsers(func(limit, offset int) ([]authDto.UserInfoResponse, error) {
	  assert.Equal(t, 10, limit)     // default
	  assert.Equal(t, 0, offset)    // (1-1)*10
	  return []authDto.UserInfoResponse{{ID:1,Username:"u"}}, nil
	})
  
	app := setupAppGetuser()
	req := httptest.NewRequest("GET", "/admin/users", nil)
	resp, _ := app.Test(req, -1)
  
	assert.Equal(t, 200, resp.StatusCode)
	var arr []authDto.UserInfoResponse
	json.NewDecoder(resp.Body).Decode(&arr)
	assert.Len(t, arr, 1)
  }
  
  func TestGetAllUsers_CustomPaging(t *testing.T) {
	Core_controllers.InitGetAllUsers(func(limit, offset int) ([]authDto.UserInfoResponse, error) {
	  assert.Equal(t, 5, limit)
	  assert.Equal(t, 10, offset)    // (3-1)*5
	  return []authDto.UserInfoResponse{}, nil
	})
  
	app := setupAppGetuser()
	req := httptest.NewRequest("GET", "/admin/users?page=3&limit=5", nil)
	resp, _ := app.Test(req, -1)
  
	assert.Equal(t, 200, resp.StatusCode)
  }
  
  func TestGetAllUsers_ServiceError(t *testing.T) {
	Core_controllers.InitGetAllUsers(func(limit, offset int) ([]authDto.UserInfoResponse, error) {
	  return nil, errors.New("boom")
	})
  
	app := setupAppGetuser()
	resp, _ := app.Test(httptest.NewRequest("GET", "/admin/users", nil), -1)
	assert.Equal(t, 500, resp.StatusCode)
	var body map[string]string
	json.NewDecoder(resp.Body).Decode(&body)
	assert.Equal(t, "failed to fetch users", body["error"])
  }

  func TestFilterUsersByRole_NotAdmin(t *testing.T) {
    app := setupAppFilterRole()
    // ไม่ใช่ SUPER_ADMIN
    req := httptest.NewRequest("GET", "/admin/users-by-role?as=STAFF&role=USER", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 403, resp.StatusCode)
}

func TestFilterUsersByRole_MissingParam(t *testing.T) {
    app := setupAppFilterRole()
    req := httptest.NewRequest("GET", "/admin/users-by-role?as=SUPER_ADMIN", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 400, resp.StatusCode)

    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Equal(t, "missing role parameter", body["error"])
}

func TestFilterUsersByRole_ServiceError(t *testing.T) {
    app := setupAppFilterRole()
    // stub service ให้ error
    Core_controllers.InitFilterUsersByRole(func(role string) ([]authDto.UserInfoResponse, error) {
        return nil, errors.New("foo error")
    })
    req := httptest.NewRequest("GET", "/admin/users-by-role?as=SUPER_ADMIN&role=USER", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 400, resp.StatusCode)

    var body map[string]string
    json.NewDecoder(resp.Body).Decode(&body)
    assert.Equal(t, "foo error", body["error"])
}

func TestFilterUsersByRole_Success(t *testing.T) {
    app := setupAppFilterRole()
    // stub service ให้ return data
    Core_controllers.InitFilterUsersByRole(func(role string) ([]authDto.UserInfoResponse, error) {
        return []authDto.UserInfoResponse{
            {ID:1, Username:"alice", Email:"a@x", Role:"USER"},
        }, nil
    })
    req := httptest.NewRequest("GET", "/admin/users-by-role?as=SUPER_ADMIN&role=USER", nil)
    resp, _ := app.Test(req)
    assert.Equal(t, 200, resp.StatusCode)

    var list []authDto.UserInfoResponse
    json.NewDecoder(resp.Body).Decode(&list)
    assert.Len(t, list, 1)
    assert.Equal(t, "alice", list[0].Username)
}