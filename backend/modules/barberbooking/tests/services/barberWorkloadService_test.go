package barberbookingServiceTest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"github.com/stretchr/testify/require"

	barberBookingModels "myapp/modules/barberbooking/models"
	coreModels "myapp/modules/core/models"
	barberbookingServices "myapp/modules/barberbooking/services"
)

func setupTestBarberWorkloadDB(t *testing.T) *gorm.DB {
    // 1. in-memory DB ใหม่ทุกครั้ง
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    require.NoError(t, err)

    // 2. สร้าง schema: branches, barbers, barber_workloads
    require.NoError(t, db.AutoMigrate(
        &coreModels.Branch{},
        &barberBookingModels.Barber{},
        &barberBookingModels.BarberWorkload{},
    ))

    // 3. Seed branches
    branches := []coreModels.Branch{
        {ID: 1, TenantID: 1, Name: "Branch One"},
        {ID: 2, TenantID: 2, Name: "Branch Two"},
    }
    for _, br := range branches {
        require.NoError(t, db.Create(&br).Error)
    }

    // 4. Seed barbers
    //    - สาขา 1: Barber ID = 1,2
    //    - สาขา 2: Barber ID = 3
    barbers := []barberBookingModels.Barber{
        {ID: 1, BranchID: 1, UserID: 1, TenantID: 1},
        {ID: 2, BranchID: 1, UserID: 2, TenantID: 1},
        {ID: 3, BranchID: 2, UserID: 3, TenantID: 2},
    }
    for _, b := range barbers {
        require.NoError(t, db.Create(&b).Error)
    }

    return db
}

 

func TestGetWorkloadByBarber(t *testing.T) {
    input := time.Date(2025, time.May, 14, 15, 30, 0, 0, time.UTC)
    truncated := input.Truncate(24 * time.Hour)

    t.Run("Found_RecordExists", func(t *testing.T) {
        db := setupTestBarberWorkloadDB(t)
        svc := barberbookingServices.NewBarberWorkloadService(db)

        // seed on 2025-05-14
        rec := barberBookingModels.BarberWorkload{
            BarberID:          7,
            Date:              input,
            TotalAppointments: 4,
            TotalHours:        9,
        }
        require.NoError(t, db.Create(&rec).Error)

        got, err := svc.GetWorkloadByBarber(context.Background(), 7, input)
        require.NoError(t, err)
        require.NotNil(t, got)

        assert.Equal(t, uint(7), got.BarberID)
        assert.True(t, got.Date.Equal(rec.Date), "expected %v, got %v", rec.Date, got.Date)
        assert.Equal(t, 4, got.TotalAppointments)
        assert.Equal(t, 9, got.TotalHours)
    })

    t.Run("NotFound_NoSuchRecord", func(t *testing.T) {
        db := setupTestBarberWorkloadDB(t)
        svc := barberbookingServices.NewBarberWorkloadService(db)

        // no seed at all
        got, err := svc.GetWorkloadByBarber(context.Background(), 42, input)
        require.NoError(t, err)
        assert.Nil(t, got)
    })

    t.Run("DifferentDay_RecordNotReturned", func(t *testing.T) {
        db := setupTestBarberWorkloadDB(t)
        svc := barberbookingServices.NewBarberWorkloadService(db)
        
        // seed on 2025-05-15 only
        future := truncated.Add(24 * time.Hour).Add(1 * time.Hour)
        rec := barberBookingModels.BarberWorkload{
            BarberID:          7,
            Date:              future,
            TotalAppointments: 2,
            TotalHours:        3,
        }
        require.NoError(t, db.Create(&rec).Error)

        // query 2025-05-14 → should not find
        got, err := svc.GetWorkloadByBarber(context.Background(), 7, input)
        require.NoError(t, err)
        assert.Nil(t, got)
    })

    t.Run("DBError_ShouldPropagate", func(t *testing.T) {
        db := setupTestBarberWorkloadDB(t)
        svc := barberbookingServices.NewBarberWorkloadService(db)

        // seed one so DB is open, then close underlying sql.DB
        require.NoError(t, db.Create(&barberBookingModels.BarberWorkload{
            BarberID: 7, Date: input,
        }).Error)
        sqlDB, _ := db.DB()
        sqlDB.Close()

        _, err := svc.GetWorkloadByBarber(context.Background(), 7, input)
        assert.Error(t, err)
    })
}


