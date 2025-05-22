package coreTest

import (
	"testing"	
	"time"
	"errors"	

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	coreModels "myapp/modules/core/models"
	coreServices "myapp/modules/core/services"
)

func setupTestDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    assert.NoError(t, err)

    err = db.AutoMigrate(&coreModels.Tenant{},&coreModels.User{},)
    assert.NoError(t, err)

    tenants := []coreModels.Tenant{
        {ID: 1, Name: "Tenant One", Domain: "one.local", IsActive: true},
        {ID: 8, Name: "Tenant Eight", Domain: "eight.local", IsActive: true},
    }
    for _, tn := range tenants {
        assert.NoError(t, db.Create(&tn).Error)
    }
    err = db.AutoMigrate(&coreModels.Branch{})
    assert.NoError(t, err)
	assert.NoError(t, db.Create(&coreModels.Tenant{ID: 2, Name: "T1", Domain: "t1.local"}).Error)

    return db
}


func TestCreateBranch(t *testing.T) {
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
	
}

func TestGetAllBranches(t *testing.T) {
	db := setupTestDB(t)
	svc := coreServices.NewBranchService(db)

	t.Run("TenantNotFound", func(t *testing.T) {
		_, err := svc.GetAllBranches(999)
		assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
	})

	t.Run("NoBranches", func(t *testing.T) {
		// Tenant 2 has no branches
		brs, err := svc.GetAllBranches(8)
		assert.NoError(t, err)
		assert.NotNil(t, brs)
		assert.Len(t, brs, 0)
	})

	t.Run("MultipleBranchesOrdered", func(t *testing.T) {
		// Create three branches for tenant 1 with distinct CreatedAt
		now := time.Now()
		b1 := coreModels.Branch{TenantID: 1, Name: "B1", CreatedAt: now.Add(-2 * time.Hour)}
		b2 := coreModels.Branch{TenantID: 1, Name: "B2", CreatedAt: now.Add(-1 * time.Hour)}
		b3 := coreModels.Branch{TenantID: 1, Name: "B3", CreatedAt: now}
		assert.NoError(t, db.Create(&b1).Error)
		assert.NoError(t, db.Create(&b2).Error)
		assert.NoError(t, db.Create(&b3).Error)

		brs, err := svc.GetAllBranches(1)
		assert.NoError(t, err)
		// Expect order: newest first: B3, B2, B1
		assert.Len(t, brs, 3)
		assert.Equal(t, "B3", brs[0].Name)
		assert.Equal(t, "B2", brs[1].Name)
		assert.Equal(t, "B1", brs[2].Name)
	})
}

func TestGetBranchByID(t *testing.T) {
	db := setupTestDB(t)
	svc := coreServices.NewBranchService(db)

	t.Run("InvalidID", func(t *testing.T) {
		// id = 0 should error ErrInvalidID
		branch, err := svc.GetBranchByID(0)
		assert.Nil(t, branch)
		assert.ErrorIs(t, err, coreServices.ErrInvalidID)
	})

	t.Run("NotFound", func(t *testing.T) {
		// Tenant seeded but no branch with ID=42
		branch, err := svc.GetBranchByID(42)
		assert.Nil(t, branch)
		assert.ErrorIs(t, err, coreServices.ErrBranchNotFound)
	})

	t.Run("Success", func(t *testing.T) {
		// Create a branch
		b := coreModels.Branch{TenantID: 1, Name: "HQ"}
		assert.NoError(t, db.Create(&b).Error)
		// Retrieve it
		got, err := svc.GetBranchByID(b.ID)
		assert.NoError(t, err)
		assert.NotNil(t, got)
		assert.Equal(t, b.ID, got.ID)
		assert.Equal(t, "HQ", got.Name)
	})

	t.Run("SoftDeleted", func(t *testing.T) {
		// Create then soft-delete
		b := coreModels.Branch{TenantID: 1, Name: "Temp"}
		assert.NoError(t, db.Create(&b).Error)
		// Soft delete
		assert.NoError(t, db.Delete(&b).Error)
		// Attempt fetch
		branch, err := svc.GetBranchByID(b.ID)
		assert.Nil(t, branch)
		assert.ErrorIs(t, err, coreServices.ErrBranchNotFound)
	})
}

