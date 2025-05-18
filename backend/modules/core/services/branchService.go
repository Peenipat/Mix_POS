package coreServices

import (
	"errors"
	"gorm.io/gorm"
	coreModels "myapp/modules/core/models"
	"strings"
)

type BranchService struct {
	DB *gorm.DB
}

func NewBranchService(db *gorm.DB) *BranchService {
	return &BranchService{DB: db}
}

// Create
func (s *BranchService) CreateBranch(branch *coreModels.Branch) error {
	//  Trim input ก่อน
	branch.Name = strings.TrimSpace(branch.Name)

	// Validate input
	if branch.TenantID == 0 || branch.Name == "" {
		return errors.New("tenant_id and name are required")
	}

	// Save
	if err := s.DB.Create(branch).Error; err != nil {
		if strings.Contains(err.Error(), "idx_tenant_name") {
			return errors.New("branch name already exists for this tenant")
		}
		return err
	}

	return nil
}

// Read All
func (s *BranchService) GetAllBranches(tenantID uint) ([]coreModels.Branch, error) {
	var branches []coreModels.Branch
	err := s.DB.Where("tenant_id = ?", tenantID).Find(&branches).Error
	return branches, err
}

// Read by ID
func (s *BranchService) GetBranchByID(id uint) (*coreModels.Branch, error) {
	var branch coreModels.Branch
	if err := s.DB.
		Where("id = ?", id).
		Where("deleted_at IS NULL"). // ❗ไม่โหลด soft-deleted
		First(&branch).Error; err != nil {
		return nil, err
	}
	return &branch, nil
}

// Update
func (s *BranchService) UpdateBranch(branch *coreModels.Branch) error {
	branch.Name = strings.TrimSpace(branch.Name)
	if branch.ID == 0 {
		return errors.New("branch ID is required")
	}
	if strings.TrimSpace(branch.Name) == "" {
		return errors.New("branch name is required")
	}

	// เช็คก่อนว่ามีจริง
	var existing coreModels.Branch
	if err := s.DB.First(&existing, branch.ID).Error; err != nil {
		return err
	}

	return s.DB.Model(&existing).Updates(branch).Error
}

// Delete (Soft Delete)
func (s *BranchService) DeleteBranch(id uint) error {
	var branch coreModels.Branch
	if err := s.DB.First(&branch, id).Error; err != nil {
		return err
	}
	return s.DB.Delete(&branch).Error
}

func (s *BranchService) GetBranchesByTenantID(tenantID uint) ([]coreModels.Branch, error) {
	var branches []coreModels.Branch
	if err := s.DB.Where("tenant_id = ?", tenantID).Find(&branches).Error; err != nil {
		return nil, err
	}

	return branches, nil
}
