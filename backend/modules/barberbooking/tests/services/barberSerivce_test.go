package barberbookingServiceTest
import (
	"context"
	"testing"
	"errors"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	coreModels "myapp/modules/core/models"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingServices "myapp/modules/barberbooking/services"
)

func setupTestBarberDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	assert.NoError(t, err)
	assert.NoError(t, db.AutoMigrate(&barberBookingModels.Barber{},coreModels.Branch{},))
	return db
}

func TestBarberService_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestBarberDB(t)
	svc := barberBookingServices.NewBarberService(db)

	// Create
	t.Run("CreateBarber", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 1, UserID: 100}
		err := svc.CreateBarber(ctx, barber)
		assert.NoError(t, err)
		assert.NotZero(t, barber.ID)
	})

	// Get
	t.Run("GetBarberByID", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 2, UserID: 200}
		_ = svc.CreateBarber(ctx, barber)

		found, err := svc.GetBarberByID(ctx, barber.ID)
		assert.NoError(t, err)
		assert.Equal(t, barber.UserID, found.UserID)
	})

	// Update
	// t.Run("UpdateBarber", func(t *testing.T) {
	// 	barber := &barberBookingModels.Barber{BranchID: 4, UserID: 400}
	// 	_ = svc.CreateBarber(ctx, barber)

	// 	updates := &barberBookingModels.Barber{BranchID: 5}
	// 	updated, err := svc.UpdateBarber(ctx, barber.ID, updates)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, uint(5), updated.BranchID)
	// })

	// List
	t.Run("ListBarbers", func(t *testing.T) {
		_ = svc.CreateBarber(ctx, &barberBookingModels.Barber{BranchID: 3, UserID: 300})
		_ = svc.CreateBarber(ctx, &barberBookingModels.Barber{BranchID: 3, UserID: 301})

		branchID := uint(3)
		barbers, err := svc.ListBarbersByBranch(ctx, &branchID) // ⬅️ แก้ตรงนี้ให้รับ *uint
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(barbers), 2)
	})

	// Delete
	t.Run("DeleteBarber", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 6, UserID: 600}
		_ = svc.CreateBarber(ctx, barber)

		err := svc.DeleteBarber(ctx, barber.ID)
		assert.NoError(t, err)

		_, err = svc.GetBarberByID(ctx, barber.ID)
		assert.Error(t, err)
	})

	t.Run("CreateBarber_DuplicateUserID", func(t *testing.T) {
		userID := uint(700)
		barber1 := &barberBookingModels.Barber{BranchID: 1, UserID: userID}
		barber2 := &barberBookingModels.Barber{BranchID: 2, UserID: userID}

		err1 := svc.CreateBarber(ctx, barber1)
		assert.NoError(t, err1)

		err2 := svc.CreateBarber(ctx, barber2)
		assert.Error(t, err2, "should fail due to unique UserID constraint")
	})

	t.Run("CreateBarber_MissingFields", func(t *testing.T) {
		missingUser := &barberBookingModels.Barber{BranchID: 1}
		missingBranch := &barberBookingModels.Barber{UserID: 800}

		err1 := svc.CreateBarber(ctx, missingUser)
		assert.Error(t, err1, "should fail due to missing UserID")

		err2 := svc.CreateBarber(ctx, missingBranch)
		assert.Error(t, err2, "should fail due to missing BranchID")
	})

	t.Run("CreateBarber_MultipleInSameBranch", func(t *testing.T) {
		branchID := uint(10)
		barbers := []barberBookingModels.Barber{
			{BranchID: branchID, UserID: 901},
			{BranchID: branchID, UserID: 902},
			{BranchID: branchID, UserID: 903},
		}

		for _, b := range barbers {
			err := svc.CreateBarber(ctx, &b)
			assert.NoError(t, err)
		}

		result, err := svc.ListBarbersByBranch(ctx, &branchID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(result), 3)
	})

	t.Run("GetBarberByID_NotFound", func(t *testing.T) {
		_, err := svc.GetBarberByID(ctx, 9999) // ID ที่ไม่มีจริง
		assert.Error(t, err, "ควร error เพราะไม่พบ barber")
	})

	t.Run("GetBarberByID_SoftDeleted", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 10, UserID: 1000}
		_ = svc.CreateBarber(ctx, barber)
		_ = svc.DeleteBarber(ctx, barber.ID)

		_, err := svc.GetBarberByID(ctx, barber.ID)
		assert.Error(t, err, "ควร error เพราะถูก soft-delete ไปแล้ว")
	})

	// t.Run("UpdateBarber_NotFound", func(t *testing.T) {
	// 	updates := &barberBookingModels.Barber{BranchID: 20, UserID: 2000}
	// 	_, err := svc.UpdateBarber(ctx, 9999, updates)
	// 	assert.Error(t, err, "ควร error เพราะไม่พบ barber ID นี้")
	// })

	// t.Run("UpdateBarber_OnlyBranchID", func(t *testing.T) {
	// 	barber := &barberBookingModels.Barber{BranchID: 30, UserID: 3000}
	// 	_ = svc.CreateBarber(ctx, barber)

	// 	updates := &barberBookingModels.Barber{BranchID: 31, UserID: 3000} // ไม่เปลี่ยน user_id
	// 	updated, err := svc.UpdateBarber(ctx, barber.ID, updates)
	// 	assert.NoError(t, err)
	// 	assert.Equal(t, uint(31), updated.BranchID)
	// 	assert.Equal(t, uint(3000), updated.UserID)
	// })

	// t.Run("UpdateBarber_DuplicateUserID", func(t *testing.T) {
	// 	// สร้าง barber คนแรก
	// 	_ = svc.CreateBarber(ctx, &barberBookingModels.Barber{BranchID: 40, UserID: 4000})

	// 	// สร้าง barber คนที่สอง
	// 	barber2 := &barberBookingModels.Barber{BranchID: 41, UserID: 4001}
	// 	_ = svc.CreateBarber(ctx, barber2)

	// 	// พยายามอัปเดต user_id ให้ซ้ำ
	// 	updates := &barberBookingModels.Barber{BranchID: 41, UserID: 4000}
	// 	_, err := svc.UpdateBarber(ctx, barber2.ID, updates)
	// 	assert.Error(t, err, "ควร error เพราะ user_id ซ้ำกับ barber คนอื่น")
	// })

	t.Run("DeleteBarber_AlreadyDeleted", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 50, UserID: 5000}
		_ = svc.CreateBarber(ctx, barber)

		// ลบครั้งที่ 1
		err1 := svc.DeleteBarber(ctx, barber.ID)
		assert.NoError(t, err1)

		// ลบซ้ำ (ควรไม่พัง แต่อาจ return error หรือ ignore)
		err2 := svc.DeleteBarber(ctx, barber.ID)
		assert.Error(t, err2, "ควร error เพราะ barber ถูกลบไปแล้ว")
	})

	t.Run("GetBarberByID_AfterDelete", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 51, UserID: 5001}
		_ = svc.CreateBarber(ctx, barber)
		_ = svc.DeleteBarber(ctx, barber.ID)

		_, err := svc.GetBarberByID(ctx, barber.ID)
		assert.Error(t, err, "ควร error เพราะ barber ถูกลบแล้ว")
	})

	t.Run("ListBarbers_Empty", func(t *testing.T) {
		// ใช้ branch ใหม่ที่ยังไม่มี barber
		branchID := uint(99)
		barbers, err := svc.ListBarbersByBranch(ctx, &branchID)
		assert.NoError(t, err)
		assert.Len(t, barbers, 0, "ควรไม่มี barber ในสาขานี้")
	})

	t.Run("CreateMultipleBarbers_ThenCheckCount", func(t *testing.T) {
		branchID := uint(60)
		for i := 0; i < 10; i++ {
			_ = svc.CreateBarber(ctx, &barberBookingModels.Barber{
				BranchID: branchID,
				UserID:   uint(6000 + i),
			})
		}

		barbers, err := svc.ListBarbersByBranch(ctx, &branchID)
		assert.NoError(t, err)
		assert.Len(t, barbers, 10, "ควรมี barber ทั้งหมด 10 คนในสาขานี้")
	})

	t.Run("CreateBarber_Additional", func(t *testing.T) {
		barber := &barberBookingModels.Barber{BranchID: 61, UserID: 6010}
		err := svc.CreateBarber(ctx, barber)
		assert.NoError(t, err)
		assert.NotZero(t, barber.ID)
	})

	t.Run("GetBarberByUser", func(t *testing.T) {
		// สร้าง Barber ก่อน
		barber := &barberBookingModels.Barber{BranchID: 20, UserID: 2000}
		err := svc.CreateBarber(ctx, barber)
		assert.NoError(t, err)

		// เรียก GetBarberByUser
		found, err := svc.GetBarberByUser(ctx, 2000)
		assert.NoError(t, err)
		assert.Equal(t, barber.ID, found.ID)
		assert.Equal(t, uint(20), found.BranchID)
		assert.Equal(t, uint(2000), found.UserID)
	})

	t.Run("GetBarberByUser_NotFound", func(t *testing.T) {
		// ลองเรียก user_id ที่ไม่มี
		_, err := svc.GetBarberByUser(ctx, 9999)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	// CreateBarber_DuplicateID (manual injection)
t.Run("CreateBarber_DuplicateID", func(t *testing.T) {
	barber1 := &barberBookingModels.Barber{ID: 999, BranchID: 50, UserID: 8000}
	err1 := svc.CreateBarber(ctx, barber1)
	assert.NoError(t, err1)

	barber2 := &barberBookingModels.Barber{ID: 999, BranchID: 51, UserID: 8001}
	err2 := svc.CreateBarber(ctx, barber2)
	assert.Error(t, err2, "should fail due to duplicate ID")
})

// CreateBarber_ReuseUserIDAfterDelete
t.Run("CreateBarber_ReuseUserIDAfterDelete", func(t *testing.T) {
	barber := &barberBookingModels.Barber{BranchID: 60, UserID: 9000}
	_ = svc.CreateBarber(ctx, barber)

	err := svc.DeleteBarber(ctx, barber.ID)
	assert.NoError(t, err)

	barberNew := &barberBookingModels.Barber{BranchID: 61, UserID: 9000}
	err = svc.CreateBarber(ctx, barberNew)
	assert.NoError(t, err, "should allow reuse of user_id after deletion")
	assert.NotZero(t, barberNew.ID)
})

// // UpdateBarber_NoChanges
// t.Run("UpdateBarber_NoChanges", func(t *testing.T) {
// 	barber := &barberBookingModels.Barber{BranchID: 70, UserID: 9100}
// 	_ = svc.CreateBarber(ctx, barber)

// 	updates := &barberBookingModels.Barber{BranchID: 70, UserID: 9100} // same data
// 	updated, err := svc.UpdateBarber(ctx, barber.ID, updates)
// 	assert.NoError(t, err)
// 	assert.Equal(t, barber.ID, updated.ID)
// 	assert.Equal(t, uint(70), updated.BranchID)
// 	assert.Equal(t, uint(9100), updated.UserID)
// })
t.Run("ListBarbersByTenant", func(t *testing.T) {
	// เตรียมข้อมูล: tenant 1 มี branch 101, 102 / tenant 2 มี branch 201
	branches := []uint{101, 102, 201}
	users := []uint{5001, 5002, 5003}
	barbers := []*barberBookingModels.Barber{
		{BranchID: branches[0], UserID: users[0]}, // tenant 1
		{BranchID: branches[1], UserID: users[1]}, // tenant 1
		{BranchID: branches[2], UserID: users[2]}, // tenant 2
	}

	db.Create(&coreModels.Branch{ID: 101, TenantID: 1, Name: "Branch A"})
db.Create(&coreModels.Branch{ID: 102, TenantID: 1, Name: "Branch B"})
db.Create(&coreModels.Branch{ID: 201, TenantID: 2, Name: "Branch C"})


	for _, b := range barbers {
		err := svc.CreateBarber(ctx, b)
		assert.NoError(t, err)
	}

	t.Run("Return barbers in correct tenant", func(t *testing.T) {
		results, err := svc.ListBarbersByTenant(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, results, 2)
		assert.ElementsMatch(t, []uint{users[0], users[1]}, []uint{results[0].UserID, results[1].UserID})
	})

	t.Run("Return empty if tenant has no branches", func(t *testing.T) {
		results, err := svc.ListBarbersByTenant(ctx, 999)
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("Soft-deleted barbers should be excluded", func(t *testing.T) {
		// ลบ barber ของ tenant 1 ไป 1 คน
		err := svc.DeleteBarber(ctx, barbers[0].ID)
		assert.NoError(t, err)

		results, err := svc.ListBarbersByTenant(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, results, 1) // เหลือแค่คนเดียว
		assert.Equal(t, users[1], results[0].UserID)
	})
})


}



