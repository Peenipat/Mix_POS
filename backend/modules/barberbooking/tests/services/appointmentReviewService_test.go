package barberbookingServiceTest

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/joho/godotenv"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingServices "myapp/modules/barberbooking/services"
	coreModels "myapp/modules/core/models"
)

func setupTestReviewDB(t *testing.T) *gorm.DB {
	_ = godotenv.Load("../../../../.env.test")

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		t.Fatal("DATABASE_URL is not set")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("could not connect to test DB: %v", err)
	}

	if err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;").Error; err != nil {
		t.Fatalf("reset schema failed: %v", err)
	}

	if err := db.AutoMigrate(
		&barberBookingModels.Service{},
		&barberBookingModels.Customer{},
		&coreModels.Branch{},
		&barberBookingModels.Appointment{},
		&barberBookingModels.AppointmentReview{},
	); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}

	db.Create(&coreModels.Tenant{
		ID:   1,
		Name: "Tenant ทดสอบ",
	})
	// Seed branch
	db.Create(&coreModels.Branch{
		ID:       1,
		TenantID: 1,
		Name:     "สาขาทดสอบ",
	})

	// Seed customer
	db.Create(&barberBookingModels.Customer{
		ID:       1,
		Name:     "ลูกค้าทดสอบ",
		Email:    "review@example.com",
		TenantID: 1,
	})

	// Seed service
	db.Create(&barberBookingModels.Service{
		ID:       1,
		TenantID: 1,
		Name:     "ตัดผมชาย",
		Price:    200,
		Duration: 30,
	})

	// Seed appointment
	ap := barberBookingModels.Appointment{
		BranchID:   1,
		TenantID:   1,
		ServiceID:  1,
		CustomerID: 1,
		StartTime:  time.Now().Add(1 * time.Hour),
		EndTime:    time.Now().Add(1*time.Hour + 30*time.Minute),
		Status:     barberBookingModels.StatusComplete,
	}
	if err := db.Create(&ap).Error; err != nil {
		t.Fatalf("seed appointment failed: %v", err)
	}

	ap2 := barberBookingModels.Appointment{
		BranchID:   1,
		TenantID:   1,
		ServiceID:  1,
		CustomerID: 1,
		StartTime:  time.Now().Add(3 * time.Hour),
		EndTime:    time.Now().Add(3*time.Hour + 30*time.Minute),
		Status:     barberBookingModels.StatusComplete,
	}
	_ = db.Create(&ap2)

	review := &barberBookingModels.AppointmentReview{
		AppointmentID: ap2.ID,
		Rating:        3,
		Comment:       "It was fine",
	}
	_ = db.Create(review)

	return db
}

