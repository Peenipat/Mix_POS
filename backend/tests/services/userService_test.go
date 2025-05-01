package testService

import (
	// "errors"
	"myapp/database"
	"myapp/modules/core/models"
	Core_authDto "myapp/modules/core/dto/auth"
	Core_userDto "myapp/modules/core/dto/user"
	"myapp/modules/core/services"
	"myapp/tests"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test การ register ด้วยตัวเองได้ Role เป็น User [Success]
func Test_CreateUser_FromRegister_Success(t *testing.T) {
    // 1) เตรียม DB ใหม่ และ override global
    db := tests.SetupTestDB()
    database.DB = db

    // 2) สร้าง Role “USER” ลงในตารางก่อน (Service จะ lookup ตามชื่อนี้)
    userRole := coreModels.Role{Name: string(coreModels.RoleNameUser)}
    require.NoError(t, db.Create(&userRole).Error)

    // 3) ทำการ register
    input := Core_authDto.RegisterInput{
        Username: "testuser",
        Email:    "test@example.com",
        Password: "12345678",
    }
    err := services.CreateUserFromRegister(input)
    require.NoError(t, err)

    // 4) ดึง User จาก DB พร้อม preload Role เพื่อเช็คชื่อบทบาท
    var user coreModels.User
    err = db.Preload("Role").First(&user, "email = ?", input.Email).Error
    require.NoError(t, err)

    // 5) Assertions
    assert.Equal(t, input.Username, user.Username)
    assert.Equal(t, input.Email,    user.Email)
    assert.Equal(t, userRole.ID,    user.RoleID)       // FK ต้องถูก
    assert.Equal(t, userRole.Name,  user.Role.Name)   // ชื่อ role ต้องคือ “USER”
    assert.NotEmpty(t, user.Password)                 // รหัสผ่านต้อง hashed
}


// Test การ register ด้วยแต่เองกรณ๊ Email ซ้ำ
func Test_CreateUser_FromRegister_EmailAlreadyUsed(t *testing.T) {
	db := tests.SetupTestDB()

	// test กรณี Email ซ้ำกัน
	db.Create(&coreModels.User{
		Username: "exist",
		Email:    "exist@example.com",
		Password: "xxx",
	})

	input := Core_authDto.RegisterInput{
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
    // 1) เตรียม DB ใน memory แล้ว override global
    db := tests.SetupTestDB()
    database.DB = db

    // 2) สร้าง Role records สำหรับทดสอบ
    branchRole := coreModels.Role{Name: string(coreModels.RoleNameBranchAdmin)}
    staffRole  := coreModels.Role{Name: string(coreModels.RoleNameStaff)}
    userRole   := coreModels.Role{Name: string(coreModels.RoleNameUser)}
    require.NoError(t, db.Create(&branchRole).Error)
    require.NoError(t, db.Create(&staffRole).Error)
    require.NoError(t, db.Create(&userRole).Error)

    // 3) กำหนด test cases พร้อม expected Role object
    testCases := []struct {
        name         string
        input        Core_userDto.CreateUserInput
        expectedRole coreModels.Role
    }{
        {
            name: "BranchAdmin",
            input: Core_userDto.CreateUserInput{
                Username: "TestUser1",
                Email:    "test1@gmail.com",
                Password: "12345678",
                Role:     string(coreModels.RoleNameBranchAdmin),
            },
            expectedRole: branchRole,
        },
        {
            name: "Staff",
            input: Core_userDto.CreateUserInput{
                Username: "TestUser2",
                Email:    "test2@gmail.com",
                Password: "12345678",
                Role:     string(coreModels.RoleNameStaff),
            },
            expectedRole: staffRole,
        },
        {
            name: "User",
            input: Core_userDto.CreateUserInput{
                Username: "TestUser3",
                Email:    "test3@gmail.com",
                Password: "12345678",
                Role:     string(coreModels.RoleNameUser),
            },
            expectedRole: userRole,
        },
    }

    // 4) รันแต่ละกรณี
    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            // เรียก service
            err := services.CreateUserFromAdmin(tc.input)
            require.NoError(t, err)

            // โหลด User พร้อม preload Role
            var u coreModels.User
            require.NoError(t, db.Preload("Role").First(&u, "email = ?", tc.input.Email).Error)

            // Assertions
            assert.Equal(t, tc.input.Username, u.Username)
            assert.Equal(t, tc.input.Email,    u.Email)
            assert.Equal(t, tc.expectedRole.ID,   u.RoleID)     // FK ต้อง match
            assert.Equal(t, tc.expectedRole.Name, u.Role.Name)  // ชื่อ role ต้องตรง
            assert.NotEmpty(t, u.Password)                    // ต้องมีการ hash เก็บ
        })
    }
}

