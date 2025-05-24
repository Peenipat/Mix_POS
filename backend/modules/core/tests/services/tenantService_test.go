package coreTest


import (
    "context"
    "testing"
	"errors"
    "time"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    coreModels "myapp/modules/core/models"
    corePorts "myapp/modules/core/port"
    coreServices "myapp/modules/core/services"
)

func setupTenantDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    require.NoError(t, db.AutoMigrate(
        &coreModels.Tenant{},
        &coreModels.TenantUser{},
        &coreModels.Branch{},
    ))
    return db
}
func ptr(s string) *string { return &s }

func TestCreateTenant(t *testing.T) {
    ctx := context.Background()

    t.Run("InvalidInput_EmptyName", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        input := corePorts.CreateTenantInput{Name: "", Domain: "example.com"}
        tenant, err := svc.CreateTenant(ctx, input)
        assert.Nil(t, tenant)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantInput)
    })

    t.Run("InvalidInput_EmptyDomain", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        input := corePorts.CreateTenantInput{Name: "Acme", Domain: " "}
        tenant, err := svc.CreateTenant(ctx, input)
        assert.Nil(t, tenant)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantInput)
    })

    t.Run("DomainTaken", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        // seed existing tenant
        existing := coreModels.Tenant{Name: "Old", Domain: "acme.local"}
        require.NoError(t, db.Create(&existing).Error)

        input := corePorts.CreateTenantInput{Name: "New", Domain: "acme.local"}
        tenant, err := svc.CreateTenant(ctx, input)
        assert.Nil(t, tenant)
        assert.ErrorIs(t, err, coreServices.ErrDomainTaken)
    })

    t.Run("UniquenessCheckError", func(t *testing.T) {
        db := setupTenantDB(t)
        sqlDB, err := db.DB()
        require.NoError(t, err)
        require.NoError(t, sqlDB.Close()) // force subsequent queries to error

        svc := coreServices.NewTenantService(db)
        input := corePorts.CreateTenantInput{Name: "Acme", Domain: "acme.local"}
        tenant, err := svc.CreateTenant(ctx, input)
        assert.Nil(t, tenant)
        require.Error(t, err)
        assert.NotErrorIs(t, err, coreServices.ErrInvalidTenantInput)
        assert.NotErrorIs(t, err, coreServices.ErrDomainTaken)
    })

    t.Run("CreateError", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        // wrap Create to fail: drop table to simulate error
        require.NoError(t, db.Migrator().DropTable(&coreModels.Tenant{}))

        input := corePorts.CreateTenantInput{Name: "Acme", Domain: "acme.local"}
        tenant, err := svc.CreateTenant(ctx, input)
        assert.Nil(t, tenant)
        require.Error(t, err)
    })

    t.Run("Success", func(t *testing.T) {
		db := setupTenantDB(t)
		svc := coreServices.NewTenantService(db)
	
		// 1) Record time window around creation
		start := time.Now()
		input := corePorts.CreateTenantInput{Name: "Acme Corp", Domain: "acme.local"}
		tenant, err := svc.CreateTenant(ctx, input)
		end := time.Now()
	
		require.NoError(t, err)
		require.NotNil(t, tenant)
	
		// 2) Basic field checks
		assert.Equal(t, "Acme Corp", tenant.Name)
		assert.Equal(t, "acme.local", tenant.Domain)
		assert.True(t, tenant.IsActive)
	
		// 3) Timestamp falls between start and end
		assert.False(t, tenant.CreatedAt.Before(start), "CreatedAt is before start")
		assert.False(t, tenant.CreatedAt.After(end),   "CreatedAt is after end")
		assert.False(t, tenant.UpdatedAt.Before(start), "UpdatedAt is before start")
		assert.False(t, tenant.UpdatedAt.After(end),   "UpdatedAt is after end")
	
		// 4) Confirm persisted in DB
		var fetched coreModels.Tenant
		require.NoError(t, db.First(&fetched, tenant.ID).Error)
		assert.Equal(t, tenant.ID,     fetched.ID)
		assert.Equal(t, tenant.Domain, fetched.Domain)
	})
	
}

