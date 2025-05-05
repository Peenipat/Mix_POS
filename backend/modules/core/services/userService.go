package coreServices

import (
	"errors"
	"myapp/database"
	"myapp/modules/core/models"
	Core_authDto "myapp/modules/core/dto/auth"
	Core_userDto "myapp/modules/core/dto/user"
	"myapp/utils"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CreateUserFromRegister(input Core_authDto.RegisterInput) error {
	// ตรวจซ้ำ email
	var existingUser coreModels.User
	database.DB.Where("email = ?", input.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return errors.New("email already in use")
	}

	var role coreModels.Role
	if err := database.DB.Where("name = ?", coreModels.RoleNameUser).First(&role).Error; err != nil {
		return errors.New("default role not found")
	}

	// hash password ด้วย bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return errors.New("failed to hash password")
	}

	user := coreModels.User{
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword),
		RoleID:   role.ID,
	}

	// save user ลง database
	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}
func CreateUserFromAdmin(input Core_userDto.CreateUserInput) error {
    // 1) ตรวจ email ซ้ำ
    var existing coreModels.User
    database.DB.Where("email = ?", input.Email).First(&existing)
    if existing.ID != 0 {
        return errors.New("email already in use")
    }

    // 2) ป้องกันสร้าง SUPER_ADMIN
    if input.Role == string(coreModels.RoleNameSaaSSuperAdmin) {
        return errors.New("cannot create SUPER_ADMIN")
    }

    // 3) ตรวจว่า role ใน input เป็นค่าที่อนุญาต
    rn := coreModels.RoleName(input.Role)
    switch rn {
    case coreModels.RoleNameTenantAdmin,
         coreModels.RoleNameBranchAdmin,
         coreModels.RoleNameAssistantManager,
         coreModels.RoleNameStaff,
         coreModels.RoleNameUser:
        // ผ่าน
    default:
        return errors.New("invalid role provided")
    }

    // 4) หา Role record ที่ตรงกับชื่อ
    var role coreModels.Role
    if err := database.DB.Where("name = ?", rn).First(&role).Error; err != nil {
        return errors.New("specified role not found")
    }

    // 5) hash password
    hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
    if err != nil {
        return errors.New("failed to hash password")
    }

    // 6) สร้าง User พร้อมตั้ง RoleID
    user := coreModels.User{
        Username: input.Username,
        Email:    input.Email,
        Password: string(hashed),
        RoleID:   role.ID,      // <-- ใส่ FK ไปยัง roles.id
    }

    if err := database.DB.Create(&user).Error; err != nil {
        return errors.New("failed to create user")
    }
    return nil
}


func ChangeRoleFromAdmin(input Core_userDto.ChangeRoleInput) error {
    return database.DB.Transaction(func(tx *gorm.DB) error {
        var user coreModels.User

        // 1) หา user
        if err := tx.First(&user, input.ID).Error; err != nil {
            return errors.New("user not found")
        }

        // 2) ห้ามเปลี่ยนเป็น SUPER_ADMIN
        if input.Role == string(coreModels.RoleNameSaaSSuperAdmin) {
            return errors.New("cannot change role to SUPER_ADMIN")
        }

        // 3) ตรวจว่าชื่อ role ใหม่เป็นค่าที่อนุญาต
        rn := coreModels.RoleName(input.Role)
        switch rn {
        case coreModels.RoleNameTenantAdmin,
             coreModels.RoleNameBranchAdmin,
             coreModels.RoleNameAssistantManager,
             coreModels.RoleNameStaff,
             coreModels.RoleNameUser:
            // ผ่าน validation
        default:
            return errors.New("invalid role provided")
        }

        // 4) หา Role record ที่ตรงกับชื่อ
        var role coreModels.Role
        if err := tx.Where("name = ?", rn).First(&role).Error; err != nil {
            return errors.New("specified role not found")
        }

        // 5) อัปเดต FK ใน user
        user.RoleID = role.ID

        if err := tx.Save(&user).Error; err != nil {
            return errors.New("failed to update user role")
        }
        return nil
    })
}


func GetAllUsers(limit int, offset int) ([]Core_authDto.UserInfoResponse, error) {
	var users []coreModels.User
	//ค้นหา user เช็ค limit และกำหนด offset
	if err := database.DB.Order("id ASC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	result := utils.MapSlice(users, func(u coreModels.User) Core_authDto.UserInfoResponse {
		return Core_authDto.UserInfoResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			RoleID:   u.RoleID,      // รหัสบทบาท
            Role:     u.Role.Name,
		}
	})

	return result, nil
}

func FilterUsersByRole(role string) ([]Core_authDto.UserInfoResponse, error) {
	// Validate role ก่อน
	validRoles := []coreModels.RoleName{
		coreModels.RoleNameTenantAdmin,
		coreModels.RoleNameBranchAdmin,
		coreModels.RoleNameAssistantManager,
		coreModels.RoleNameStaff,
		coreModels.RoleNameUser,
	}

	isValid := false
	for _, r := range validRoles {
		if coreModels.RoleName(role) == r {
			isValid = true
			break
		}
	}
	if !isValid {
		return nil, errors.New("invalid role")
	}

	// ดึง users
	var users []coreModels.User
	if err := database.DB.Where("role = ?", role).Find(&users).Error; err != nil {
		return nil, errors.New("failed to fetch users")
	}

	// Map เป็น UserResponse
	var result []Core_authDto.UserInfoResponse
	for _, u := range users {
		result = append(result, Core_authDto.UserInfoResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			RoleID:   u.RoleID,      
            Role:     u.Role.Name,
		})
	}

	return result, nil
}
