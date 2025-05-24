package coreServices

import (
	"errors"
	"gorm.io/gorm"
	coreModels "myapp/modules/core/models"
	"fmt"
    "strings"

    "github.com/jackc/pgconn"
)


type BranchPort interface {
	CreateBranch(b *coreModels.Branch) error
	GetAllBranches() ([]coreModels.Branch, error) // saas admin line
	GetBranchByID(id uint) (*coreModels.Branch, error)
    UpdateBranch(*coreModels.Branch) error  
    DeleteBranch(uint) error 
    GetBranchesByTenantID(tenantID uint) ([]coreModels.Branch, error)
  }

type BranchService struct {
	DB *gorm.DB
}

func NewBranchService(db *gorm.DB) *BranchService {
	return &BranchService{DB: db}
}

var (
    ErrInvalidInput        = errors.New("tenant_id and name are required")
    ErrTenantNotFound      = errors.New("tenant not found")
	ErrBranchNotFound      = errors.New("branch not found")
    ErrBranchExists        = errors.New("branch name already exists for this tenant")
    ErrForeignKey          = errors.New("foreign key violation")
	ErrInvalidID           = errors.New("invalid branch ID")
	ErrFetchBranchesFailed = errors.New("failed to fetch branches")
    ErrBranchInUse         = errors.New("branch cannot be deleted: still has dependent records")
    ErrInvalidTenantID     = errors.New("tenant ID is required")
    ErrNameRequired        = errors.New("branch name is required")
)

// Create
// CreateBranch trims input, validates, checks tenant existence, and handles DB errors robustly.
func (s *BranchService) CreateBranch(branch *coreModels.Branch) error {
    // Trim name whitespace
    branch.Name = strings.TrimSpace(branch.Name)

    // Validate input
    if branch.TenantID == 0 || branch.Name == "" {
        return ErrInvalidInput
    }
    // Enforce a reasonable max length on name
    if len(branch.Name) > 100 {
        return fmt.Errorf("name too long: maximum is 100 characters")
    }

    // Verify tenant exists
    var tenant coreModels.Tenant
    if err := s.DB.First(&tenant, branch.TenantID).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrTenantNotFound
        }
        return fmt.Errorf("failed to verify tenant: %w", err)
    }

    // Attempt to save
    if err := s.DB.Create(branch).Error; err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23505": // unique_violation
				return ErrBranchExists
			case "23503": // foreign_key_violation
				return fmt.Errorf("%w: %s", ErrForeignKey, pgErr.ConstraintName)
			}
		}
		return fmt.Errorf("failed to create branch: %w", err)
	}

    return nil
}

// Read All
func (s *BranchService) GetAllBranches() ([]coreModels.Branch, error) {
    var branches []coreModels.Branch
    if err := s.DB.
        Order("created_at DESC").
        Find(&branches).Error; err != nil {
        return nil, fmt.Errorf("%w: %v", ErrFetchBranchesFailed, err)
    }
    if branches == nil {
        branches = make([]coreModels.Branch, 0)
    }
    return branches, nil
}

// Read by ID
func (s *BranchService) GetBranchByID(id uint) (*coreModels.Branch, error) {
    // 1) Validate the input ID
    if id == 0 {
        return nil, ErrInvalidID
    }

    // 2) Query for the branch, skipping soft-deleted rows
    var branch coreModels.Branch
    err := s.DB.
        Where("id = ?", id).
        Where("deleted_at IS NULL").
        First(&branch).Error

    // 3) Handle not-found vs other errors
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrBranchNotFound
        }
        return nil, fmt.Errorf("error fetching branch %d: %w", id, err)
    }

    // 4) Success
    return &branch, nil
}

func (s *BranchService) UpdateBranch(branch *coreModels.Branch) error {
    // 1) Validate ID
    if branch.ID == 0 {
        return ErrInvalidID
    }

    // 2) Trim and validate Name
    branch.Name = strings.TrimSpace(branch.Name)
    if branch.Name == "" {
        return ErrNameRequired
    }

    // 3) Fetch existing record (skip soft-deleted)
    var existing coreModels.Branch
    err := s.DB.
        Where("id = ?", branch.ID).
        Where("deleted_at IS NULL").
        First(&existing).Error

    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrBranchNotFound
        }
        return fmt.Errorf("error fetching branch %d: %w", branch.ID, err)
    }

    // 4) Perform update
    if err := s.DB.
        Model(&existing).
        Select("name").   // ป้องกัน field อื่นถูกเปลี่ยนถ้าไม่ได้ใส่มา
        Updates(coreModels.Branch{Name: branch.Name}).Error; err != nil {
        return fmt.Errorf("error updating branch %d: %w", branch.ID, err)
    }

    return nil
}

// Delete (Soft Delete)
func (s *BranchService) DeleteBranch(id uint) error {
    if id == 0 {
        return ErrInvalidID
    }

    var branch coreModels.Branch
    err := s.DB.
        Where("id = ?", id).
        Where("deleted_at IS NULL").
        First(&branch).Error
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return ErrBranchNotFound
        }
        return fmt.Errorf("error fetching branch %d: %w", id, err)
    }

    var count int64
    if err := s.DB.Model(&coreModels.User{}).
        Where("branch_id = ?", id).
        Count(&count).Error; err != nil {
        return fmt.Errorf("error checking branch dependencies: %w", err)
    }
    if count > 0 {
        return ErrBranchInUse
    }

    // 4) Soft-delete
    if err := s.DB.Delete(&branch).Error; err != nil {
        return fmt.Errorf("error deleting branch %d: %w", id, err)
    }

    return nil
}

func (s *BranchService) GetBranchesByTenantID(tenantID uint) ([]coreModels.Branch, error) {
    // 1) Validate the input ID
    if tenantID == 0 {
        return nil, ErrInvalidTenantID
    }

    // 2) Check tenant exists (no need to order here)
    var tenant coreModels.Tenant
    if err := s.DB.
        Select("id").
        Where("id = ?", tenantID).
        Where("deleted_at IS NULL").
        First(&tenant).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, ErrTenantNotFound
        }
        return nil, fmt.Errorf("error fetching tenant %d: %w", tenantID, err)
    }

    // 3) Query branches and order by created_at DESC
    var branches []coreModels.Branch
    if err := s.DB.
        Where("tenant_id = ?", tenantID).
        Where("deleted_at IS NULL").
        Order("created_at DESC").    // ← ใส่ที่นี่
        Find(&branches).Error; err != nil {
        return nil, fmt.Errorf("error fetching branches for tenant %d: %w", tenantID, err)
    }

    return branches, nil
}