// Test การสร้าง User ผ่าน SuperAdmin กรณีใส่ Role ผิด เช่น สร้าง SuperAdmin หรือใส่ role ที่ไม่มีจริง
func Test_CreateUser_FromAdmin_InvalidRole(t *testing.T) {
	tests.SetupTestDB()
	testCases := []struct {
		name        string
		input       Core_userDto.CreateUserInput
		expectedErr string
	}{
		{
			name: "SuperAdmin",
			input: Core_userDto.CreateUserInput{
				Username: "TestSuperAdmin",
				Email:    "test_super_admin@gmail.com",
				Password: "12345678",
				Role:     string(coreModels.RoleNameSaaSSuperAdmin), // สร้าง User ที่เป็น Role SuperAdmin
			},
			expectedErr: "cannot create SUPER_ADMIN",
		},
		{
			name: "AnotherRole",
			input: Core_userDto.CreateUserInput{
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
	// 1) สร้าง DB สำรอง
	db := tests.SetupTestDB()

	// 2) แทนที่ database.DB ให้เป็นตัวนี้
	database.DB = db

	// 3) สร้าง Role record
	userRole := coreModels.Role{Name: string(coreModels.RoleNameUser)}
	staffRole := coreModels.Role{Name: string(coreModels.RoleNameStaff)}
	require.NoError(t, db.Create(&userRole).Error)
	require.NoError(t, db.Create(&staffRole).Error)

	// 4) สร้าง User ที่มี role_id = userRole.ID
	user := coreModels.User{
		Username: "ChangeUser",
		Email:    "change_user@example.com",
		Password: "12345678",  // ไม่ตรวจ hash ใน test
		RoleID:   userRole.ID, // ตั้ง FK
	}
	require.NoError(t, db.Create(&user).Error)

	// 5) เรียก service เปลี่ยน role
	input := Core_userDto.ChangeRoleInput{
		ID:   user.ID,
		Role: string(coreModels.RoleNameStaff),
	}
	err := services.ChangeRoleFromAdmin(input)
	require.NoError(t, err)

	// 6) โหลด User ซ้ำพร้อม Preload(Role)
	var updated coreModels.User
	require.NoError(t, db.Preload("Role").First(&updated, user.ID).Error)

	// 7) ยืนยันว่ามันเปลี่ยน RoleID และ Role.Name ถูกต้อง
	assert.Equal(t, staffRole.ID, updated.RoleID)
	assert.Equal(t, string(coreModels.RoleNameStaff), updated.Role.Name)
}

// Test การเปลี่ยน Role ผ่าน SuperAdmin กรณีพยายามเปลี่ยนเป็น SuperAdmin และ Role ที่ไม่มีจริง
func Test_ChangeRole_FromAdmin_InvalidRole(t *testing.T) {
	db := tests.SetupTestDB()

	staffRole := coreModels.Role{Name: string(coreModels.RoleNameStaff)}
	db.Create(&staffRole)
	user := coreModels.User{
		Username: "ChangeSuperAdmin",
		Email:    "changesuperadmin@gmail.com",
		Password: "12345678",
		RoleID:   staffRole.ID, 
	}
	db.Create(&user)

	testCases := []struct {
		name        string
		input       Core_userDto.ChangeRoleInput
		expectedErr string
	}{
		{
			name: "SuperAdmin",
			input: Core_userDto.ChangeRoleInput{
				ID:   user.ID,
				Role: string(coreModels.RoleNameSaaSSuperAdmin), // เปลี่ยน User ที่เป็น Role Staff เป็น SuperAdmin
			},
			expectedErr: "cannot change role to SUPER_ADMIN",
		},
		{
			name: "AnotherRole",
			input: Core_userDto.ChangeRoleInput{
				ID:   user.ID,
				Role: "HACKER", // เปลี่ยน User ที่เป็น Role Staff เป็น Role อื่นๆ
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
func Test_GetAllUser_limitData(t *testing.T) {
	db := tests.SetupTestDB()

	userRole := coreModels.Role{Name: string(coreModels.RoleNameStaff)}
	db.Create(&userRole)
	// Mock Data
	users := []coreModels.User{
		{Username: "User1", Email: "user1@example.com", Password: "xx", RoleID: userRole.ID},
		{Username: "User2", Email: "user2@example.com", Password: "xx", RoleID: userRole.ID},
		{Username: "User3", Email: "user3@example.com", Password: "xx", RoleID: userRole.ID},
		{Username: "User4", Email: "user4@example.com", Password: "xx", RoleID: userRole.ID},
		{Username: "User5", Email: "user5@example.com", Password: "xx", RoleID: userRole.ID},
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
