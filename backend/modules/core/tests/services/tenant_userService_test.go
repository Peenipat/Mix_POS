package coreTest

import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    coreModels "myapp/modules/core/models"
    coreServices "myapp/modules/core/services"
)

func setupTenantUserDB(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)
    require.NoError(t, db.AutoMigrate(
        &coreModels.Tenant{},
        &coreModels.User{},
        &coreModels.TenantUser{},
    ))
    return db
}

func TestAddUserToTenant(t *testing.T) {
    ctx := context.Background()

    t.Run("InvalidTenantID", func(t *testing.T) {
        svc := coreServices.NewTenantUserService(setupTenantUserDB(t))
        err := svc.AddUserToTenant(ctx, 0, 1)
        assert.ErrorIs(t, err, coreServices.ErrInvalidTenantID)
    })

    t.Run("InvalidUserID", func(t *testing.T) {
        svc := coreServices.NewTenantUserService(setupTenantUserDB(t))
        err := svc.AddUserToTenant(ctx, 1, 0)
        assert.ErrorIs(t, err, coreServices.ErrInvalidUserID)
    })

    t.Run("TenantNotFound", func(t *testing.T) {
        db := setupTenantUserDB(t)
        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 42, 1)
        assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
    })

    t.Run("UserNotFound", func(t *testing.T) {
        db := setupTenantUserDB(t)
        // create tenant only
        require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 1, 99)
        assert.ErrorIs(t, err, coreServices.ErrUserNotFound)
    })

    t.Run("AlreadyAssigned", func(t *testing.T) {
        db := setupTenantUserDB(t)
        tnt := coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}
        usr := coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}
        require.NoError(t, db.Create(&tnt).Error)
        require.NoError(t, db.Create(&usr).Error)
        // pre-create assignment
        require.NoError(t, db.Create(&coreModels.TenantUser{TenantID: 1, UserID: 2}).Error)

        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 1, 2)
        assert.ErrorIs(t, err, coreServices.ErrUserAlreadyAssigned)
    })

    t.Run("DBErrorOnFetchTenant", func(t *testing.T) {
        db := setupTenantUserDB(t)
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close()) // simulate error
        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 1, 1)
        require.Error(t, err)
    })

    t.Run("DBErrorOnFetchUser", func(t *testing.T) {
        db := setupTenantUserDB(t)
        // create valid tenant
        require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close()) // simulate error before user fetch
        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 1, 1)
        require.Error(t, err)
    })

    t.Run("DBErrorOnCheckAssignment", func(t *testing.T) {
        db := setupTenantUserDB(t)
        // create tenant and user
        require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
        require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close()) // simulate error before checking assignment
        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 1, 2)
        require.Error(t, err)
    })

    t.Run("DBErrorOnCreate", func(t *testing.T) {
        db := setupTenantUserDB(t)
        require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
        require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)
        svcCloseDB := coreServices.NewTenantUserService(db)
        sqlDB, _ := db.DB()
        require.NoError(t, sqlDB.Close()) // simulate error after checks
        err := svcCloseDB.AddUserToTenant(ctx, 1, 2)
        require.Error(t, err)
    })

    t.Run("Success", func(t *testing.T) {
        db := setupTenantUserDB(t)
        require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
        require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)

        svc := coreServices.NewTenantUserService(db)
        err := svc.AddUserToTenant(ctx, 1, 2)
        require.NoError(t, err)

        // confirm in DB
        var tu coreModels.TenantUser
        err = db.Where("tenant_id = ? AND user_id = ?", 1, 2).First(&tu).Error
        require.NoError(t, err)
        assert.Equal(t, uint(1), tu.TenantID)
        assert.Equal(t, uint(2), tu.UserID)
    })
}


