package barberbookingServiceTest

import (
	"testing"
	"context"
	"github.com/stretchr/testify/assert"
	"time"
	"gorm.io/gorm"
	"os"
	"gorm.io/driver/postgres"
	"github.com/joho/godotenv"
	"fmt"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingServices "myapp/modules/barberbooking/services"
	barberBookingDto "myapp/modules/barberbooking/dto"
	coreModels "myapp/modules/core/models"
)

func setupTestWorkingHourDB(t *testing.T) *gorm.DB {
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
		&barberBookingModels.WorkingHour{},
		&coreModels.Branch{},

	)
	
	if err != nil {
		t.Fatalf("migration failed: %v", err)
	}

	return db
}

func parseTimeToDateToday(s string) time.Time {
	t, err := time.Parse("15:04", s)
	if err != nil {
		panic(fmt.Sprintf("invalid time format: %s", s))
	}

	now := time.Now()
	return time.Date(now.Year(), now.Month(), now.Day(), t.Hour(), t.Minute(), 0, 0, now.Location())
}


func TestWorkingHourService(t *testing.T) {
	ctx := context.Background()
	db := setupTestWorkingHourDB(t)
	svc := barberBookingServices.NewWorkingHourService(db)

	tenant := coreModels.Tenant{
        ID:       1,
        Name:     "TestTenant",
        Domain:   "test.local",
        IsActive: true,
        // หาก struct มีฟิลด์อื่น เช่น CreatedAt, ก็ควรใส่ default หรือเว้นไว้
    }
    if err := db.Create(&tenant).Error; err != nil {
        t.Fatalf("failed to create tenant: %v", err)
    }

    // 2) สร้าง Branch ที่อ้างถึง TenantID=1
	addr := "123 Test St."
    branch := coreModels.Branch{
        ID:       1,
        TenantID: tenant.ID,
        Name:     "Test Branch",
        Address:  &addr,
    }
    if err := db.Create(&branch).Error; err != nil {
        t.Fatalf("failed to create branch: %v", err)
    }
	t.Run("UpdateWorkingHours_Success", func(t *testing.T) {
		input := []barberBookingDto.WorkingHourInput{
			{Weekday: 0, StartTime: parseTimeToDateToday("09:00"), EndTime: parseTimeToDateToday("18:00")},
			{Weekday: 1, StartTime: parseTimeToDateToday("10:00"), EndTime: parseTimeToDateToday("17:00")},
		}
		err := svc.UpdateWorkingHours(ctx, 1, 1,input)
		assert.NoError(t, err)

		var results []barberBookingModels.WorkingHour
		assert.NoError(t, db.Where("branch_id = ?", 1).Find(&results).Error)
		assert.Equal(t, 2, len(results))
	})

	t.Run("GetWorkingHours_ReturnsSorted", func(t *testing.T) {
		results, err := svc.GetWorkingHours(ctx, 1,1)
		assert.NoError(t, err)
		assert.Equal(t, 0, results[0].Weekday)
		assert.Equal(t, 1, results[1].Weekday)
	})

	

	

	t.Run("UpdateWorkingHours_InvalidWeekday_Fail", func(t *testing.T) {
		input := []barberBookingDto.WorkingHourInput{
			{Weekday: 7, StartTime: parseTimeToDateToday("09:00"), EndTime: parseTimeToDateToday("17:00")}, //  invalid
		}
		err := svc.UpdateWorkingHours(ctx, 1, 1,input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid weekday")
	})

	t.Run("GetWorkingHours_NoData", func(t *testing.T) {
		results, err := svc.GetWorkingHours(ctx, 9999,1) // ใช้ branchID ที่ไม่มี
		assert.NoError(t, err)
		assert.Len(t, results, 0)
	})

	t.Run("CreateWorkingHours_Success", func(t *testing.T) {
		input := barberBookingDto.WorkingHourInput{
			Weekday:   4, // Wednesday
			StartTime: parseTimeToDateToday("9:00"),
			EndTime:   parseTimeToDateToday("17:00"),
		}
		err := svc.CreateWorkingHours(ctx, 1, input)
		assert.NoError(t, err)
	
		var result barberBookingModels.WorkingHour
		err = db.Where("branch_id = ? AND weekday = ?", 1, 4).First(&result).Error
		assert.NoError(t, err)
		assert.Equal(t, input.StartTime.Hour(), result.StartTime.Hour())
	})
	
	t.Run("CreateWorkingHours_Duplicate", func(t *testing.T) {
        input := barberBookingDto.WorkingHourInput{
            Weekday:   3,
            StartTime: parseTimeToDateToday("10:00"),
            EndTime:   parseTimeToDateToday("18:00"),
        }

        // 1) สร้าง entry แรก
        if err := svc.CreateWorkingHours(ctx, branch.ID, input); err != nil {
            t.Fatalf("expected first insert to succeed, got error: %v", err)
        }

        // 2) สร้างซ้ำอีกครั้ง → ควร error "already exists"
        err := svc.CreateWorkingHours(ctx, branch.ID, input)
        assert.Error(t, err)
        assert.Contains(t, err.Error(), "already exists")
    })
	
	t.Run("CreateWorkingHours_InvalidWeekday", func(t *testing.T) {
		input := barberBookingDto.WorkingHourInput{
			Weekday:   7, //  invalid
			StartTime: parseTimeToDateToday("09:00"),
			EndTime:   parseTimeToDateToday("17:00"),
		}
		err := svc.CreateWorkingHours(ctx, 1, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid weekday")
	})
	
	t.Run("CreateWorkingHours_StartAfterEnd_Fail", func(t *testing.T) {
		input := barberBookingDto.WorkingHourInput{
			Weekday:   4,
			StartTime: parseTimeToDateToday("18:00"),
			EndTime:   parseTimeToDateToday("09:00"),
		}
		err := svc.CreateWorkingHours(ctx, 1, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start time must be before")
	})

	t.Run("CreateWorkingHours_ZeroDuration_Fail", func(t *testing.T) {
		input := barberBookingDto.WorkingHourInput{
			Weekday:   5, // Friday
			StartTime: parseTimeToDateToday("09:00"),
			EndTime:   parseTimeToDateToday("09:00"),
		}
		err := svc.CreateWorkingHours(ctx, 1, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start time must be before end time")
	})

	t.Run("CreateWorkingHours_OverMidnight_Fail", func(t *testing.T) {
		input := barberBookingDto.WorkingHourInput{
			Weekday:   6, // Saturday
			StartTime: parseTimeToDateToday("22:00"),
			EndTime:   parseTimeToDateToday("01:00"),
		}
		err := svc.CreateWorkingHours(ctx, 1, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start time must be before end time")
	})

	t.Run("CreateWorkingHours_OverMidnight_Fail", func(t *testing.T) {
		input := barberBookingDto.WorkingHourInput{
			Weekday:   6, // Saturday
			StartTime: parseTimeToDateToday("22:00"),
			EndTime:   parseTimeToDateToday("01:00"),
		}
		err := svc.CreateWorkingHours(ctx, 1, input)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "start time must be before end time")
	})
	
	
	
}