func TestAppointmentReviewService_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestReviewDB(t)
	svc := barberBookingServices.NewAppointmentReviewService(db)

	t.Run("CreateReview_Success", func(t *testing.T) {
		rev := &barberBookingModels.AppointmentReview{
			AppointmentID: 1,
			Rating:        5,
			Comment:       "Great service!",
		}
		got, err := svc.CreateReview(ctx, rev)
		assert.NoError(t, err)
		assert.NotZero(t, got.ID)
		assert.Equal(t, 5, got.Rating)
		assert.Equal(t, "Great service!", got.Comment)
		assert.Equal(t, uint(1), got.AppointmentID)
	})

	t.Run("CreateReview_MissingRequiredFields", func(t *testing.T) {
		rev := &barberBookingModels.AppointmentReview{} // ไม่มีข้อมูลใดๆ
		_, err := svc.CreateReview(ctx, rev)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid review input: appointmentID and rating (1-5) are required")
	})

	t.Run("CreateReview_InvalidAppointment", func(t *testing.T) {
		rev := &barberBookingModels.AppointmentReview{
			AppointmentID: 999, // ไม่มี appointment นี้
			Rating:        4,
			Comment:       "Okay",
		}
		_, err := svc.CreateReview(ctx, rev)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "appointment with ID 999 not found")
	})

	t.Run("GetReviewByID_Success", func(t *testing.T) {
		// สร้างรีวิวก่อน

		ap3 := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(5 * time.Hour),
			EndTime:    time.Now().Add(5*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap3)

		rev := barberBookingModels.AppointmentReview{
			ID:            uint(time.Now().Unix()),
			AppointmentID: ap3.ID,
			Rating:        3,
			Comment:       "It was fine",
		}
		if err := db.Create(&rev).Error; err != nil {
			t.Fatalf("seed review failed: %v", err)
		}

		got, err := svc.GetByID(ctx, rev.ID)
		assert.NoError(t, err)
		assert.Equal(t, rev.ID, got.ID)
		assert.Equal(t, "It was fine", got.Comment)
	})

	t.Run("GetReviewByID_NotFound", func(t *testing.T) {
		_, err := svc.GetByID(ctx, 9999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "review with ID 9999 not found")
	})

	t.Run("UpdateReview_Success", func(t *testing.T) {
		// สร้าง review ที่จะอัปเดต

		ap3 := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(6 * time.Hour),
			EndTime:    time.Now().Add(6*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap3)

		rev := &barberBookingModels.AppointmentReview{
			AppointmentID: ap3.ID,
			Rating:        3,
			Comment:       "It was fine",
		}
		err := db.Create(rev).Error
		assert.NoError(t, err)

		// Prepare input
		update := &barberBookingModels.AppointmentReview{
			Rating:  5,
			Comment: "เปลี่ยนใจแล้ว ดีมาก!",
		}

		updated, err := svc.UpdateReview(ctx, rev.ID, update)
		assert.NoError(t, err)
		assert.Equal(t, 5, updated.Rating)
		assert.Equal(t, "เปลี่ยนใจแล้ว ดีมาก!", updated.Comment)
	})

	t.Run("UpdateReview_NotFound", func(t *testing.T) {
		update := &barberBookingModels.AppointmentReview{
			Rating:  4,
			Comment: "ไม่พบรีวิวนี้",
		}

		_, err := svc.UpdateReview(ctx, 99999, update)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "review with ID 99999 not found")
	})

	t.Run("UpdateReview_InvalidRating", func(t *testing.T) {
		// สร้างรีวิวก่อน
		rev := &barberBookingModels.AppointmentReview{
			AppointmentID: 2,
			Rating:        3,
			Comment:       "เดี๋ยวจะใส่เรตผิด",
		}
		_ = db.Create(rev)

		update := &barberBookingModels.AppointmentReview{
			Rating:  6, // invalid
			Comment: "คะแนนเกิน",
		}

		_, err := svc.UpdateReview(ctx, rev.ID, update)
		assert.Error(t, err)
		assert.EqualError(t, err, "invalid rating: must be between 1 and 5")
	})

	t.Run("GetReviewByAppointment_Success", func(t *testing.T) {
		// เตรียม appointment ใหม่
		ap := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(7 * time.Hour),
			EndTime:    time.Now().Add(7*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap)

		// สร้าง review
		rev := &barberBookingModels.AppointmentReview{
			AppointmentID: ap.ID,
			Rating:        4,
			Comment:       "ดีมากครับ",
		}
		err := db.Create(rev).Error
		assert.NoError(t, err)

		// เรียก service
		got, err := svc.GetReviewByAppointment(ctx, ap.ID)
		assert.NoError(t, err)
		assert.Equal(t, rev.AppointmentID, got.AppointmentID)
		assert.Equal(t, rev.Comment, got.Comment)
		assert.Equal(t, 4, got.Rating)
	})

	t.Run("GetReviewByAppointment_NotFound", func(t *testing.T) {
		_, err := svc.GetReviewByAppointment(ctx, 99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "review for appointment ID 99999 not found")
	})

	t.Run("DeleteReview_ByOwner_Success", func(t *testing.T) {
		// สร้าง appointment ใหม่
		ap := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(8 * time.Hour),
			EndTime:    time.Now().Add(8*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap)

		// สร้าง review โดย CustomerID = 1
		review := barberBookingModels.AppointmentReview{
			AppointmentID: ap.ID,
			CustomerID:    ptrUint(1),
			Rating:        5,
			Comment:       "จะลบละ",
		}
		_ = db.Create(&review)

		err := svc.DeleteReview(ctx, review.ID, 1, string(coreModels.RoleNameUser))
		assert.NoError(t, err)

		// ตรวจสอบว่าลบแล้วจริง (soft delete)
		var check barberBookingModels.AppointmentReview
		tx := db.Unscoped().First(&check, review.ID)
		assert.NoError(t, tx.Error)
		assert.NotNil(t, check.DeletedAt.Valid)
	})

	t.Run("DeleteReview_ByOtherUser_Fail", func(t *testing.T) {
		// review เดิมยังค้างอยู่
		ap := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(8 * time.Hour),
			EndTime:    time.Now().Add(8*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap)
		review := barberBookingModels.AppointmentReview{
			AppointmentID: ap.ID,
			CustomerID:    ptrUint(1),
			Rating:        3,
			Comment:       "to be deleted",
		}
		err := db.Create(&review).Error
		assert.NoError(t, err) // ป้องกัน ID = 0

		// 🔥 ทดสอบลบโดย user ผิด
		err = svc.DeleteReview(ctx, review.ID, 9999, string(coreModels.RoleNameUser))
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not authorized")

	})

	t.Run("DeleteReview_ByAdmin_Success", func(t *testing.T) {
		ap := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(9 * time.Hour),
			EndTime:    time.Now().Add(9*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		_ = db.Create(&ap)

		review := barberBookingModels.AppointmentReview{
			AppointmentID: ap.ID,
			CustomerID:    ptrUint(1),
			Rating:        5,
			Comment:       "admin จะลบ",
		}
		_ = db.Create(&review)

		err := svc.DeleteReview(ctx, review.ID, 42, string(coreModels.RoleNameBranchAdmin)) // admin ลบ
		assert.NoError(t, err)
	})

	t.Run("DeleteReview_UnknownRole_Fail", func(t *testing.T) {

		ap := barberBookingModels.Appointment{
			BranchID:   1,
			TenantID:   1,
			ServiceID:  1,
			CustomerID: 1,
			StartTime:  time.Now().Add(6 * time.Hour),
			EndTime:    time.Now().Add(6*time.Hour + 30*time.Minute),
			Status:     barberBookingModels.StatusComplete,
		}
		err := db.Create(&ap).Error
		assert.NoError(t, err)

		review := barberBookingModels.AppointmentReview{
			AppointmentID: ap.ID, 
			CustomerID:    ptrUint(1),
			Rating:        5,
			Comment:       "unknown role จะลบ",
		}
		err = db.Create(&review).Error
		assert.NoError(t, err)

		err = svc.DeleteReview(ctx, review.ID, 1, "UNKNOWN_ROLE")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not authorized")
	})

	t.Run("GetAverageRatingByBarber_Success", func(t *testing.T) {
        barberID := uint(1)
    
        // Seed barber 
        barber := barberBookingModels.Barber{
            ID:       barberID,
            BranchID: 1,
            TenantID: 1,
            UserID:   101, // สมมุติ ID ผู้ใช้
        }
        err := db.Create(&barber).Error
        assert.NoError(t, err)
    
        // Appointment 1
        ap1 := barberBookingModels.Appointment{
            BranchID:   1,
            TenantID:   1,
            ServiceID:  1,
            CustomerID: 1,
            BarberID:   &barberID,
            StartTime:  time.Now().Add(1 * time.Hour),
            EndTime:    time.Now().Add(1*time.Hour + 30*time.Minute),
            Status:     barberBookingModels.StatusComplete,
        }
        err = db.Create(&ap1).Error
        assert.NoError(t, err)
        assert.NotZero(t, ap1.ID)
    
        //  Review 1
        rv1 := barberBookingModels.AppointmentReview{
            AppointmentID: ap1.ID,
            CustomerID:    ptrUint(1),
            Rating:        4,
            Comment:       "ดี",
        }
        assert.NoError(t, db.Create(&rv1).Error)
    
        //  Appointment 2
        ap2 := barberBookingModels.Appointment{
            BranchID:   1,
            TenantID:   1,
            ServiceID:  1,
            CustomerID: 1,
            BarberID:   &barberID,
            StartTime:  time.Now().Add(2 * time.Hour),
            EndTime:    time.Now().Add(2*time.Hour + 30*time.Minute),
            Status:     barberBookingModels.StatusComplete,
        }
        err = db.Create(&ap2).Error
        assert.NoError(t, err)
        assert.NotZero(t, ap2.ID)
    
        //  Review 2
        rv2 := barberBookingModels.AppointmentReview{
            AppointmentID: ap2.ID,
            CustomerID:    ptrUint(1),
            Rating:        5,
            Comment:       "เยี่ยม",
        }
        assert.NoError(t, db.Create(&rv2).Error)
    
        //  Run test
        avg, err := svc.GetAverageRatingByBarber(ctx, barberID)
        assert.NoError(t, err)
        assert.InEpsilon(t, 4.5, avg, 0.01) // (4+5)/2 = 4.5
    })
    

	t.Run("GetAverageRatingByBarber_NoReview", func(t *testing.T) {
		barberID := uint(99) // ไม่มีรีวิวเลย
		avg, err := svc.GetAverageRatingByBarber(ctx, barberID)
		assert.NoError(t, err)
		assert.Equal(t, 0.0, avg)
	})
}

