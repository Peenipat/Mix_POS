package testService

import (
	// "errors"
	"myapp/database"
	authDto "myapp/dto/auth"
	userDto "myapp/dto/user"
	"myapp/models"
	"myapp/services"
	"testing"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

//เชื่อมต่อ database เข้า memory  
func setupTestDB() *gorm.DB {
	db, _ := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	db.AutoMigrate(&models.User{})
	database.DB = db
	return db
}

// Test การ register ด้วยตัวเองได้ Role เป็น User [Success]
func Test_CreateUser_FromRegister_Success(t *testing.T) {
	setupTestDB()

	input := authDto.RegisterInput{
		Username: "testuser",
		Email:    "test@example.com",
		Password: "12345678",
	}

	err := services.CreateUserFromRegister(input) // เรียก service มาลอง test
	assert.Nil(t, err)

	var user models.User
	database.DB.First(&user, "email = ?", input.Email)
	assert.Equal(t, input.Username, user.Username)
	assert.Equal(t, input.Email, user.Email)    // เช็ค DB ว่า email ตรงกับ input
	assert.Equal(t, models.RoleUser, user.Role) // เช็คว่า role ต้องเป็น User
	assert.NotEmpty(t, user.Password)           // ต้องมีการ hash
}

// Test การ register ด้วยแต่เองกรณ๊ Email ซ้ำ
func Test_CreateUser_FromRegister_EmailAlreadyUsed(t *testing.T) {
	db := setupTestDB()

	// test กรณี Email ซ้ำกัน
	db.Create(&models.User{
		Username: "exist",
		Email:    "exist@example.com",
		Password: "xxx",
	})

	input := authDto.RegisterInput{
		Username: "newuser",
		Email:    "exist@example.com",
		Password: "12345678",
	}

	err := services.CreateUserFromRegister(input)
	assert.NotNil(t, err)
	assert.Equal(t, "email already in use", err.Error())
}

// Test การสร้าง User ผ่าน Admin ได้ Role เป็น BranchAdmin, Staff, User [Success]
func Test_CreateUser_FromAdmin_Success(t *testing.T) {
	db := setupTestDB()
	testCases := []struct {
		name  string
		input userDto.CreateUserInput
	}{
		{
			name: "BranchAdmin",
			input: userDto.CreateUserInput{
				Username: "TestUser1",
				Email:    "test1@gmail.com",
				Password: "12345678",
				Role:     string(models.RoleBranchAdmin), // สร้าง User ที่เป็น Role BranchAdmin
			},
		},
		{
			name: "Staff",
			input: userDto.CreateUserInput{
				Username: "TestUser2",
				Email:    "test2@gmail.com",
				Password: "12345678",
				Role:     string(models.RoleStaff), // สร้าง User ที่เป็น Role Staff
			},
		},
		{
			name: "User",
			input: userDto.CreateUserInput{
				Username: "TestUser3",
				Email:    "test3@gmail.com",
				Password: "12345678",
				Role:     string(models.RoleUser), // สร้าง User ที่เป็น Role User
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := services.CreateUserFromAdmin(tc.input)
			assert.Nil(t, err)

			var user models.User
			db.First(&user, "email = ?", tc.input.Email)
			assert.Equal(t, tc.input.Username, user.Username)
			assert.Equal(t, tc.input.Email, user.Email)            // เช็ค DB ว่า email ตรงกับ input
			assert.Equal(t, models.Role(tc.input.Role), user.Role) // เช็คว่า role ต้องเป็น BranchAdmin,Staff,User
			assert.NotEmpty(t, user.Password)

		})
	}
}

// Test การสร้าง User ผ่าน SuperAdmin กรณีใส่ Role ผิด เช่น สร้าง SuperAdmin หรือใส่ role ที่ไม่มีจริง
func Test_CreateUser_FromAdmin_InvalidRole(t *testing.T) {
	setupTestDB()
	testCases := []struct {
		name        string
		input       userDto.CreateUserInput
		expectedErr string
	}{
		{
			name: "SuperAdmin",
			input: userDto.CreateUserInput{
				Username: "TestSuperAdmin",
				Email:    "test_super_admin@gmail.com",
				Password: "12345678",
				Role:     string(models.RoleSuperAdmin), // สร้าง User ที่เป็น Role SuperAdmin
			},
			expectedErr: "cannot create SUPER_ADMIN",
		},
		{
			name: "AnotherRole",
			input: userDto.CreateUserInput{
				Username: "TestUserAnother",
				Email:    "test_another_role@gmail.com",
				Password: "12345678",
				Role:     "HACKER", // สร้าง User ที่เป็น Role อื่น ๆ
			},
			expectedErr: "invalid role provided",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := services.CreateUserFromAdmin(tc.input)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedErr, err.Error())
		})
	}
}