func TestGetTenantByID(t *testing.T) {
    ctx := context.Background()

    t.Run("InvalidID", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        tenant, err := svc.GetTenantByID(ctx, 0)
        assert.Nil(t, tenant)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantID)
    })

    t.Run("NotFound", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        tenant, err := svc.GetTenantByID(ctx, 123)
        assert.Nil(t, tenant)
        assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
    })

    t.Run("DBError", func(t *testing.T) {
        db := setupTenantDB(t)
        // Force a DB error by closing the underlying sql.DB
        sqlDB, err := db.DB()
        require.NoError(t, err)
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewTenantService(db)
        tenant, err := svc.GetTenantByID(ctx, 1)
        assert.Nil(t, tenant)
        require.Error(t, err)
        // It should not be ErrInvalidTenantID or ErrTenantNotFound
        assert.False(t, errors.Is(err, coreServices.ErrInvalidTenantID))
        assert.False(t, errors.Is(err, coreServices.ErrTenantNotFound))
    })

    t.Run("Success", func(t *testing.T) {
        db := setupTenantDB(t)
        // Seed a tenant
        created := coreModels.Tenant{Name: "Test", Domain: "test.local", IsActive: true}
        require.NoError(t, db.Create(&created).Error)

        svc := coreServices.NewTenantService(db)
        tenant, err := svc.GetTenantByID(ctx, created.ID)
        require.NoError(t, err)
        require.NotNil(t, tenant)

        assert.Equal(t, created.ID, tenant.ID)
        assert.Equal(t, "Test", tenant.Name)
        assert.Equal(t, "test.local", tenant.Domain)
        assert.True(t, tenant.IsActive)
    })
}

func TestListTenants(t *testing.T) {
    ctx := context.Background()

    t.Run("Empty_All", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        tenants, err := svc.ListTenants(ctx, false)
        require.NoError(t, err)
        assert.NotNil(t, tenants)
        assert.Len(t, tenants, 0)
    })

    t.Run("Empty_ActiveOnly", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)

        tenants, err := svc.ListTenants(ctx, true)
        require.NoError(t, err)
        assert.NotNil(t, tenants)
        assert.Len(t, tenants, 0)
    })

    t.Run("Multiple_MixedActive", func(t *testing.T) {
        db := setupTenantDB(t)

        now := time.Now()
        t1 := coreModels.Tenant{
            Name:      "A",
            Domain:    "a.local",
            IsActive:  true,
            CreatedAt: now.Add(-3 * time.Hour),
            UpdatedAt: now.Add(-3 * time.Hour),
        }
        t2 := coreModels.Tenant{
            Name:      "B",
            Domain:    "b.local",
            IsActive:  false, // ต้องการ inactive
            CreatedAt: now.Add(-2 * time.Hour),
            UpdatedAt: now.Add(-2 * time.Hour),
        }
        t3 := coreModels.Tenant{
            Name:      "C",
            Domain:    "c.local",
            IsActive:  true,
            CreatedAt: now.Add(-1 * time.Hour),
            UpdatedAt: now.Add(-1 * time.Hour),
        }

        // 1) สร้าง t1 และ t3 ตามปกติ
        require.NoError(t, db.Create(&t1).Error)
        require.NoError(t, db.Create(&t3).Error)

        // 2) สร้าง t2 แล้วอัพเดตคอลัมน์ is_active ให้เป็น false
        require.NoError(t, db.Create(&t2).Error)
        require.NoError(t, db.Model(&t2).
            Select("is_active").
            Update("is_active", false).Error)

        svc := coreServices.NewTenantService(db)

        // onlyActive = true → ควรได้ C ก่อน A
        active, err := svc.ListTenants(ctx, true)
        require.NoError(t, err)
        require.Len(t, active, 2)
        assert.Equal(t, "C", active[0].Name)
        assert.Equal(t, "A", active[1].Name)

        // onlyActive = false → ควรได้ C, B, A
        all, err := svc.ListTenants(ctx, false)
        require.NoError(t, err)
        require.Len(t, all, 3)
        assert.Equal(t, "C", all[0].Name)
        assert.Equal(t, "B", all[1].Name)
        assert.Equal(t, "A", all[2].Name)
    })

    t.Run("SkipSoftDeleted", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "ToDelete", Domain: "d.local", IsActive: true}
        require.NoError(t, db.Create(&tenant).Error)
        // soft-delete
        require.NoError(t, db.Delete(&tenant).Error)

        svc := coreServices.NewTenantService(db)
        tenants, err := svc.ListTenants(ctx, false)
        require.NoError(t, err)
        assert.Len(t, tenants, 0)
    })

    t.Run("DBError", func(t *testing.T) {
        db := setupTenantDB(t)
        sqlDB, err := db.DB()
        require.NoError(t, err)
        require.NoError(t, sqlDB.Close()) // force error

        svc := coreServices.NewTenantService(db)
        tenants, err := svc.ListTenants(ctx, false)
        assert.Nil(t, tenants)
        require.Error(t, err)
        assert.Contains(t, err.Error(), "failed to fetch tenants")
    })
}

