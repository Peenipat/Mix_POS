package barberbookingServiceTest

import (
	"context"
	"os"
	"sync/atomic"
	"testing"
	"time"

	"sync"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingService "myapp/modules/barberbooking/services"
)

func setupTestUnavailabilityDB(t *testing.T) *gorm.DB {
	_ = godotenv.Load("../../../../.env.test")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is not set.")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect test DB: %v", err)
	}

	err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error
	if err != nil {
		t.Fatalf("failed to reset schema: %v", err)
	}

	err = db.AutoMigrate(
		&barberBookingModels.Barber{},
		&barberBookingModels.Unavailability{},
	)
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	db.Create(&barberBookingModels.Barber{
		ID:       1,
		BranchID: 1,
		UserID:   101,
		TenantID: 1,
	})

	return db
}

func TestUnavailabilityService_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestUnavailabilityDB(t)
	svc := barberBookingService.NewUnavailabilityService(db)

	t.Run("CreateUnavailability_Success", func(t *testing.T) {
		input := &barberBookingModels.Unavailability{
			BranchID: ptrUint(1),
			BarberID: ptrUint(1),
			Date:     time.Now().AddDate(0, 0, 1).Truncate(24 * time.Hour),
			Reason:   "พักผ่อน",
		}
		result, err := svc.CreateUnavailability(ctx, input)
		assert.NoError(t, err)
		assert.NotZero(t, result.ID)
	})

	t.Run("CreateUnavailability_Duplicate", func(t *testing.T) {
		date := time.Now().AddDate(0, 0, 2).Truncate(24 * time.Hour)

		input1 := &barberBookingModels.Unavailability{
			BranchID: ptrUint(1),
			BarberID: ptrUint(1),
			Date:     date,
			Reason:   "เหตุผล A",
		}
		_, _ = svc.CreateUnavailability(ctx, input1)

		input2 := &barberBookingModels.Unavailability{
			BranchID: ptrUint(1),
			BarberID: ptrUint(1),
			Date:     date,
			Reason:   "เหตุผล B",
		}
		_, err := svc.CreateUnavailability(ctx, input2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})

	t.Run("GetUnavailabilitiesByBranch", func(t *testing.T) {
		from := time.Now().AddDate(0, 0, -1)
		to := time.Now().AddDate(0, 0, 10)
		results, err := svc.GetUnavailabilitiesByBranch(ctx, 1, from, to)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 1)
	})

	t.Run("UpdateUnavailability", func(t *testing.T) {
		u := &barberBookingModels.Unavailability{
			BranchID: ptrUint(1),
			BarberID: ptrUint(1),
			Date:     time.Now().AddDate(0, 0, 3),
			Reason:   "ก่อนแก้",
		}
		db.Create(u)

		err := svc.UpdateUnavailability(ctx, u.ID, map[string]interface{}{
			"reason": "แก้ไขแล้ว",
		})
		assert.NoError(t, err)

		var updated barberBookingModels.Unavailability
		_ = db.First(&updated, u.ID)
		assert.Equal(t, "แก้ไขแล้ว", updated.Reason)
	})

	t.Run("DeleteUnavailability", func(t *testing.T) {
		u := &barberBookingModels.Unavailability{
			BranchID: ptrUint(1),
			BarberID: ptrUint(1),
			Date:     time.Now().AddDate(0, 0, 4),
			Reason:   "จะลบ",
		}
		db.Create(u)

		err := svc.DeleteUnavailability(ctx, u.ID)
		assert.NoError(t, err)

		var deleted barberBookingModels.Unavailability
		err = db.Unscoped().First(&deleted, u.ID).Error
		assert.NoError(t, err)
		assert.NotNil(t, deleted.DeletedAt)
	})
	t.Run("GetUnavailabilitiesByBarber_SingleResult", func(t *testing.T) {
		date := time.Now().AddDate(0, 0, 5).Truncate(24 * time.Hour)

		db.Create(&barberBookingModels.Unavailability{
			BarberID: ptrUint(1),
			BranchID: ptrUint(1),
			Date:     date,
			Reason:   "ไปทำธุระ",
		})

		from := date.AddDate(0, 0, -1)
		to := date.AddDate(0, 0, 1)

		results, err := svc.GetUnavailabilitiesByBarber(ctx, 1, from, to)
		assert.NoError(t, err)
		assert.Len(t, results, 1)
		assert.Equal(t, "ไปทำธุระ", results[0].Reason)
	})

	t.Run("GetUnavailabilitiesByBarber_NoResultInRange", func(t *testing.T) {
		from := time.Now().AddDate(0, 0, -30)
		to := time.Now().AddDate(0, 0, -25)

		results, err := svc.GetUnavailabilitiesByBarber(ctx, 1, from, to)
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("GetUnavailabilitiesByBarber_InvalidBarber", func(t *testing.T) {
		// ใช้ barberID ที่ไม่มีอยู่
		from := time.Now().AddDate(0, 0, -1)
		to := time.Now().AddDate(0, 0, 10)

		results, err := svc.GetUnavailabilitiesByBarber(ctx, 9999, from, to)
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("GetUnavailabilitiesByBarber_OverlappingDateRanges", func(t *testing.T) {
		// สร้างวันหยุดหลายวัน
		startDate := time.Now().AddDate(0, 0, 10).Truncate(24 * time.Hour)
		for i := 0; i < 3; i++ {
			db.Create(&barberBookingModels.Unavailability{
				BarberID: ptrUint(1),
				BranchID: ptrUint(1),
				Date:     startDate.AddDate(0, 0, i),
				Reason:   "พักร้อน",
			})
		}

		from := startDate
		to := startDate.AddDate(0, 0, 2)
		results, err := svc.GetUnavailabilitiesByBarber(ctx, 1, from, to)
		assert.NoError(t, err)
		assert.Len(t, results, 3)
	})

	t.Run("GetUnavailabilitiesByBarber_BarberIDNil", func(t *testing.T) {
		// กรณี barber_id == NULL (เช่น วันหยุดของสาขา)
		// ต้องไม่ return ใน getByBarber

		db.Create(&barberBookingModels.Unavailability{
			BarberID: nil,
			BranchID: ptrUint(1),
			Date:     time.Now().AddDate(0, 0, 7),
			Reason:   "วันหยุดสาขา",
		})

		from := time.Now().AddDate(0, 0, 6)
		to := time.Now().AddDate(0, 0, 8)

		results, err := svc.GetUnavailabilitiesByBarber(ctx, 1, from, to)
		assert.NoError(t, err)
		for _, u := range results {
			assert.NotNil(t, u.BarberID)
		}
	})

	t.Run("CreateUnavailability_Concurrent_DuplicateProtection", func(t *testing.T) {
		date := time.Now().AddDate(0, 0, 15).Truncate(24 * time.Hour)
		n := 10 // จำนวน concurrent goroutine ที่จะยิงพร้อมกัน
	

		var wg sync.WaitGroup
		var successCount int32

		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				input := &barberBookingModels.Unavailability{
					BarberID: ptrUint(1),
					BranchID: ptrUint(1),
					Date:     date,
					Reason:   "หยุดพร้อมกัน",
				}

				_, err := svc.CreateUnavailability(ctx, input)
				if err == nil {
					atomic.AddInt32(&successCount, 1)
				}
			}()
		}
		wg.Wait()

		assert.Equal(t, int32(1), successCount, "ควรมีเพียง 1 รายการที่ถูกสร้างสำเร็จ")
	})

}

func ptrUint(v uint) *uint {
	return &v
}
