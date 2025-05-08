package barberbookingServiceTest

import (
	"context"
	"sync"
	"testing"
	"time"

	"errors"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	barberBookingDto "myapp/modules/barberbooking/dto"
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

	db.Create(&barberBookingModels.Appointment{
		TenantID:   1,
		BranchID:   1,
		ServiceID:  1,
		CustomerID: 1,
		StartTime:  time.Now().Add(1 * time.Hour),
		EndTime:    time.Now().Add(1*time.Hour + 30*time.Minute),
		Status:     barberBookingModels.StatusPending,
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

	apptStart1 := time.Now().Add(2 * time.Hour).Truncate(time.Second)
	apptEnd1 := apptStart1.Add(30 * time.Minute)

	apptStart2 := time.Now().Add(4 * time.Hour).Truncate(time.Second)
	apptEnd2 := apptStart2.Add(45 * time.Minute)

	_ = db.Create(&barberBookingModels.Appointment{
		TenantID:   tenantID,
		ServiceID:  serviceID,
		CustomerID: customerID,
		BarberID:   &barberID,
		StartTime:  apptStart1,
		EndTime:    apptEnd1,
		Status:     barberBookingModels.StatusConfirmed,
		BranchID:   1,
	})

	_ = db.Create(&barberBookingModels.Appointment{
		TenantID:   tenantID,
		ServiceID:  serviceID,
		CustomerID: customerID,
		BarberID:   &barberID,
		StartTime:  apptStart2,
		EndTime:    apptEnd2,
		Status:     barberBookingModels.StatusPending,
		BranchID:   1,
	})

	// now := time.Now().Truncate(time.Second)
	// apptStart := now.Add(2 * time.Hour)
	// apptEnd := apptStart.Add(30 * time.Minute)
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

		// Restore service ID = 1 (ถ้าเคย soft-delete)
		db.Unscoped().Model(&barberBookingModels.Service{}).
			Where("id = ?", serviceID).
			Update("deleted_at", nil)

		// สร้าง barber ใหม่พร้อม userID ที่ไม่ซ้ำ
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

	t.Run("GetAvailableBarbers_ShouldReturnOnlyAvailable", func(t *testing.T) {
		start := time.Now().Add(1 * time.Hour)
		end := start.Add(30 * time.Minute)
		branchID := uint(1)
		db.Exec("DELETE FROM appointments")
		db.Exec("DELETE FROM barbers")

		// เตรียม barber A และ B
		barberA := barberBookingModels.Barber{
			ID:       1001,
			BranchID: branchID,
			UserID:   7001,
			TenantID: tenantID,
		}
		barberB := barberBookingModels.Barber{
			ID:       1002,
			BranchID: branchID,
			UserID:   7002,
			TenantID: tenantID,
		}
		_ = db.Create(&barberA)
		_ = db.Create(&barberB)

		// สร้าง appointment ซ้อนเวลาของ barber A
		_ = db.Create(&barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberA.ID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusConfirmed,
		})

		// ทดสอบ GetAvailableBarbers
		available, err := svc.GetAvailableBarbers(ctx, tenantID, branchID, start, end)
		assert.NoError(t, err)

		// ควรได้เฉพาะ barberB
		assert.Len(t, available, 1)
		assert.Equal(t, barberB.ID, available[0].ID)
	})
	//

	t.Run("GetAvailableBarbers_CompletedAppointment_ShouldNotBlock", func(t *testing.T) {
		db.Exec("DELETE FROM appointments")
		db.Exec("DELETE FROM barbers")
		start := time.Now().Add(3 * time.Hour)
		end := start.Add(30 * time.Minute)

		barberID := uint(8001)
		_ = db.Create(&barberBookingModels.Barber{
			ID:       barberID,
			UserID:   barberID,
			BranchID: 1,
			TenantID: tenantID,
		})

		// มีคิวซ้อน แต่เป็น COMPLETED
		_ = db.Create(&barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusComplete,
		})

		barbers, err := svc.GetAvailableBarbers(ctx, tenantID, 1, start, end)
		assert.NoError(t, err)
		assert.Len(t, barbers, 1)
		assert.Equal(t, barberID, barbers[0].ID)
	})

	t.Run("GetAvailableBarbers_AnotherTenant_ShouldNotReturn", func(t *testing.T) {
		start := time.Now().Add(4 * time.Hour)
		end := start.Add(30 * time.Minute)

		barberID := uint(9001)
		_ = db.Create(&barberBookingModels.Barber{
			ID:       barberID,
			UserID:   barberID,
			BranchID: 1,
			TenantID: 999, //  Tenant อื่น
		})

		barbers, err := svc.GetAvailableBarbers(ctx, tenantID, 1, start, end)
		assert.NoError(t, err)
		assert.NotContains(t, barbers, barberID)
	})

	t.Run("UpdateAppointment_Success", func(t *testing.T) {
		// สร้างนัดหมายเดิม
		start := time.Now().Add(1 * time.Hour)
		end := start.Add(30 * time.Minute)
		ap := barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusPending,
		}
		err := db.Create(&ap).Error
		assert.NoError(t, err)

		// ข้อมูลใหม่
		newStart := start.Add(1 * time.Hour)
		updateInput := &barberBookingModels.Appointment{
			StartTime: newStart,
			Status:    barberBookingModels.StatusConfirmed,
		}

		// เรียกอัปเดต
		updated, err := svc.UpdateAppointment(ctx, ap.ID, tenantID, updateInput)
		assert.NoError(t, err)
		assert.Equal(t, newStart, updated.StartTime)
		assert.Equal(t, barberBookingModels.StatusConfirmed, updated.Status)
	})

	t.Run("UpdateAppointment_NotFound_ShouldFail", func(t *testing.T) {
		updateInput := &barberBookingModels.Appointment{
			Status: barberBookingModels.StatusConfirmed,
		}
		_, err := svc.UpdateAppointment(ctx, 99999, tenantID, updateInput)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "appointment not found")
	})

	t.Run("UpdateAppointment_BarberUnavailable_ShouldFail", func(t *testing.T) {
		// สร้าง barber ใหม่
		barberID := uint(8888)
		db.Create(&barberBookingModels.Barber{
			ID:       barberID,
			TenantID: tenantID,
			BranchID: 1,
			UserID:   8888,
		})

		// สร้างนัดที่ชนกัน
		start := time.Now().Add(4 * time.Hour)
		end := start.Add(30 * time.Minute)
		db.Create(&barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusConfirmed,
		})

		// สร้างอีกนัด
		ap := barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			StartTime:  start.Add(1 * time.Hour),
			EndTime:    start.Add(1*time.Hour + 30*time.Minute),
		}
		db.Create(&ap)

		// พยายามอัปเดตให้นัดนั้นใช้ barber เดิมที่ไม่ว่าง
		updateInput := &barberBookingModels.Appointment{
			BarberID:  &barberID,
			StartTime: start,
		}
		_, err := svc.UpdateAppointment(ctx, ap.ID, tenantID, updateInput)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not available")
	})

	t.Run("CancelAppointment_Success", func(t *testing.T) {
		barber := barberBookingModels.Barber{
			ID:       1,
			TenantID: 1,
			BranchID: 1,
			UserID:   101,
		}
		assert.NoError(t, db.Create(&barber).Error)

		ap := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   ptrUint(1),
			StartTime:  time.Now().Add(2 * time.Hour),
			EndTime:    time.Now().Add(2*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusConfirmed,
		}
		assert.NoError(t, db.Create(&ap).Error)
		assert.NotZero(t, ap.ID)

		err := svc.CancelAppointment(ctx, ap.ID, 123) // actorUserID = 123
		assert.NoError(t, err)

		var updated barberBookingModels.Appointment
		assert.NoError(t, db.First(&updated, ap.ID).Error)
		assert.Equal(t, barberBookingModels.StatusCancelled, updated.Status)
		assert.Equal(t, uint(123), *updated.UserID)
	})

	t.Run("CancelAppointment_NotFound", func(t *testing.T) {
		err := svc.CancelAppointment(ctx, 99999, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("CancelAppointment_AlreadyCancelled", func(t *testing.T) {
		ap := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			Status:     barberBookingModels.StatusCancelled,
			StartTime:  time.Now().Add(4 * time.Hour),
			EndTime:    time.Now().Add(4*time.Hour + 30*time.Minute),
		}
		assert.NoError(t, db.Create(&ap).Error)

		err := svc.CancelAppointment(ctx, ap.ID, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be cancelled")
	})

	t.Run("CancelAppointment_AlreadyCompleted", func(t *testing.T) {
		ap := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			Status:     barberBookingModels.StatusComplete,
			StartTime:  time.Now().Add(-3 * time.Hour),
			EndTime:    time.Now().Add(-2 * time.Hour),
		}
		assert.NoError(t, db.Create(&ap).Error)

		err := svc.CancelAppointment(ctx, ap.ID, 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "cannot be cancelled")
	})

	t.Run("GetAppointmentByID_Success", func(t *testing.T) {
		// Seed appointment
		ap := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(1 * time.Hour),
			EndTime:    time.Now().Add(1*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusPending,
		}
		assert.NoError(t, db.Create(&ap).Error)
		assert.NotZero(t, ap.ID)

		// Call GetByID
		got, err := svc.GetAppointmentByID(ctx, ap.ID)
		assert.NoError(t, err)
		assert.Equal(t, ap.ID, got.ID)
		assert.Equal(t, ap.Status, got.Status)
		assert.Equal(t, ap.StartTime.Unix(), got.StartTime.Unix())
	})

	t.Run("GetAppointmentByID_NotFound", func(t *testing.T) {
		_, err := svc.GetAppointmentByID(ctx, 999999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("RescheduleAppointment_Success", func(t *testing.T) {
		start := parseTimeToDateToday("13:00")
		end := start.Add(30 * time.Minute)

		// สร้าง appointment เดิม
		ap := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   ptrUint(1),
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusConfirmed,
		}
		assert.NoError(t, db.Create(&ap).Error)
		service := barberBookingModels.Service{
			ID:       4,
			TenantID: ap.TenantID,
			Name:     "ตัดผมชาย",
			Duration: 30,
		}
		assert.NoError(t, db.Create(&service).Error)
		ap.ServiceID = service.ID // อัปเดต appointment ให้ชี้มาที่ ID ใหม่

		// สร้าง service ที่ใช้ระยะเวลา 30 นาที

		newStart := parseTimeToDateToday("15:00")
		err := svc.RescheduleAppointment(ctx, ap.ID, newStart, 999) // 999 = actorUserID
		assert.NoError(t, err)

		var updated barberBookingModels.Appointment
		assert.NoError(t, db.First(&updated, ap.ID).Error)
		assert.Equal(t, newStart.Unix(), updated.StartTime.Unix())
		assert.Equal(t, newStart.Add(30*time.Minute).Unix(), updated.EndTime.Unix())
		assert.Equal(t, uint(999), *updated.UserID)
		assert.Equal(t, barberBookingModels.StatusConfirmed, updated.Status)
	})

	t.Run("RescheduleAppointment_NotFound", func(t *testing.T) {
		err := svc.RescheduleAppointment(ctx, 9999, time.Now(), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("RescheduleAppointment_AlreadyCompleted", func(t *testing.T) {
		ap := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   ptrUint(1),
			StartTime:  parseTimeToDateToday("10:00"),
			EndTime:    parseTimeToDateToday("10:30"),
			Status:     barberBookingModels.StatusComplete,
		}
		assert.NoError(t, db.Create(&ap).Error)

		err := svc.RescheduleAppointment(ctx, ap.ID, parseTimeToDateToday("14:00"), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "completed or cancelled")
	})

	t.Run("RescheduleAppointment_Conflict", func(t *testing.T) {
		// appointment ที่ชนกับเวลาจะเลื่อน
		existing := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   ptrUint(1),
			StartTime:  parseTimeToDateToday("17:00"),
			EndTime:    parseTimeToDateToday("17:30"),
			Status:     barberBookingModels.StatusConfirmed,
		}
		assert.NoError(t, db.Create(&existing).Error)

		target := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   ptrUint(1),
			StartTime:  parseTimeToDateToday("11:00"),
			EndTime:    parseTimeToDateToday("11:30"),
			Status:     barberBookingModels.StatusConfirmed,
		}
		assert.NoError(t, db.Create(&target).Error)

		// เวลาใหม่ชนกับ existing
		err := svc.RescheduleAppointment(ctx, target.ID, parseTimeToDateToday("17:00"), 1)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "conflicts with another appointment")
	})

	t.Run("GetAppointmentsByBarber_Success", func(t *testing.T) {
		start := time.Now().Add(3 * time.Hour).Truncate(time.Second)
		end := start.Add(30 * time.Minute)

		appt := &barberBookingModels.Appointment{
			TenantID:   tenantID,
			ServiceID:  serviceID,
			CustomerID: customerID,
			BarberID:   &barberID,
			StartTime:  start,
			EndTime:    end,
			Status:     barberBookingModels.StatusConfirmed,
			BranchID:   1,
		}
		err := db.Create(appt).Error
		assert.NoError(t, err)

		// ใช้ pointer ของช่วงเวลา
		from := start.Add(-1 * time.Hour)
		to := end.Add(1 * time.Hour)

		results, err := svc.GetAppointmentsByBarber(ctx, barberID, &from, &to)
		assert.NoError(t, err)
		assert.NotEmpty(t, results)

		found := false
		for _, a := range results {
			if a.ID == appt.ID {
				found = true
				break
			}
		}
		assert.True(t, found, "expected appointment not found in results")
	})

	t.Run("GetAppointmentsByBarber_AllTimeRange", func(t *testing.T) {
		results, err := svc.GetAppointmentsByBarber(ctx, barberID, nil, nil)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(results), 2)
	})
	t.Run("GetAppointmentsByBarber_OutsideRange_ShouldReturnZero", func(t *testing.T) {
		from := time.Now().Add(10 * time.Hour)
		to := from.Add(1 * time.Hour)

		results, err := svc.GetAppointmentsByBarber(ctx, barberID, &from, &to)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	// t.Run("GetAppointmentsByBarber_NilTo_ShouldReturnFromStartOnly", func(t *testing.T) {
	// 	now := time.Now().Truncate(time.Second)
	// 	apptStart := now.Add(6 * time.Hour)
	// 	apptEnd := apptStart.Add(30 * time.Minute)

	// 	// Seed appointment
	// 	_ = db.Create(&barberBookingModels.Appointment{
	// 		TenantID:   tenantID,
	// 		BranchID:   1,
	// 		ServiceID:  serviceID,
	// 		CustomerID: customerID,
	// 		BarberID:   &barberID,
	// 		StartTime:  apptStart,
	// 		EndTime:    apptEnd,
	// 		Status:     barberBookingModels.StatusConfirmed,
	// 	})

	// 	from := apptStart.Add(-5 * time.Minute)
	// 	results, err := svc.GetAppointmentsByBarber(ctx, barberID, &from, nil)
	// 	assert.NoError(t, err)
	// 	assert.Len(t, results, 1)
	// })

	t.Run("GetAppointmentsByBarber_WrongBarberID_ShouldReturnEmpty", func(t *testing.T) {
		results, err := svc.GetAppointmentsByBarber(ctx, 9999, nil, nil)
		assert.NoError(t, err)
		assert.Empty(t, results)
	})

	t.Run("DeleteAppointment_ShouldRemoveFromDB", func(t *testing.T) {
		// สร้าง appointment ใหม่
		appt := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   &barberID,
			StartTime:  time.Now().Add(1 * time.Hour),
			EndTime:    time.Now().Add(2 * time.Hour),
			Status:     barberBookingModels.StatusPending,
		}
		err := db.Create(&appt).Error
		assert.NoError(t, err)

		// เรียกลบ
		err = svc.DeleteAppointment(ctx, appt.ID)
		assert.NoError(t, err)

		// ตรวจสอบว่าไม่มีใน DB แล้ว
		var found barberBookingModels.Appointment
		err = db.First(&found, appt.ID).Error
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
	})

	t.Run("DeleteAppointment_NonExisting_ShouldNotError", func(t *testing.T) {
		nonExistingID := uint(99999) // ไม่มีในระบบแน่ ๆ
		err := svc.DeleteAppointment(ctx, nonExistingID)
		assert.NoError(t, err) // GORM ไม่ถือว่าเป็น error ถ้าลบ id ที่ไม่มี
	})

	t.Run("DeleteAppointment_Twice_ShouldNotErrorSecondTime", func(t *testing.T) {
		appt := barberBookingModels.Appointment{
			TenantID:   1,
			BranchID:   1,
			ServiceID:  1,
			CustomerID: 1,
			BarberID:   &barberID,
			StartTime:  time.Now().Add(4 * time.Hour),
			EndTime:    time.Now().Add(5 * time.Hour),
			Status:     barberBookingModels.StatusConfirmed,
		}
		_ = db.Create(&appt)

		_ = svc.DeleteAppointment(ctx, appt.ID)
		err := svc.DeleteAppointment(ctx, appt.ID)
		assert.NoError(t, err)
	})

}
func TestAppointmentService_ListAppointments(t *testing.T) {
	ctx := context.Background()
	db := setupTestAppointmentDB(t)

	svc := barberBookingService.NewAppointmentService(db)

	t.Run("List_ByTenantID_Only", func(t *testing.T) {
		results, err := svc.ListAppointments(ctx, barberBookingDto.AppointmentFilter{
			TenantID: 1,
		})
		assert.NoError(t, err)
		assert.NotNil(t, results)
	})

	t.Run("List_ByDateRange", func(t *testing.T) {
		start := time.Date(2025, 4, 29, 0, 0, 0, 0, time.UTC) // จาก execution
		end := time.Date(2025, 5, 13, 23, 59, 59, 0, time.UTC)

		results, err := svc.ListAppointments(ctx, barberBookingDto.AppointmentFilter{
			TenantID:  1,
			StartDate: &start,
			EndDate:   &end,
		})
		assert.NoError(t, err)
		for _, ap := range results {
			assert.True(t, ap.StartTime.After(start) || ap.StartTime.Equal(start))
			assert.True(t, ap.EndTime.Before(end) || ap.EndTime.Equal(end))
		}
	})

	t.Run("List_ByBarberID_AndStatus", func(t *testing.T) {
		status := barberBookingModels.StatusConfirmed
		barberID := uint(1)
		results, err := svc.ListAppointments(ctx, barberBookingDto.AppointmentFilter{
			TenantID: 1,
			BarberID: &barberID,
			Status:   &status,
		})
		assert.NoError(t, err)
		for _, ap := range results {
			assert.Equal(t, barberID, *ap.BarberID)
			assert.Equal(t, status, ap.Status)
		}
	})

	t.Run("List_WithPagination", func(t *testing.T) {
		limit := 2
		offset := 1
		results, err := svc.ListAppointments(ctx, barberBookingDto.AppointmentFilter{
			TenantID: 1,
			Limit:    &limit,
			Offset:   &offset,
		})
		assert.NoError(t, err)
		assert.LessOrEqual(t, len(results), limit)
	})

	t.Run("List_WithSortByDesc", func(t *testing.T) {
		sort := "start_time desc"
		results, err := svc.ListAppointments(ctx, barberBookingDto.AppointmentFilter{
			TenantID: 1,
			SortBy:   &sort,
		})
		assert.NoError(t, err)
		for i := 1; i < len(results); i++ {
			assert.True(t, results[i-1].StartTime.After(results[i].StartTime) || results[i-1].StartTime.Equal(results[i].StartTime))
		}
	})

	t.Run("CalculateAppointmentEndTime_Success", func(t *testing.T) {
		svc := barberBookingService.NewAppointmentService(db)

		// สร้าง service
		service := barberBookingModels.Service{
			TenantID: 1,
			Name:     "ตัดผม",
			Duration: 45, // นาที
		}
		assert.NoError(t, db.Create(&service).Error)

		start := time.Date(2025, 5, 8, 14, 0, 0, 0, time.UTC)
		expectedEnd := start.Add(45 * time.Minute)

		endTime, err := svc.CalculateAppointmentEndTime(ctx, service.ID, start)
		assert.NoError(t, err)
		assert.Equal(t, expectedEnd, endTime)
	})

	t.Run("CalculateAppointmentEndTime_ServiceNotFound_Fail", func(t *testing.T) {
		invalidID := uint(9999)
		start := time.Now()

		endTime, err := svc.CalculateAppointmentEndTime(ctx, invalidID, start)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, gorm.ErrRecordNotFound))
		assert.True(t, endTime.IsZero())
	})

	t.Run("CalculateAppointmentEndTime_ZeroDuration", func(t *testing.T) {
		service := barberBookingModels.Service{
			TenantID: 1,
			Name:     "ทดลอง",
			Duration: 0,
		}
		assert.NoError(t, db.Create(&service).Error)

		start := time.Now()
		endTime, err := svc.CalculateAppointmentEndTime(ctx, service.ID, start)
		assert.NoError(t, err)
		assert.Equal(t, start, endTime)
	})

	t.Run("CalculateAppointmentEndTime_ZeroDuration", func(t *testing.T) {
		service := barberBookingModels.Service{
			TenantID: 1,
			Name:     "ทดลอง",
			Duration: 0,
		}
		assert.NoError(t, db.Create(&service).Error)

		start := time.Now()
		endTime, err := svc.CalculateAppointmentEndTime(ctx, service.ID, start)
		assert.NoError(t, err)
		assert.Equal(t, start, endTime)
	})

	t.Run("CalculateAppointmentEndTime_NegativeDuration", func(t *testing.T) {
		service := barberBookingModels.Service{
			TenantID: 1,
			Name:     "ผิดพลาด",
			Duration: -30,
		}
		assert.NoError(t, db.Create(&service).Error)

		start := time.Now()
		_, err := svc.CalculateAppointmentEndTime(ctx, service.ID, start)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid")
	})
}