func TestUpdateBranch(t *testing.T) {
    t.Run("InvalidID", func(t *testing.T) {
        db := setupTestDB(t)
        svc := coreServices.NewBranchService(db)

        err := svc.UpdateBranch(&coreModels.Branch{ID: 0, Name: "X"})
        assert.ErrorIs(t, err, coreServices.ErrInvalidID)
    })

    t.Run("NameRequired", func(t *testing.T) {
        db := setupTestDB(t)
        svc := coreServices.NewBranchService(db)

        // even if ID non-zero, blank name after trim
        err := svc.UpdateBranch(&coreModels.Branch{ID: 1, Name: "   "})
        assert.ErrorIs(t, err, coreServices.ErrNameRequired)
    })

    t.Run("NotFound", func(t *testing.T) {
        db := setupTestDB(t)
        svc := coreServices.NewBranchService(db)

        // no branch with ID=42
        err := svc.UpdateBranch(&coreModels.Branch{ID: 42, Name: "New"})
        assert.ErrorIs(t, err, coreServices.ErrBranchNotFound)
    })

    t.Run("Success", func(t *testing.T) {
        db := setupTestDB(t)

        // create a branch to update
        original := coreModels.Branch{ID: 5, TenantID: 1, Name: "OldName", CreatedAt: time.Now()}
        assert.NoError(t, db.Create(&original).Error)

        svc := coreServices.NewBranchService(db)
        toUpdate := &coreModels.Branch{ID: 5, Name: " NewName "}
        err := svc.UpdateBranch(toUpdate)
        assert.NoError(t, err)

        // fetch back and verify
        var updated coreModels.Branch
        assert.NoError(t, db.First(&updated, 5).Error)
        assert.Equal(t, "NewName", updated.Name)
    })

    t.Run("DBErrorOnFetch", func(t *testing.T) {
        // simulate fetch error by closing the DB connection
        db := setupTestDB(t)
        // forcibly close underlying connection:
        sqlDB, _ := db.DB()
        sqlDB.Close()

        svc := coreServices.NewBranchService(db)
        err := svc.UpdateBranch(&coreModels.Branch{ID: 1, Name: "X"})
        assert.Error(t, err)
        assert.False(t, errors.Is(err, coreServices.ErrBranchNotFound),
            "should wrap underlying DB error, not ErrBranchNotFound")
    })

    t.Run("DBErrorOnUpdate", func(t *testing.T) {
        // Use a GORM callback to inject an error on Update
        db := setupTestDB(t)
        assert.NoError(t, db.Create(&coreModels.Branch{ID: 7, Name: "A"}).Error)

        db.Callback().Update().Replace("gorm:before_update", func(db *gorm.DB) {
            db.AddError(errors.New("update failed"))
        })

        svc := coreServices.NewBranchService(db)
        err := svc.UpdateBranch(&coreModels.Branch{ID: 7, Name: "New"})
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "error updating branch 7")
    })
}

func TestDeleteBranch(t *testing.T) {
    db := setupTestDB(t)
    svc := coreServices.NewBranchService(db)

    t.Run("InvalidID", func(t *testing.T) {
        err := svc.DeleteBranch(0)
        assert.ErrorIs(t, err, coreServices.ErrInvalidID)
    })

    t.Run("NotFound", func(t *testing.T) {
        // no branch with ID 99
        err := svc.DeleteBranch(99)
        assert.ErrorIs(t, err, coreServices.ErrBranchNotFound)
    })

    t.Run("InUse", func(t *testing.T) {
        // create branch + an attached user
        br := coreModels.Branch{TenantID: 1, Name: "InUse"}
        assert.NoError(t, db.Create(&br).Error)
        user := coreModels.User{Username: "u1", BranchID: &br.ID}
        assert.NoError(t, db.Create(&user).Error)

        err := svc.DeleteBranch(br.ID)
        assert.ErrorIs(t, err, coreServices.ErrBranchInUse)
    })

    t.Run("Success_SoftDelete", func(t *testing.T) {
        // create branch without users
        br := coreModels.Branch{TenantID: 1, Name: "ToDelete", CreatedAt: time.Now()}
        assert.NoError(t, db.Create(&br).Error)

        err := svc.DeleteBranch(br.ID)
        assert.NoError(t, err)

        // after soft-delete, deleted_at should be set
        var fetched coreModels.Branch
        lookupErr := db.First(&fetched, br.ID).Error
        assert.ErrorIs(t, lookupErr, gorm.ErrRecordNotFound)

        // but raw query (unscoped) should find it
        var raw coreModels.Branch
        assert.NoError(t, db.Unscoped().First(&raw, br.ID).Error)
        assert.NotNil(t, raw.DeletedAt)
    })
}

func TestGetBranchesByTenantID(t *testing.T) {
    t.Run("InvalidID", func(t *testing.T) {
        db := setupTestDB(t)
        svc := coreServices.NewBranchService(db)

        _, err := svc.GetBranchesByTenantID(0)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantID)
    })

    t.Run("TenantNotFound", func(t *testing.T) {
        db := setupTestDB(t)
        svc := coreServices.NewBranchService(db)

        _, err := svc.GetBranchesByTenantID(999)
        assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
    })

    t.Run("NoBranches", func(t *testing.T) {
        db := setupTestDB(t)
        svc := coreServices.NewBranchService(db)

        brs, err := svc.GetBranchesByTenantID(8) // tenant 8 ไม่มี branch
        assert.NoError(t, err)
        assert.Len(t, brs, 0)
    })

    t.Run("MultipleBranchesOrderedByCreatedAt", func(t *testing.T) {
        db := setupTestDB(t)

        now := time.Now()
        b1 := coreModels.Branch{TenantID: 1, Name: "Old", CreatedAt: now.Add(-2 * time.Hour)}
        b2 := coreModels.Branch{TenantID: 1, Name: "Mid", CreatedAt: now.Add(-1 * time.Hour)}
        b3 := coreModels.Branch{TenantID: 1, Name: "New", CreatedAt: now}
        assert.NoError(t, db.Create(&b1).Error)
        assert.NoError(t, db.Create(&b2).Error)
        assert.NoError(t, db.Create(&b3).Error)

        svc := coreServices.NewBranchService(db)
        brs, err := svc.GetBranchesByTenantID(1)
        assert.NoError(t, err)
        assert.Len(t, brs, 3)
        assert.Equal(t, "New", brs[0].Name)
        assert.Equal(t, "Mid", brs[1].Name)
        assert.Equal(t, "Old", brs[2].Name)
    })

    t.Run("SkipSoftDeleted", func(t *testing.T) {
        db := setupTestDB(t)

        b := coreModels.Branch{TenantID: 2, Name: "ToDelete", CreatedAt: time.Now()}
        assert.NoError(t, db.Create(&b).Error)
        // soft-delete
        assert.NoError(t, db.Delete(&b).Error)

        svc := coreServices.NewBranchService(db)
        brs, err := svc.GetBranchesByTenantID(2)
        assert.NoError(t, err)
        assert.Len(t, brs, 0)
    })
}


