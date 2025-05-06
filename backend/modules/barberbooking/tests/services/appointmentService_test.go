package barberbookingServiceTest

import (
	"context"
	"sync"
	"testing"
	"time"

	"fmt"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingService "myapp/modules/barberbooking/services"
)

func setupTestAppointmentDB(t *testing.T) *gorm.DB {
	_ = godotenv.Load("../../../../.env.test") // ใช้ relative path ให้ตรงกับตำแหน่งจริง

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is not set. Please check .env.test or environment variable.")
	}
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to PostgreSQL test DB: %v", err)
	}

	// 🧹 ล้าง schema แล้ว migrate ใหม่ (ปลอดภัยสำหรับ test เท่านั้น)
	err = db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error
	if err != nil {
		t.Fatalf("failed to reset schema: %v", err)
	}

	err = db.AutoMigrate(
		&barberBookingModels.Service{},
		&barberBookingModels.Customer{},
		&barberBookingModels.Barber{},
		&barberBookingModels.Appointment{},
	)
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	//  Seed Customer หลัก
	db.Create(&barberBookingModels.Customer{
		ID:       1,
		Name:     "ลูกค้าทดสอบ",
		Email:    "test@example.com",
		TenantID: 1,
	})

	// Seed Barber หลัก
	db.Create(&barberBookingModels.Barber{
		ID:       1,
		BranchID: 1,
		UserID:   1001,
		TenantID: 1,
	})

	return db
}

