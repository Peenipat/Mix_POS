package barberbookingServiceTest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	barberbookingmodels "myapp/modules/barberbooking/models"
	barberbookingServices "myapp/modules/barberbooking/services"
)

func setupTestBarberWorkloadDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&barberbookingmodels.Barber{},
		&barberbookingmodels.BarberWorkload{},
	)
	assert.NoError(t, err)

	// Seed barber 1
	err = db.Create(&barberbookingmodels.Barber{
		ID:       1,
		BranchID: 1,
		UserID:   1,
		TenantID: 1,
	}).Error
	assert.NoError(t, err)

	// Seed barber 2
	err = db.Create(&barberbookingmodels.Barber{
		ID:       2,
		BranchID: 1,
		UserID:   2,
		TenantID: 1,
	}).Error
	assert.NoError(t, err)

	return db
}

func TestBarberWorkloadService(t *testing.T) {
	ctx := context.Background()
	db := setupTestBarberWorkloadDB(t)
	svc := barberbookingServices.NewBarberWorkloadService(db)

	today := time.Now().UTC().Truncate(24 * time.Hour)


	t.Run("UpsertBarberWorkload_InsertNew", func(t *testing.T) {
		err := svc.UpsertBarberWorkload(ctx, 1, today, 3, 5)
		assert.NoError(t, err)

		var got barberbookingmodels.BarberWorkload
		err = db.First(&got, "barber_id = ? AND strftime('%Y-%m-%d', date) = ?", 1, today.Format("2006-01-02")).Error
		assert.NoError(t, err)
		assert.Equal(t, 3, got.TotalAppointments)
		assert.Equal(t, 5, got.TotalHours)
	})

	t.Run("UpsertBarberWorkload_UpdateExisting", func(t *testing.T) {
		err := svc.UpsertBarberWorkload(ctx, 1, today, 6, 8)
		assert.NoError(t, err)
	
		got, err := svc.GetWorkloadByBarber(ctx, 1, today)
		assert.NoError(t, err)
		assert.NotNil(t, got) // ✅ ป้องกัน panic จาก got == nil
		assert.Equal(t, 6, got.TotalAppointments)
		assert.Equal(t, 8, got.TotalHours)
	})
	
	t.Run("GetWorkloadByBarber_NotFound", func(t *testing.T) {
		got, err := svc.GetWorkloadByBarber(ctx, 99, today)
		assert.NoError(t, err)
		assert.Nil(t, got)
	})

	t.Run("GetWorkloadByDate_Success", func(t *testing.T) {
		err := svc.UpsertBarberWorkload(ctx, 2, today, 2, 4)
		assert.NoError(t, err)

		workloads, err := svc.GetWorkloadByDate(ctx, today)
		assert.NoError(t, err)
		assert.Len(t, workloads, 2)
	})
}
