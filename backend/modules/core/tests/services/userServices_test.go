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
    require.NoError(t, db.AutoMigrate(&coreModels.User{}))
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