func TestGetWorkloadSummaryByBranch(t *testing.T) {
    db := setupTestBarberWorkloadDB(t)
    svc := barberbookingServices.NewBarberWorkloadService(db)

    testDate := time.Date(2025, time.May, 14, 0, 0, 0, 0, time.UTC)
    nextDay  := testDate.Add(24 * time.Hour)

    // Seed workloads
    require.NoError(t, db.Create(&barberBookingModels.BarberWorkload{
        BarberID: 1, Date: testDate.Add(2 * time.Hour),
    }).Error)
    require.NoError(t, db.Create(&barberBookingModels.BarberWorkload{
        BarberID: 2, Date: testDate.Add(5 * time.Hour),
    }).Error)
    // workload นอกช่วง date
    require.NoError(t, db.Create(&barberBookingModels.BarberWorkload{
        BarberID: 1, Date: nextDay.Add(1 * time.Hour),
    }).Error)

    t.Run("NoFilters_ReturnAllBranches", func(t *testing.T) {
        sums, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 0, 0)
        require.NoError(t, err)
        // ควรได้ทั้ง 2 สาขา
        assert.Len(t, sums, 2)
    })

    t.Run("FilterByTenant1_ReturnOnlyBranch1", func(t *testing.T) {
        sums, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 1, 0)
        require.NoError(t, err)
        // tenant 1 มีแค่ branch 1
        assert.Len(t, sums, 1)
        s := sums[0]
        assert.Equal(t, uint(1), s.TenantID)
        assert.Equal(t, uint(1), s.BranchID)
        assert.Equal(t, 2, s.NumWorked)     // workloads ของ barber 1&2
        assert.Equal(t, 2, s.TotalBarbers)  // barbers 1&2
    })

    t.Run("FilterByBranch2_ReturnOnlyBranch2", func(t *testing.T) {
        sums, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 0, 2)
        require.NoError(t, err)
        // branch 2 อยู่ tenant 2 แต่ไม่มี workloads
        assert.Len(t, sums, 1)
        s := sums[0]
        assert.Equal(t, uint(2), s.TenantID)
        assert.Equal(t, uint(2), s.BranchID)
        assert.Equal(t, 0, s.NumWorked)
        assert.Equal(t, 1, s.TotalBarbers)  // barber 3
    })

    t.Run("FilterByTenant2Branch1_NoMatch", func(t *testing.T) {
        sums, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 2, 1)
        require.NoError(t, err)
        // tenant 2 กับ branch 1 ไม่สัมพันธ์กัน
        assert.Empty(t, sums)
    })

    t.Run("FilterByInvalidTenant_ReturnEmpty", func(t *testing.T) {
        sums, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 999, 0)
        require.NoError(t, err)
        assert.Empty(t, sums)
    })

    t.Run("FilterByInvalidBranch_ReturnEmpty", func(t *testing.T) {
        sums, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 0, 999)
        require.NoError(t, err)
        assert.Empty(t, sums)
    })

    t.Run("ServiceError_ShouldPropagate", func(t *testing.T) {
        sqlDB, _ := db.DB()
        sqlDB.Close()
        _, err := svc.GetWorkloadSummaryByBranch(context.Background(), testDate, 0, 0)
        assert.Error(t, err)
    })
}

func TestUpsertBarberWorkload(t *testing.T) {
    db := setupTestBarberWorkloadDB(t)
    svc := barberbookingServices.NewBarberWorkloadService(db)

    ctx := context.Background()
    date := time.Date(2025, time.May, 14, 15, 30, 0, 0, time.UTC)
    truncated := date.Truncate(24 * time.Hour)

    t.Run("InsertNewRecord", func(t *testing.T) {
        // ยังไม่มี record เดิม
        err := svc.UpsertBarberWorkload(ctx, 42, date, 3, 5)
        require.NoError(t, err)

        // ตรวจใน DB ว่ามี record เกิดขึ้น
        var w barberBookingModels.BarberWorkload
        err = db.First(&w, "barber_id = ? AND date = ?", 42, truncated).Error
        require.NoError(t, err)

        assert.Equal(t, uint(42), w.BarberID)
        assert.True(t, w.Date.Equal(truncated))
        assert.Equal(t, 3, w.TotalAppointments)
        assert.Equal(t, 5, w.TotalHours)
    })

    t.Run("AccumulateExistingRecord", func(t *testing.T) {
        // เรียกอีกครั้งบนเดียวกัน: เพิ่ม appointments=2, hours=1
        err := svc.UpsertBarberWorkload(ctx, 42, date, 2, 1)
        require.NoError(t, err)

        // ควรสะสมจาก (3,5) ➞ (5,6)
        var w barberBookingModels.BarberWorkload
        err = db.First(&w, "barber_id = ? AND date = ?", 42, truncated).Error
        require.NoError(t, err)

        assert.Equal(t, 5, w.TotalAppointments)
        assert.Equal(t, 6, w.TotalHours)
    })

    t.Run("SeparateDatesCreateDistinct", func(t *testing.T) {
        // ใช้วันถัดไป จะต้องสร้าง record ใหม่ ไม่ไปกระทบของเดิม
        nextDay := truncated.Add(24 * time.Hour)
        err := svc.UpsertBarberWorkload(ctx, 42, nextDay, 1, 2)
        require.NoError(t, err)

        // ดึงสองเรคอร์ด
        var workloads []barberBookingModels.BarberWorkload
        err = db.Find(&workloads, "barber_id = ?", 42).Error
        require.NoError(t, err)

        // ควรมี 2 เรคอร์ด: หนึ่งสำหรับ truncated, หนึ่งสำหรับ nextDay
        assert.Len(t, workloads, 2)

        // แยกดูแต่ละวัน
        var gotTrunc, gotNext barberBookingModels.BarberWorkload
        for _, w := range workloads {
            if w.Date.Equal(truncated) {
                gotTrunc = w
            } else if w.Date.Equal(nextDay) {
                gotNext = w
            }
        }
        assert.Equal(t, 5, gotTrunc.TotalAppointments)
        assert.Equal(t, 6, gotTrunc.TotalHours)
        assert.Equal(t, 1, gotNext.TotalAppointments)
        assert.Equal(t, 2, gotNext.TotalHours)
    })

    t.Run("DBError_ShouldPropagate", func(t *testing.T) {
        // ปิด DB เพื่อให้เกิด error
        sqlDB, err := db.DB()
        require.NoError(t, err)
        sqlDB.Close()

        err = svc.UpsertBarberWorkload(ctx, 99, date, 1, 1)
        assert.Error(t, err)
    })
}


