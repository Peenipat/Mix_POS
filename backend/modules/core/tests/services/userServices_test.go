package coreTest
import (
    "context"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "golang.org/x/crypto/bcrypt"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"

    coreModels "myapp/modules/core/models"
    coreServices "myapp/modules/core/services"
)

func setupUserDBTest(t *testing.T) *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // Migrate only the User model for these tests
    require.NoError(t, db.AutoMigrate(&coreModels.Tenant{}))
    require.NoError(t, db.AutoMigrate(&coreModels.Branch{}))
    require.NoError(t, db.AutoMigrate(&coreModels.User{}))
    require.NoError(t, db.AutoMigrate(&coreModels.TenantUser{}))
    return db
}

func TestChangePassword(t *testing.T) {
    ctx := context.Background()

    t.Run("UserNotFound", func(t *testing.T) {
        db := setupUserDBTest(t)
        svc := coreServices.NewUserService(db)

        err := svc.ChangePassword(ctx, 999, "old", "new")
        assert.ErrorIs(t, err, coreServices.ErrUserNotFound)
    })

    t.Run("InvalidOldPassword", func(t *testing.T) {
        db := setupUserDBTest(t)
        svc := coreServices.NewUserService(db)

        // Seed a user with a known password
        hashed, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
        u := coreModels.User{Username: "u", Email: "e@e", Password: string(hashed)}
        require.NoError(t, db.Create(&u).Error)

        err := svc.ChangePassword(ctx, u.ID, "wrong-old", "newpass")
        assert.ErrorIs(t, err, coreServices.ErrInvalidOldPassword)
    })

    t.Run("UpdateError", func(t *testing.T) {
        db := setupUserDBTest(t)
        sqlDB, _ := db.DB()
        // Close the underlying *sql.DB to force an update error
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewUserService(db)
        // Create a dummy user so user lookup succeeds
        db.Exec("INSERT INTO users (id, username, email, password, created_at, updated_at) VALUES (?, ?, ?, ?, datetime('now'), datetime('now'))",
            1, "u", "e@e", "$2a$10$7EqJtq98hPqEX7fNZaFWoO") // bcrypt hash for "password"

        err := svc.ChangePassword(ctx, 1, "password", "newpass")
        require.Error(t, err)
        assert.NotEqual(t, coreServices.ErrUserNotFound, err)
    })

    t.Run("Success", func(t *testing.T) {
        db := setupUserDBTest(t)
        svc := coreServices.NewUserService(db)

        // Seed a user with a known password
        hashed, _ := bcrypt.GenerateFromPassword([]byte("oldpass"), bcrypt.DefaultCost)
        u := coreModels.User{Username: "u", Email: "e@e", Password: string(hashed)}
        require.NoError(t, db.Create(&u).Error)

        err := svc.ChangePassword(ctx, u.ID, "oldpass", "brandnew")
        require.NoError(t, err)

        // Reload and verify the password has been updated
        var updated coreModels.User
        require.NoError(t, db.First(&updated, u.ID).Error)

        // Ensure the new hash matches the new password
        assert.NoError(t, bcrypt.CompareHashAndPassword(
            []byte(updated.Password),
            []byte("brandnew"),
        ))
    })
}

func TestUserService_Me(t *testing.T) {
    ctx := context.Background()

    t.Run("UserNotFound", func(t *testing.T) {
        db := setupUserDBTest(t)
        svc := coreServices.NewUserService(db)

        dto, err := svc.Me(ctx, 999)
        require.NoError(t, err)
        assert.Nil(t, dto)
    })

    t.Run("DBError", func(t *testing.T) {
        db := setupUserDBTest(t)
        sqlDB, err := db.DB()
        require.NoError(t, err)
        require.NoError(t, sqlDB.Close())

        svc := coreServices.NewUserService(db)
        dto, err := svc.Me(ctx, 1)
        require.Error(t, err)
        assert.Nil(t, dto)
    })

    t.Run("Success_NoAffiliations", func(t *testing.T) {
        db := setupUserDBTest(t)
        svc := coreServices.NewUserService(db)

        // Seed user ไม่มี branch และ tenant_users
        u := coreModels.User{
            Username: "alice",
            Email:    "alice@example.com",
            Password: "irrelevant",
        }
        require.NoError(t, db.Create(&u).Error)

        dto, err := svc.Me(ctx, u.ID)
        require.NoError(t, err)
        require.NotNil(t, dto)

        assert.Equal(t, u.ID, dto.ID)
        assert.Equal(t, u.Username, dto.Username)
        assert.Equal(t, u.Email, dto.Email)
        assert.Nil(t, dto.BranchID)
        assert.Len(t, dto.TenantIDs, 0)
    })

    t.Run("Success_WithBranchAndTenants", func(t *testing.T) {
        db := setupUserDBTest(t)
        svc := coreServices.NewUserService(db)

        // สร้าง Tenant สองอัน
        tenant1 := coreModels.Tenant{Name: "T1", Domain: "t1.local", IsActive: true}
        tenant2 := coreModels.Tenant{Name: "T2", Domain: "t2.local", IsActive: true}
        require.NoError(t, db.Create(&tenant1).Error)
        require.NoError(t, db.Create(&tenant2).Error)

        // สร้าง Branch
        branch := coreModels.Branch{Name: "B1", TenantID: tenant1.ID}
        require.NoError(t, db.Create(&branch).Error)

        // สร้าง User ผูกกับ Branch
        u := coreModels.User{
            Username: "bob",
            Email:    "bob@example.com",
            Password: "irrelevant",
            BranchID: &branch.ID,
        }
        require.NoError(t, db.Create(&u).Error)

        // สร้าง TenantUser ผูก Bob กับ 2 Tenant
        tu1 := coreModels.TenantUser{UserID: u.ID, TenantID: tenant1.ID}
        tu2 := coreModels.TenantUser{UserID: u.ID, TenantID: tenant2.ID}
        require.NoError(t, db.Create(&tu1).Error)
        require.NoError(t, db.Create(&tu2).Error)

        dto, err := svc.Me(ctx, u.ID)
        require.NoError(t, err)
        require.NotNil(t, dto)

        assert.Equal(t, u.ID, dto.ID)
        assert.Equal(t, u.Username, dto.Username)
        assert.Equal(t, u.Email, dto.Email)
        require.NotNil(t, dto.BranchID)
        assert.Equal(t, branch.ID, *dto.BranchID)
        assert.ElementsMatch(t, []uint{tenant1.ID, tenant2.ID}, dto.TenantIDs)
    })
}