// Test การเปลี่ยน Role ผ่าน SuperAdmin ได้ Role เป็น BranchAdmin, Staff, User [Success]
func Test_ChangeRole_FromAdmin_Success(t *testing.T) {
	db := setupTestDB()

	user := models.User{
		Username: "ChangeUser",
		Email:    "change_user@example.com",
		Password: "12345678",
		Role:     models.RoleUser,
	}
	db.Create(&user)

	input := userDto.ChangeRoleInput{
		ID:   user.ID,
		Role: string(models.RoleStaff),
	}

	err := services.ChangeRoleFromAdmin(input)
	assert.Nil(t, err)

	var updated models.User
	db.First(&updated, user.ID)
	assert.Equal(t, models.RoleStaff, updated.Role)
}

// Test การเปลี่ยน Role ผ่าน SuperAdmin กรณีพยายามเปลี่ยนเป็น SuperAdmin และ Role ที่ไม่มีจริง
func Test_ChangeRole_FromAdmin_InvalidRole(t *testing.T) {
	db := setupTestDB()

	user := models.User{
		Username: "ChangeSuperAdmin",
		Email:    "changesuperadmin@gmail.com",
		Password: "12345678",
		Role:     models.RoleStaff,
	}
	db.Create(&user)

	testCases := []struct {
		name        string
		input       userDto.ChangeRoleInput
		expectedErr string
	}{
		{
			name: "SuperAdmin",
			input: userDto.ChangeRoleInput{
			ID: user.ID,
			Role: string(models.RoleSuperAdmin), // เปลี่ยน User ที่เป็น Role Staff เป็น SuperAdmin
			},
			expectedErr: "cannot change role to SUPER_ADMIN",
		},
		{
			name: "AnotherRole",
			input: userDto.ChangeRoleInput{
			ID: user.ID,
			Role:   "HACKER", // เปลี่ยน User ที่เป็น Role Staff เป็น Role อื่นๆ
			},
			expectedErr: "invalid role provided",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := services.ChangeRoleFromAdmin(tc.input)
			assert.NotNil(t, err)
			assert.Equal(t, tc.expectedErr, err.Error())
		})
	}

}

// Test การดึงข้อมูล User โดยที่สามารถใส่ limit ได้ [Success] 
func Test_GetAllUser_limitData(t *testing.T){
	db := setupTestDB()

	// Mock Data
	users := []models.User{
		{Username: "User1", Email: "user1@example.com", Password: "xx", Role: models.RoleUser},
		{Username: "User2", Email: "user2@example.com", Password: "xx", Role: models.RoleUser},
		{Username: "User3", Email: "user3@example.com", Password: "xx", Role: models.RoleUser},
		{Username: "User4", Email: "user4@example.com", Password: "xx", Role: models.RoleUser},
		{Username: "User5", Email: "user5@example.com", Password: "xx", Role: models.RoleUser},
	}
	for _, u := range users {
		db.Create(&u)
	}
	// call service
	result, err := services.GetAllUsers(3, 0)

	// check
	assert.Nil(t, err)
	assert.Len(t, result, 3)
	assert.Equal(t, "User1", result[0].Username)
	assert.Equal(t, "User2", result[1].Username)
	assert.Equal(t, "User3", result[2].Username)
}