func TestRemoveUserFromTenant(t *testing.T) {
	ctx := context.Background()

	t.Run("InvalidTenantID", func(t *testing.T) {
		svc := coreServices.NewTenantUserService(setupTenantUserDB(t))
		err := svc.RemoveUserFromTenant(ctx, 0, 1)
		assert.ErrorIs(t, err, coreServices.ErrInvalidTenantID)
	})

	t.Run("InvalidUserID", func(t *testing.T) {
		svc := coreServices.NewTenantUserService(setupTenantUserDB(t))
		err := svc.RemoveUserFromTenant(ctx, 1, 0)
		assert.ErrorIs(t, err, coreServices.ErrInvalidUserID)
	})

	t.Run("TenantNotFound", func(t *testing.T) {
		db := setupTenantUserDB(t)
		svc := coreServices.NewTenantUserService(db)
		err := svc.RemoveUserFromTenant(ctx, 42, 1)
		assert.ErrorIs(t, err, coreServices.ErrTenantNotFound)
	})

	t.Run("UserNotFound", func(t *testing.T) {
		db := setupTenantUserDB(t)
		require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
		svc := coreServices.NewTenantUserService(db)
		err := svc.RemoveUserFromTenant(ctx, 1, 99)
		assert.ErrorIs(t, err, coreServices.ErrUserNotFound)
	})

	t.Run("UserNotAssigned", func(t *testing.T) {
		db := setupTenantUserDB(t)
		require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
		require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)
		svc := coreServices.NewTenantUserService(db)
		err := svc.RemoveUserFromTenant(ctx, 1, 2)
		assert.ErrorIs(t, err, coreServices.ErrUserNotAssigned)
	})

	t.Run("DBErrorOnFetchTenant", func(t *testing.T) {
		db := setupTenantUserDB(t)
		svc := coreServices.NewTenantUserService(db)
		sqlDB, _ := db.DB()
		require.NoError(t, sqlDB.Close())
		err := svc.RemoveUserFromTenant(ctx, 1, 1)
		require.Error(t, err)
	})

	t.Run("DBErrorOnFetchUser", func(t *testing.T) {
		db := setupTenantUserDB(t)
		require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
		svc := coreServices.NewTenantUserService(db)
		sqlDB, _ := db.DB()
		require.NoError(t, sqlDB.Close())
		err := svc.RemoveUserFromTenant(ctx, 1, 1)
		require.Error(t, err)
	})

	t.Run("DBErrorOnFetchAssignment", func(t *testing.T) {
		db := setupTenantUserDB(t)
		require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
		require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)
		svc := coreServices.NewTenantUserService(db)
		sqlDB, _ := db.DB()
		require.NoError(t, sqlDB.Close())
		err := svc.RemoveUserFromTenant(ctx, 1, 2)
		require.Error(t, err)
	})

	t.Run("DBErrorOnDelete", func(t *testing.T) {
		db := setupTenantUserDB(t)
		require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
		require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)
		require.NoError(t, db.Create(&coreModels.TenantUser{TenantID: 1, UserID: 2}).Error)
		svcClose := coreServices.NewTenantUserService(db)
		sqlDB, _ := db.DB()
		require.NoError(t, sqlDB.Close())
		err := svcClose.RemoveUserFromTenant(ctx, 1, 2)
		require.Error(t, err)
	})

	t.Run("Success", func(t *testing.T) {
		db := setupTenantUserDB(t)
		require.NoError(t, db.Create(&coreModels.Tenant{ID: 1, Name: "T", Domain: "d"}).Error)
		require.NoError(t, db.Create(&coreModels.User{ID: 2, Username: "u", Email: "u@example.com"}).Error)
		require.NoError(t, db.Create(&coreModels.TenantUser{TenantID: 1, UserID: 2}).Error)

		svc := coreServices.NewTenantUserService(db)
		err := svc.RemoveUserFromTenant(ctx, 1, 2)
		require.NoError(t, err)

		var tu coreModels.TenantUser
		err = db.Where("tenant_id = ? AND user_id = ?", 1, 2).First(&tu).Error
		// record should be gone (ErrRecordNotFound)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})
}