func TestUpdateTenant(t *testing.T) {
    ctx := context.Background()
    svc := coreServices.NewTenantService(setupTenantDB(t))

    t.Run("InvalidID", func(t *testing.T) {
        err := svc.UpdateTenant(ctx, corePorts.UpdateTenantInput{ID: 0})
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantID)
    })

    t.Run("EmptyName", func(t *testing.T) {
        input := corePorts.UpdateTenantInput{
            ID:   1,
            Name: ptr("   "),
        }
        err := svc.UpdateTenant(ctx, input)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantInput)
    })

    t.Run("EmptyDomain", func(t *testing.T) {
        input := corePorts.UpdateTenantInput{
            ID:     1,
            Domain: ptr(" "),
        }
        err := svc.UpdateTenant(ctx, input)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantInput)
    })

    t.Run("TenantNotFound", func(t *testing.T) {
        input := corePorts.UpdateTenantInput{
            ID:   99,
            Name: ptr("X"),
        }
        err := svc.UpdateTenant(ctx, input)
        assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
    })

    t.Run("DomainTaken", func(t *testing.T) {
        db := setupTenantDB(t)
        // seed two tenants
        t1 := coreModels.Tenant{Name: "A", Domain: "a.local"}
        t2 := coreModels.Tenant{Name: "B", Domain: "b.local"}
        require.NoError(t, db.Create(&t1).Error)
        require.NoError(t, db.Create(&t2).Error)

        svc := coreServices.NewTenantService(db)
        input := corePorts.UpdateTenantInput{
            ID:     t1.ID,
            Domain: ptr("b.local"),
        }
        err := svc.UpdateTenant(ctx, input)
        assert.ErrorIs(t, err, coreServices.ErrDomainTaken)
    })

    t.Run("DBErrorOnFetch", func(t *testing.T) {
        db := setupTenantDB(t)
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewTenantService(db)
        input := corePorts.UpdateTenantInput{ID: 1, Name: ptr("X")}
        err := svc.UpdateTenant(ctx, input)
        require.Error(t, err)
    })

    t.Run("DBErrorOnUpdate", func(t *testing.T) {
		db := setupTenantDB(t)
		// seed one tenant
		tenant := coreModels.Tenant{Name: "A", Domain: "a.local"}
		require.NoError(t, db.Create(&tenant).Error)
	
		// ปิด DB เพื่อจำลอง failure
		sqlDB, _ := db.DB()
		require.NoError(t, sqlDB.Close())
	
		svc := coreServices.NewTenantService(db)
		input := corePorts.UpdateTenantInput{
			ID:   tenant.ID,
			Name: ptr("NewName"),
		}
		err := svc.UpdateTenant(ctx, input)
		// แค่ตรวจว่ามี error เท่านั้น
		require.Error(t, err)
	})

    t.Run("NoChanges", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "A", Domain: "a.local"}
        require.NoError(t, db.Create(&tenant).Error)

        svc := coreServices.NewTenantService(db)
        input := corePorts.UpdateTenantInput{ID: tenant.ID}
        err := svc.UpdateTenant(ctx, input)
        assert.NoError(t, err)
    })

    t.Run("Success_AllFields", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "A", Domain: "a.local", IsActive: false}
        require.NoError(t, db.Create(&tenant).Error)

        svc := coreServices.NewTenantService(db)
        newName := "NewA"
        newDomain := "newa.local"
        newActive := true
        input := corePorts.UpdateTenantInput{
            ID:       tenant.ID,
            Name:     &newName,
            Domain:   &newDomain,
            IsActive: &newActive,
        }
        err := svc.UpdateTenant(ctx, input)
        require.NoError(t, err)

        var updated coreModels.Tenant
        require.NoError(t, db.First(&updated, tenant.ID).Error)
        assert.Equal(t, newName, updated.Name)
        assert.Equal(t, newDomain, updated.Domain)
        assert.Equal(t, newActive, updated.IsActive)
    })
}

