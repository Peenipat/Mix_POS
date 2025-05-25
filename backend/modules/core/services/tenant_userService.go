package coreServices

import (
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"

	coreModels "myapp/modules/core/models"
	corePort "myapp/modules/core/port"
)

var (
	ErrInvalidUserID       = errors.New("invalid user ID")
	ErrUserAlreadyAssigned = errors.New("user already assigned to tenant")
	ErrUserNotAssigned = errors.New("user not assigned to tenant")
    ErrNoTenantsAssigned  = errors.New("no tenants assigned to user")
)

// TenantUserService handles M2M between tenants and users.
type TenantUserService struct {
	DB *gorm.DB
}

// IsUserInTenant implements corePort.ITenantUser.
func (s *TenantUserService) IsUserInTenant(ctx context.Context, tenantID uint, userID uint) (bool, error) {
	panic("unimplemented")
}

func (s *TenantUserService) ListTenantsByUser(ctx context.Context, userID uint) ([]coreModels.Tenant, error) {
    // 1. ตรวจ userID เบื้องต้น
    if userID == 0 {
        return nil, ErrInvalidUserID
    }

    // 2. ตรวจว่ามี user ตัวนั้นอยู่ในระบบหรือไม่
    var cnt int64
    if err := s.DB.
        WithContext(ctx).
        Model(&coreModels.User{}).
        Where("id = ?", userID).
        Count(&cnt).
        Error; err != nil {
        return nil, fmt.Errorf("checking user existence: %w", err)
    }
    if cnt == 0 {
        return nil, ErrUserNotFound
    }

    // 3. Join ดึง tenants ที่ user ถูก assign
    var tenants []coreModels.Tenant
    err := s.DB.
        WithContext(ctx).
        Model(&coreModels.Tenant{}).
        Joins("JOIN tenant_users tu ON tu.tenant_id = tenants.id").
        Where("tu.user_id = ?", userID).
        Find(&tenants).
        Error
    if err != nil {
        return nil, fmt.Errorf("querying assigned tenants: %w", err)
    }

    // 4. ถ้าไม่มี tenant ไหนถูก assign เลย
    if len(tenants) == 0 {
        return nil, ErrNoTenantsAssigned
    }

    // 5. คืนผลลัพธ์
    return tenants, nil
}

// ListUsersByTenant implements corePort.ITenantUser.
func (s *TenantUserService) ListUsersByTenant(ctx context.Context, tenantID uint) ([]coreModels.User, error) {
	panic("unimplemented")
}




// NewTenantUserService constructs a new service.
func NewTenantUserService(db *gorm.DB) corePort.ITenantUser {
	return &TenantUserService{DB: db}
}

// AddUserToTenant associates a user with a tenant.
func (s *TenantUserService) AddUserToTenant(ctx context.Context, tenantID, userID uint) error {
	// 1) Validate IDs
	if tenantID == 0 {
		return ErrInvalidTenantID
	}
	if userID == 0 {
		return ErrInvalidUserID
	}

	// 2) Check tenant exists (and not soft-deleted)
	var t coreModels.Tenant
	if err := s.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", tenantID).
		First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrTenantNotFound
		}
		return fmt.Errorf("fetch tenant %d: %w", tenantID, err)
	}

	// 3) Check user exists (and not soft-deleted)
	var u coreModels.User
	if err := s.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", userID).
		First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return fmt.Errorf("fetch user %d: %w", userID, err)
	}

	// 4) Check existing assignment
	var tu coreModels.TenantUser
	if err := s.DB.WithContext(ctx).
		Where("tenant_id = ? AND user_id = ?", tenantID, userID).
		First(&tu).Error; err == nil {
		return ErrUserAlreadyAssigned
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("check existing assignment: %w", err)
	}

	// 5) Create assignment
	newTU := coreModels.TenantUser{
		TenantID: tenantID,
		UserID:   userID,
	}
	if err := s.DB.WithContext(ctx).Create(&newTU).Error; err != nil {
		return fmt.Errorf("assign user %d to tenant %d: %w", userID, tenantID, err)
	}

	return nil
}

func (s *TenantUserService) RemoveUserFromTenant(ctx context.Context, tenantID, userID uint) error {
    // 1) Validate IDs
    if tenantID == 0 {
        return ErrInvalidTenantID
    }
    if userID == 0 {
        return ErrInvalidUserID
    }

    // 2) Ensure tenant exists
    var t coreModels.Tenant
    if err := s.DB.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", tenantID).
        First(&t).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrTenantNotFound
        }
        return fmt.Errorf("fetch tenant %d: %w", tenantID, err)
    }

    // 3) Ensure user exists
    var u coreModels.User
    if err := s.DB.WithContext(ctx).
        Where("id = ? AND deleted_at IS NULL", userID).
        First(&u).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrUserNotFound
        }
        return fmt.Errorf("fetch user %d: %w", userID, err)
    }

    // 4) Find mapping
    var tu coreModels.TenantUser
    if err := s.DB.WithContext(ctx).
        Where("tenant_id = ? AND user_id = ?", tenantID, userID).
        First(&tu).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrUserNotAssigned
        }
        return fmt.Errorf("fetch assignment: %w", err)
    }

    // 5) Delete mapping
    if err := s.DB.WithContext(ctx).Delete(&tu).Error; err != nil {
        return fmt.Errorf("remove user %d from tenant %d: %w", userID, tenantID, err)
    }

    return nil
}
