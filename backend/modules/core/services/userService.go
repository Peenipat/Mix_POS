package coreServices

import (
	"errors"
	"context"
	"myapp/database"
	"myapp/modules/core/models"
	corePort "myapp/modules/core/port"
	"myapp/utils"
	"fmt"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func NewUserService(db *gorm.DB) corePort.IUser {
	return &UserService{DB: db}
}

var (
    ErrUserNotFound       = errors.New("user not found")
    ErrInvalidOldPassword = errors.New("old password is incorrect")
)


func (u *UserService) CreateUserFromRegister(input corePort.RegisterInput) error {

	var existingUser coreModels.User
	database.DB.Where("email = ?", input.Email).First(&existingUser)
	if existingUser.ID != 0 {
		return errors.New("email already in use")
	}

	var role coreModels.Role
	if err := database.DB.Where("name = ?", coreModels.RoleNameUser).First(&role).Error; err != nil {
		return errors.New("default role not found")
	}

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

	if err := database.DB.Create(&user).Error; err != nil {
		return errors.New("failed to create user")
	}

	return nil
}


func (u *UserService) CreateUserFromAdmin(input corePort.CreateUserInput) error {
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


func (u *UserService) ChangeRoleFromAdmin(input corePort.ChangeRoleInput) error {
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

func (u *UserService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
    if userID == 0 {
        return fmt.Errorf("user ID is required")
    }

    var user coreModels.User
    if err := u.DB.WithContext(ctx).
        Where("id = ?", userID).
        First(&user).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrUserNotFound
        }
        return fmt.Errorf("fetch user %d: %w", userID, err)
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
        return ErrInvalidOldPassword
    }

    hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
    if err != nil {
        return fmt.Errorf("hash new password: %w", err)
    }

    // 4) อัพเดตใน DB
    if err := u.DB.WithContext(ctx).
        Model(&user).
        Update("password", string(hashed)).Error; err != nil {
        return fmt.Errorf("update password: %w", err)
    }

    return nil
}

func (u *UserService) GetAllUsers(limit int, offset int) ([]corePort.UserInfoResponse, error) {
	var users []coreModels.User
	//ค้นหา user เช็ค limit และกำหนด offset
	if err := database.DB.Order("id ASC").Limit(limit).Offset(offset).Find(&users).Error; err != nil {
		return nil, err
	}

	result := utils.MapSlice(users, func(u coreModels.User) corePort.UserInfoResponse {
		return corePort.UserInfoResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			RoleID:   u.RoleID,      // รหัสบทบาท
            Role:     u.Role.Name,
		}
	})

	return result, nil
}

func (u *UserService) FilterUsersByRole(role string) ([]corePort.UserInfoResponse, error) {
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
	var result []corePort.UserInfoResponse
	for _, u := range users {
		result = append(result, corePort.UserInfoResponse{
			ID:       u.ID,
			Username: u.Username,
			Email:    u.Email,
			RoleID:   u.RoleID,      
            Role:     u.Role.Name,
		})
	}

	return result, nil
}


//ChangePassword เอาไว้ก่อน
//Authenticate (ยืนยันรหัสผ่าน) เอาไว้ก่อน
//ListUsers (ที่มี UserFilter) filter by role/branch/tenant
// type UserFilter struct {
//     RoleID   *uint
//     BranchID *uint
//     TenantID *uint  
//     Active   *bool  
//     Page     int
//     PageSize int
// }