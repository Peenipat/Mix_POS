package coreServices

import (
	"context"
	"gorm.io/gorm"
	"errors"
	"strings"
	"fmt"
	"time"
	coreModels "myapp/modules/core/models"
	corePort "myapp/modules/core/port"
)

type TenantService struct {
	DB *gorm.DB
}

func NewTenantService(db *gorm.DB) corePort.ITenant {
	return &TenantService{DB: db}
}

var (
    ErrInvalidTenantInput = errors.New("invalid tenant input")
    ErrDomainTaken        = errors.New("domain already in use")
	ErrFetchTenantsFailed = errors.New("failed to fetch tenants")
    ErrUpdateFailed       = errors.New("failed to update tenant")
    ErrTenantInUse       = errors.New("tenant still in use by users or branches")
    ErrDeleteTenantFail  = errors.New("failed to delete tenant")
)


func (s *TenantService) CreateTenant(ctx context.Context, input corePort.CreateTenantInput) (*coreModels.Tenant, error) {
    name := strings.TrimSpace(input.Name)
    domain := strings.TrimSpace(input.Domain)
    if name == "" || domain == "" {
        return nil, ErrInvalidTenantInput
    }
    var exists int64
    if err := s.DB.WithContext(ctx).
        Model(&coreModels.Tenant{}).
        Where("domain = ?", domain).
        Count(&exists).Error; err != nil {
        return nil, fmt.Errorf("check domain uniqueness: %w", err)
    }
    if exists > 0 {
        return nil, ErrDomainTaken
    }

    tenant := &coreModels.Tenant{
        Name:      name,
        Domain:    domain,
        IsActive:  true,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    if err := s.DB.WithContext(ctx).Create(tenant).Error; err != nil {
        return nil, fmt.Errorf("create tenant: %w", err)
    }

    return tenant, nil
}

func (s *TenantService) GetTenantByID(ctx context.Context, id uint) (*coreModels.Tenant, error) {
    // 1) Validate input
    if id == 0 {
        return nil, ErrInvalidTenantID
    }

    // 2) Query for tenant, skip soft-deleted
    var tenant coreModels.Tenant
    err := s.DB.
        WithContext(ctx).
        Where("id = ?", id).
        Where("deleted_at IS NULL").
        First(&tenant).Error

    // 3) Handle not-found vs other errors
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrTenantNotFound
        }
        return nil, fmt.Errorf("fetch tenant %d: %w", id, err)
    }

    // 4) Success
    return &tenant, nil
}


func (s *TenantService) DeleteTenant(ctx context.Context, id uint) error {
    // 1) Validate ID
    if id == 0 {
        return ErrInvalidTenantID
    }

    // 2) Load tenant (skip soft‐deleted)
    var t coreModels.Tenant
    if err := s.DB.
        WithContext(ctx).
        Where("id = ?", id).
        Where("deleted_at IS NULL").
        First(&t).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrTenantNotFound
        }
        return fmt.Errorf("fetch tenant %d: %w", id, err)
    }

    // 3) Check for any users still assigned to this tenant
    var usageCount int64
    if err := s.DB.
        WithContext(ctx).
        Model(&coreModels.TenantUser{}).
        Where("tenant_id = ?", id).
        Count(&usageCount).Error; err != nil {
        return fmt.Errorf("check tenant users: %w", err)
    }
    if usageCount > 0 {
        return ErrTenantInUse
    }

    // 4) (Optional) Check for branches
    var branchCount int64
    if err := s.DB.
        WithContext(ctx).
        Model(&coreModels.Branch{}).
        Where("tenant_id = ?", id).
        Count(&branchCount).Error; err != nil {
        return fmt.Errorf("check tenant branches: %w", err)
    }
    if branchCount > 0 {
        return ErrTenantInUse
    }

    // 5) Soft‐delete
    if err := s.DB.
        WithContext(ctx).
        Delete(&t).Error; err != nil {
        return fmt.Errorf("%w: %v", ErrDeleteTenantFail, err)
    }

    return nil
}


func (s *TenantService) ListTenants(ctx context.Context, onlyActive bool) ([]coreModels.Tenant, error) {
    var tenants []coreModels.Tenant
    db := s.DB.WithContext(ctx).
        Model(&coreModels.Tenant{}).
        Where("deleted_at IS NULL")

    if onlyActive {
        db = db.Where("is_active = ?", true)
    }

    if err := db.
        Order("created_at DESC").
        Find(&tenants).Error; err != nil {
        return nil, fmt.Errorf("%w: %v", ErrFetchTenantsFailed, err)
    }

    // ensure non-nil slice
    if tenants == nil {
        tenants = make([]coreModels.Tenant, 0)
    }
    return tenants, nil
}

func (s *TenantService) UpdateTenant(ctx context.Context, input corePort.UpdateTenantInput) error {
    // 1) Validate ID
    if input.ID == 0 {
        return ErrInvalidTenantID
    }

    // 2) Validate fields
    if input.Name != nil {
        *input.Name = strings.TrimSpace(*input.Name)
        if *input.Name == "" {
            return ErrInvalidTenantInput
        }
    }
    if input.Domain != nil {
        *input.Domain = strings.TrimSpace(*input.Domain)
        if *input.Domain == "" {
            return ErrInvalidTenantInput
        }
    }

    // 3) Check domain uniqueness if changing domain
    if input.Domain != nil {
        var count int64
        if err := s.DB.WithContext(ctx).
            Model(&coreModels.Tenant{}).
            Where("domain = ?", *input.Domain).
            Where("id <> ?", input.ID).
            Count(&count).Error; err != nil {
            return fmt.Errorf("domain check: %w", err)
        }
        if count > 0 {
            return ErrDomainTaken
        }
    }

    // 4) Load existing
    var tenant coreModels.Tenant
    err := s.DB.WithContext(ctx).
        Where("id = ?", input.ID).
        Where("deleted_at IS NULL").
        First(&tenant).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrTenantNotFound
        }
        return fmt.Errorf("fetch tenant: %w", err)
    }

    // 5) Apply updates
    updates := make(map[string]interface{})
    if input.Name != nil {
        updates["name"] = *input.Name
    }
    if input.Domain != nil {
        updates["domain"] = *input.Domain
    }
    if input.IsActive != nil {
        updates["is_active"] = *input.IsActive
    }
    if len(updates) == 0 {
        // nothing to update
        return nil
    }

    if err := s.DB.WithContext(ctx).
        Model(&tenant).
        Updates(updates).Error; err != nil {
        return fmt.Errorf("%w: %v", ErrUpdateFailed, err)
    }

    return nil
}


