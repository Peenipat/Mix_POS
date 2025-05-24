package barberbookingServiceTest

import (
	"context"
	"testing"
	"time"
	"strings"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	barberbookingmodels "myapp/modules/barberbooking/models"
	barberbookingservices "myapp/modules/barberbooking/services"
)

func setupAppointmentLogTestDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(
		&barberbookingmodels.Appointment{},
		&barberbookingmodels.Customer{},
		&barberbookingmodels.AppointmentStatusLog{},
	)
	assert.NoError(t, err)

	// Seed appointment
	err = db.Create(&barberbookingmodels.Appointment{
		ID:        1,
		TenantID:  1,
		CustomerID: 1,
		Status:    barberbookingmodels.StatusPending,
	}).Error
	assert.NoError(t, err)

	// Seed customer
	err = db.Create(&barberbookingmodels.Customer{
		ID:       1,
		TenantID: 1,
		Email:    "test@example.com",
	}).Error
	assert.NoError(t, err)

	

	return db
}

func TestAppointmentStatusLogService(t *testing.T) {
	ctx := context.Background()
	db := setupAppointmentLogTestDB(t)
	svc := barberbookingservices.NewAppointmentStatusLogService(db)

	appointment := barberbookingmodels.Appointment{
		TenantID:   1,
		BranchID:   1,
		ServiceID:  1,
		CustomerID: 1,
		StartTime:  time.Now().Add(1 * time.Hour),
		EndTime:    time.Now().Add(2 * time.Hour),
		Status:     "PENDING",
	}
	assert.NoError(t, db.Create(&appointment).Error)
	appointmentID := appointment.ID

	t.Run("LogStatusChange_Success", func(t *testing.T) {
		err := svc.LogStatusChange(ctx, 1, "PENDING", "CONFIRMED", ptrUint(99), nil, "approved by admin")
		assert.NoError(t, err)
	})

	t.Run("GetLogsForAppointment_Success", func(t *testing.T) {
		logs, err := svc.GetLogsForAppointment(ctx, 1)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Equal(t, "PENDING", logs[0].OldStatus)
		assert.Equal(t, "CONFIRMED", logs[0].NewStatus)
	})

	t.Run("DeleteLogsByAppointmentID_Success", func(t *testing.T) {
		err := svc.DeleteLogsByAppointmentID(ctx, 1)
		assert.NoError(t, err)

		logs, err := svc.GetLogsForAppointment(ctx, 1)
		assert.NoError(t, err)
		assert.Empty(t, logs)
	})

	t.Run("GetLogsForAppointment_NotFound_ShouldReturnEmpty", func(t *testing.T) {
		logs, err := svc.GetLogsForAppointment(ctx, 9999) // สมมุติไม่มี
		assert.NoError(t, err)
		assert.Empty(t, logs)
	})

	t.Run("DeleteLogsByAppointmentID_WithoutLogs_ShouldNotFail", func(t *testing.T) {
		err := svc.DeleteLogsByAppointmentID(ctx, 123456) // appointment ที่ไม่มี log
		assert.NoError(t, err)
	})

	t.Run("LogStatusChange_WithNilUserAndCustomer_ShouldSucceed", func(t *testing.T) {
		
		err := svc.LogStatusChange(ctx, appointmentID, "PENDING", "CANCELLED", nil, nil, "system cancelled")
		assert.NoError(t, err)
	
		logs, err := svc.GetLogsForAppointment(ctx, appointmentID)
		assert.NoError(t, err)
		assert.Len(t, logs, 1)
		assert.Nil(t, logs[0].ChangedByUserID)
		assert.Nil(t, logs[0].ChangedByCustomerID)
	})

	t.Run("LogStatusChange_SameOldAndNewStatus_ShouldAllow", func(t *testing.T) {
		userID := uint(101)
		err := svc.LogStatusChange(ctx, appointmentID, "CONFIRMED", "CONFIRMED", &userID, nil, "no change but forced log")
		assert.NoError(t, err)
	
		logs, err := svc.GetLogsForAppointment(ctx, appointmentID)
		assert.NoError(t, err)
		assert.NotEmpty(t, logs)
	})

	t.Run("LogStatusChange_WithVeryLongNote_ShouldSucceed", func(t *testing.T) {
		longNote := strings.Repeat("a", 10000) // 10,000 ตัวอักษร
		err := svc.LogStatusChange(ctx, appointmentID, "PENDING", "CONFIRMED", nil, nil, longNote)
		assert.NoError(t, err)
	})
	
	
}
