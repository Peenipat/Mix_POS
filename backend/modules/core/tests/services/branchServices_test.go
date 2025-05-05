package coreTest

import (
	"testing"
	"strings"
	

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	coreModels "myapp/modules/core/models"
	coreServices "myapp/modules/core/services"
)

func setupTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&coreModels.Branch{})
	assert.NoError(t, err)

	return db
}

func TestBranchService(t *testing.T) {
	db := setupTestDB(t)
	service := coreServices.NewBranchService(db)

	t.Run("CreateBranch", func(t *testing.T) {
		branch := &coreModels.Branch{
			TenantID: 1,
			Name:     "Main Branch",
		}
		err := service.CreateBranch(branch)
		assert.NoError(t, err)
		assert.NotZero(t, branch.ID)
	})

	t.Run("CreateBranch_InvalidData", func(t *testing.T) {
		branch := &coreModels.Branch{
			TenantID: 1,
			Name:     "", //  ว่าง
		}
		err := service.CreateBranch(branch)
		assert.Error(t, err)
	})

	t.Run("CreateBranch_MissingTenantID", func(t *testing.T) {
		branch := &coreModels.Branch{Name: "NoTenant"}
		err := service.CreateBranch(branch)
		assert.Error(t, err)
	})

	t.Run("CreateBranch_TrimmedName", func(t *testing.T) {
		branch := &coreModels.Branch{
			TenantID: 8,
			Name:     "   Trimmed Name   ",
		}
		err := service.CreateBranch(branch)
		assert.NoError(t, err)
	
		result, _ := service.GetBranchByID(branch.ID)
		assert.Equal(t, "Trimmed Name", result.Name)
	})

	t.Run("CreateBranch_DuplicateNamePerTenant", func(t *testing.T) {
		tenantID := uint(55)
		branch1 := &coreModels.Branch{TenantID: tenantID, Name: "Dup Branch"}
		branch2 := &coreModels.Branch{TenantID: tenantID, Name: "Dup Branch"}
	
		_ = service.CreateBranch(branch1)
		err := service.CreateBranch(branch2)
		assert.Error(t, err) // expect duplicate error
	})
	

	t.Run("GetBranchByID", func(t *testing.T) {
		branch := &coreModels.Branch{
			TenantID: 2,
			Name:     "Second Branch",
		}
		_ = service.CreateBranch(branch)

		result, err := service.GetBranchByID(branch.ID)
		assert.NoError(t, err)
		assert.Equal(t, branch.Name, result.Name)
	})

	t.Run("GetBranchByID_NotFound", func(t *testing.T) {
		_, err := service.GetBranchByID(9999999)
		assert.Error(t, err)
	})

	t.Run("GetBranchByID_Deleted", func(t *testing.T) {
		branch := &coreModels.Branch{TenantID: 1, Name: "ToDelete"}
		_ = service.CreateBranch(branch)
		_ = service.DeleteBranch(branch.ID)
	
		_, err := service.GetBranchByID(branch.ID)
		assert.Error(t, err)
	})

	t.Run("GetBranchesByTenantID", func(t *testing.T) {
		tenantID := uint(99)
	
		// เตรียม data
		branch1 := &coreModels.Branch{TenantID: tenantID, Name: "Branch A"}
		branch2 := &coreModels.Branch{TenantID: tenantID, Name: "Branch B"}
		_ = service.CreateBranch(branch1)
		_ = service.CreateBranch(branch2)
	
		// อีก tenant
		_ = service.CreateBranch(&coreModels.Branch{TenantID: tenantID + 1, Name: "Other Branch"})
	
		branches, err := service.GetBranchesByTenantID(tenantID)
		assert.NoError(t, err)
		assert.Len(t, branches, 2)
		assert.ElementsMatch(t, []string{"Branch A", "Branch B"}, []string{branches[0].Name, branches[1].Name})
	})

	t.Run("GetBranchesByTenantID_EmptyResult", func(t *testing.T) {
		branches, err := service.GetBranchesByTenantID(999999) // unlikely tenant ID
		assert.NoError(t, err)
		assert.Len(t, branches, 0)
	})
	
	t.Run("UpdateBranch", func(t *testing.T) {
		branch := &coreModels.Branch{
			TenantID: 3,
			Name:     "Old Name",
		}
		_ = service.CreateBranch(branch)

		branch.Name = "New Name"
		err := service.UpdateBranch(branch)
		assert.NoError(t, err)

		result, _ := service.GetBranchByID(branch.ID)
		assert.Equal(t, "New Name", result.Name)
	})

	t.Run("UpdateBranch_NotExist", func(t *testing.T) {
		branch := &coreModels.Branch{
			ID:       9999999,
			TenantID: 1,
			Name:     "Ghost Branch",
		}
		err := service.UpdateBranch(branch)
		assert.Error(t, err)
	})

	t.Run("UpdateBranch_TrimmedName", func(t *testing.T) {
		branch := &coreModels.Branch{TenantID: 2, Name: "   MyBranch   "}
		_ = service.CreateBranch(branch)
	
		branch.Name = "   Updated Name   "
		err := service.UpdateBranch(branch)
		assert.NoError(t, err)
	
		result, _ := service.GetBranchByID(branch.ID)
		assert.Equal(t, strings.TrimSpace(branch.Name), result.Name)
	})

	t.Run("UpdateBranch_MissingRequiredField", func(t *testing.T) {
		branch := &coreModels.Branch{TenantID: 77, Name: "ToEdit"}
		_ = service.CreateBranch(branch)
	
		branch.Name = ""
		err := service.UpdateBranch(branch)
		assert.Error(t, err)
	})
	

	t.Run("DeleteBranch", func(t *testing.T) {
		branch := &coreModels.Branch{
			TenantID: 4,
			Name:     "To Be Deleted",
		}
		_ = service.CreateBranch(branch)

		err := service.DeleteBranch(branch.ID)
		assert.NoError(t, err)

		_, err = service.GetBranchByID(branch.ID)
		assert.Error(t, err) // Not found
	})

	

	


	t.Run("DeleteBranch_NotExist", func(t *testing.T) {
		err := service.DeleteBranch(9999999)
		assert.Error(t, err)
	})

	

	

	
	
	
}