func TestDeleteTenant(t *testing.T) {
    ctx := context.Background()

    t.Run("InvalidID", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, 0)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantID)
    })

    t.Run("NotFound", func(t *testing.T) {
        db := setupTenantDB(t)
        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, 123)
        assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
    })

    t.Run("InUseByUser", func(t *testing.T) {
        db := setupTenantDB(t)
        // create tenant + a TenantUser record
        tenant := coreModels.Tenant{Name: "T", Domain: "d"}
        require.NoError(t, db.Create(&tenant).Error)
        tu := coreModels.TenantUser{TenantID: tenant.ID, UserID: 42}
        require.NoError(t, db.Create(&tu).Error)

        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, tenant.ID)
        assert.ErrorIs(t, err, coreServices.ErrTenantInUse)
    })

    t.Run("InUseByBranch", func(t *testing.T) {
        db := setupTenantDB(t)
        // create tenant + a Branch record
        tenant := coreModels.Tenant{Name: "T", Domain: "d"}
        require.NoError(t, db.Create(&tenant).Error)
        branch := coreModels.Branch{TenantID: tenant.ID, Name: "B"}
        require.NoError(t, db.Create(&branch).Error)

        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, tenant.ID)
        assert.ErrorIs(t, err, coreServices.ErrTenantInUse)
    })

    t.Run("DBErrorOnFetch", func(t *testing.T) {
        db := setupTenantDB(t)
        // corrupt DB
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, 1)
        require.Error(t, err)
    })

    t.Run("DBErrorOnTenantUserCheck", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "T", Domain: "d"}
        require.NoError(t, db.Create(&tenant).Error)
        // simulate error by closing low-level
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, tenant.ID)
        require.Error(t, err)
    })

    t.Run("DBErrorOnBranchCheck", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "T", Domain: "d"}
        require.NoError(t, db.Create(&tenant).Error)
        // first check (TenantUser) passes; now corrupt DB before branch check
        require.NoError(t, db.Delete(&coreModels.TenantUser{TenantID: tenant.ID}).Error)
        // simulate closing
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, tenant.ID)
        require.Error(t, err)
    })

    t.Run("DBErrorOnDelete", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "T", Domain: "d"}
        require.NoError(t, db.Create(&tenant).Error)

        // ensure no dependencies
        require.NoError(t, db.Where("tenant_id = ?", tenant.ID).Delete(&coreModels.TenantUser{}).Error)
        require.NoError(t, db.Where("tenant_id = ?", tenant.ID).Delete(&coreModels.Branch{}).Error)

        // corrupt DB before delete
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewTenantService(db)
		err := svc.DeleteTenant(ctx, tenant.ID)
		require.Error(t, err)
    })

    t.Run("Success", func(t *testing.T) {
        db := setupTenantDB(t)
        tenant := coreModels.Tenant{Name: "T", Domain: "d"}
        require.NoError(t, db.Create(&tenant).Error)

        // no TenantUser, no Branch
        svc := coreServices.NewTenantService(db)
        err := svc.DeleteTenant(ctx, tenant.ID)
        require.NoError(t, err)

		var fetched coreModels.Tenant
		err = db.First(&fetched, tenant.ID).Error
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
    })
}