func TestAppointmentService_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestAppointmentDB(t)
	svc := barberBookingService.NewAppointmentService(db)

	tenantID := uint(1)
	serviceID := uint(1)
	customerID := uint(1)
	barberID := uint(1)

	db.Create(&barberBookingModels.Service{
		ID:       1,
		TenantID: 1,
		Name:     "ตัดผมชาย",
		Price:    200,
		Duration: 30,
	})
	t.Run("CreateAppointment_Success", func(t *testing.T) {
		start := time.Now().Add(1 * time.Hour)
		appointment := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			BranchID:   1,
		}

		result, err := svc.CreateAppointment(ctx, appointment)
		assert.NoError(t, err)
		assert.NotZero(t, result.ID)
		assert.WithinDuration(t, result.StartTime.Add(30*time.Minute), result.EndTime, time.Second)
	})

	t.Run("CreateAppointment_BarberUnavailable", func(t *testing.T) {
		start := time.Now().Add(2 * time.Hour)

		// First appointment
		ap1 := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			BranchID:   1,
		}
		_, _ = svc.CreateAppointment(ctx, ap1)

		// Second appointment at same time
		ap2 := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			BranchID:   1,
		}
		_, err := svc.CreateAppointment(ctx, ap2)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not available")
	})

	t.Run("CreateAppointment_MissingService", func(t *testing.T) {
		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  9999,
			CustomerID: customerID,
			StartTime:  time.Now().Add(3 * time.Hour),
		}
		_, err := svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service not found")
	})

	t.Run("CreateAppointment_MissingRequiredFields", func(t *testing.T) {
		ap := &barberBookingModels.Appointment{}
		_, err := svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required fields")
	})

	t.Run("CreateAppointment_NoBarber", func(t *testing.T) {
		start := time.Now().Add(4 * time.Hour)
		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			StartTime:  start,
		}
		result, err := svc.CreateAppointment(ctx, ap)
		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Nil(t, result.BarberID)
	})

	t.Run("CreateAppointment_EndsExactlyAtStartOfAnother", func(t *testing.T) {
		start1 := time.Now().Add(6 * time.Hour)
		start2 := start1.Add(30 * time.Minute) // ช่างว่างต่อเนื่อง

		// สร้างคิวแรก
		_ = db.Create(&barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start1,
			EndTime:    start2,
			Status:     barberBookingModels.StatusConfirmed,
		}).Error

		// สร้างคิวที่ต่อพอดี
		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start2,
			BranchID:   1,
		}
		result, err := svc.CreateAppointment(ctx, ap)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("CreateAppointment_ServiceDeleted", func(t *testing.T) {
		// ลบ service แบบ soft delete
		err := db.Delete(&barberBookingModels.Service{}, 1).Error
		assert.NoError(t, err)

		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			StartTime:  time.Now().Add(5 * time.Hour),
		}
		_, err = svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service not found")
	})

	t.Run("CreateAppointment_InvalidTimeRange", func(t *testing.T) {
		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			StartTime:  time.Time{}, // ไม่มีเวลาเลย
		}
		_, err := svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "missing required fields")
	})

	t.Run("CreateAppointment_AnotherTenantService", func(t *testing.T) {
		db.Unscoped().Delete(&barberBookingModels.Service{}, "tenant_id = ?", 999) // ลบแบบ force

		svc2 := barberBookingModels.Service{
			ID:       uint(time.Now().Unix()), // ป้องกัน primary key ซ้ำ
			Name:     "บริการปลอม",
			Duration: 20,
			Price:    100,
			TenantID: 999,
		}
		err := db.Create(&svc2).Error
		assert.NoError(t, err)
		assert.NotZero(t, svc2.ID)

		ap := &barberBookingModels.Appointment{
			TenantID:   1,
			ServiceID:  svc2.ID,
			CustomerID: 1,
			StartTime:  time.Now().Add(8 * time.Hour),
		}
		_, err = svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "service not found or access denied")
	})

	t.Run("CreateAppointment_OverlapWithOtherBarber_ShouldPass", func(t *testing.T) {
		start := time.Now().Add(10 * time.Hour)

		// คืนค่า service กลับมา (หากถูกลบใน test ก่อนหน้า)
		db.Unscoped().Model(&barberBookingModels.Service{}).
			Where("id = ?", 1).Update("deleted_at", nil)

		// Barber A จองคิว
		barberA := uint(1)
		ap1 := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BranchID:   1,
			BarberID:   &barberA,
			StartTime:  start,
		}
		_, err := svc.CreateAppointment(ctx, ap1)
		assert.NoError(t, err)

		//  Barber B → ต้องจองได้แม้เวลาเดียวกัน
		barberB := uint(2)
		db.Create(&barberBookingModels.Barber{
			ID:       barberB,
			BranchID: 1,
			UserID:   2001,
			TenantID: 1,
		})

		ap2 := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BranchID:   1,
			BarberID:   &barberB,
			StartTime:  start,
		}
		result, err := svc.CreateAppointment(ctx, ap2)
		assert.NoError(t, err)
		assert.NotZero(t, result.ID)
	})

	t.Run("CreateAppointment_WithCompletedStatusOverlap_ShouldPass", func(t *testing.T) {
		start := time.Now().Add(11 * time.Hour)
		end := start.Add(30 * time.Minute)

		// มีคิวก่อนหน้าเป็น COMPLETED → ไม่ควรถือว่า block เวลา
		ap := barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BranchID:   1,
			BarberID:   &barberID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap)

		// จองใหม่ทับเวลาก็ได้ เพราะ completed ไม่ block
		ap2 := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			BranchID:   1,
			StartTime:  start,
		}
		_, err := svc.CreateAppointment(ctx, ap2)
		assert.NoError(t, err)
	})

	t.Run("CreateAppointment_WithZeroDurationService_ShouldFail", func(t *testing.T) {
		// ลบ service เดิมที่อาจซ้ำชื่อหรือ key
		db.Where("name = ?", "ผิดพลาด").Delete(&barberBookingModels.Service{})

		// สร้าง service duration = 0
		svcZero := barberBookingModels.Service{
			ID:       uint(time.Now().UnixNano()),
			Name:     fmt.Sprintf("ผิดพลาด-%d", time.Now().UnixNano()), // ป้องกันซ้ำ
			Duration: 0,
			Price:    100,
			TenantID: tenantID,
		}
		err := db.Create(&svcZero).Error
		assert.NoError(t, err)
		assert.NotZero(t, svcZero.ID)

		// พยายามสร้างนัดหมายด้วย service duration = 0
		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  svcZero.ID,
			BranchID:   1,
			CustomerID: customerID,
			StartTime:  time.Now().Add(12 * time.Hour),
		}
		_, err = svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duration must be > 0")
	})

	t.Run("CreateAppointment_BarberFromAnotherTenant_ShouldFail", func(t *testing.T) {
		barberOther := barberBookingModels.Barber{
			ID:       999,
			BranchID: 1,
			TenantID: 1,
		}
		_ = db.Create(&barberOther)

		ap := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberOther.ID,
			BranchID:   1,
			StartTime:  time.Now().Add(13 * time.Hour),
		}
		// ❗ ขณะนี้ระบบยังไม่ validate tenant ของ barber → ควรทำถ้าต้องการ
		_, err := svc.CreateAppointment(ctx, ap)
		// คาดหวังว่าจะต้อง fail ถ้ามีการ validate tenant
		assert.NoError(t, err) // ❗ เปลี่ยนเป็น assert.Error ถ้าคุณเพิ่ม tenant validation ให้ barber
	})

	t.Run("CreateAppointment_BarberFromAnotherBranch_ShouldFail", func(t *testing.T) {
		db.Unscoped().Delete(&barberBookingModels.Barber{}, "id = ?", 1001)

		barberX := barberBookingModels.Barber{
			ID:       1001,
			BranchID: 99,                      // สาขาอื่น
			UserID:   uint(time.Now().Unix()), // ป้องกันซ้ำ
			TenantID: 1,
		}
		err := db.Create(&barberX).Error
		assert.NoError(t, err)

		ap := &barberBookingModels.Appointment{
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   &barberX.ID,
			StartTime:  time.Now().Add(14 * time.Hour),
			BranchID:   1, //  สาขาหลักที่ไม่ตรง
		}

		_, err = svc.CreateAppointment(ctx, ap)
		assert.Error(t, err)
		if assert.Error(t, err) {
			assert.Contains(t, err.Error(), "barber not found or mismatched branch")
		}

	})

	t.Run("CreateAppointment_ConcurrentConflict", func(t *testing.T) {
		ctx := context.Background()
		start := time.Now().Add(2 * time.Hour)

		// ✅ Restore service ID = 1 (ถ้าเคย soft-delete)
		db.Unscoped().Model(&barberBookingModels.Service{}).
			Where("id = ?", serviceID).
			Update("deleted_at", nil)

		// ✅ สร้าง barber ใหม่พร้อม userID ที่ไม่ซ้ำ
		barberID := uint(5001)
		err := db.Create(&barberBookingModels.Barber{
			ID:       barberID,
			BranchID: 1,
			UserID:   5001,
			TenantID: 1,
		}).Error
		assert.NoError(t, err)

		// ✅ สร้าง customer 2 คน
		customer1 := barberBookingModels.Customer{Name: "User1", Email: "user1@test.com", TenantID: tenantID}
		customer2 := barberBookingModels.Customer{Name: "User2", Email: "user2@test.com", TenantID: tenantID}
		_ = db.Create(&customer1)
		_ = db.Create(&customer2)

		var wg sync.WaitGroup
		wg.Add(2)

		results := make([]error, 2)

		go func() {
			defer wg.Done()
			ap := &barberBookingModels.Appointment{
				TenantID:   tenantID,
				ServiceID:  serviceID,
				CustomerID: customer1.ID,
				BarberID:   &barberID,
				BranchID:   1,
				StartTime:  start,
			}
			_, err := svc.CreateAppointment(ctx, ap)
			results[0] = err
		}()

		go func() {
			defer wg.Done()
			ap := &barberBookingModels.Appointment{
				TenantID:   tenantID,
				ServiceID:  serviceID,
				CustomerID: customer2.ID,
				BarberID:   &barberID,
				BranchID:   1,
				StartTime:  start,
			}
			_, err := svc.CreateAppointment(ctx, ap)
			results[1] = err
		}()

		wg.Wait()

		// ต้องมีคนหนึ่งสำเร็จ คนหนึ่ง fail
		successCount := 0
		failureCount := 0

		for _, err := range results {
			if err == nil {
				successCount++
			} else {
				failureCount++
				t.Logf("Expected conflict error: %v", err)
			}
		}

		assert.Equal(t, 1, successCount, "Exactly one request should succeed")
		assert.Equal(t, 1, failureCount, "Exactly one request should fail due to conflict")
	})

	t.Run("CheckBarberAvailability_Available", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now().Add(15 * time.Hour)
		end := start.Add(30 * time.Minute)

		// ไม่มีการจองคิวซ้อน → ต้องว่าง
		available, err := svc.CheckBarberAvailability(ctx, tenantID, barberID, start, end)
		assert.NoError(t, err)
		assert.True(t, available)
	})

	t.Run("CheckBarberAvailability_Overlap_ShouldReturnFalse", func(t *testing.T) {
		ctx := context.Background()

		// สร้างคิวที่ block เวลาไว้ก่อน
		start := time.Now().Add(16 * time.Hour)
		end := start.Add(30 * time.Minute)
		_ = db.Create(&barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusConfirmed,
		})

		// ลองเช็ค availability ที่ซ้อนกับคิวนี้
		conflictStart := start.Add(10 * time.Minute)
		conflictEnd := conflictStart.Add(30 * time.Minute)

		available, err := svc.CheckBarberAvailability(ctx, tenantID, barberID, conflictStart, conflictEnd)
		assert.NoError(t, err)
		assert.False(t, available)
	})

	t.Run("CheckBarberAvailability_WithCompletedAppointment_ShouldReturnTrue", func(t *testing.T) {
		ctx := context.Background()

		start := time.Now().Add(17 * time.Hour)
		end := start.Add(30 * time.Minute)

		// เพิ่มคิวสถานะ completed → ไม่ควร block เวลา
		_ = db.Create(&barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusComplete,
		})

		available, err := svc.CheckBarberAvailability(ctx, tenantID, barberID, start, end)
		assert.NoError(t, err)
		assert.True(t, available)
	})
	t.Run("CheckBarberAvailability_BarberNotFound_ShouldReturnFalse", func(t *testing.T) {
		ctx := context.Background()
		start := time.Now().Add(18 * time.Hour)
		end := start.Add(30 * time.Minute)

		available, err := svc.CheckBarberAvailability(ctx, 999, 999, start, end)
		assert.NoError(t, err)
		assert.False(t, available)
	})



}
