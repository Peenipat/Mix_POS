package seeds

import (
	"errors"
	"time"

	"gorm.io/gorm"

	bookingModels "myapp/modules/barberbooking/models"
)

func SeedAppointmentReviews(db *gorm.DB) error {
	// 1. โหลด Appointment ที่มีสถานะ Completed
	var appointments []bookingModels.Appointment
	if err := db.Where("status = ?", bookingModels.StatusComplete).Limit(5).Find(&appointments).Error; err != nil {
		return errors.New("failed to find appointments: " + err.Error())
	}
	if len(appointments) == 0 {
		return errors.New("no completed appointments found")
	}

	// 2. สร้าง Seed Data ของ Review
	for i, appt := range appointments {
		review := bookingModels.AppointmentReview{
			AppointmentID: appt.ID,
			CustomerID:    &appt.CustomerID,
			Rating:        4 + i%2, // สลับระหว่าง 4 และ 5
			Comment:       "ขอบคุณสำหรับบริการดี ๆ",
			CreatedAt:     time.Now(),
		}

		// ตรวจสอบว่ามี review อยู่แล้วหรือยัง
		var existing bookingModels.AppointmentReview
		if err := db.Where("appointment_id = ?", appt.ID).First(&existing).Error; err == nil {
			continue // มีอยู่แล้ว → ข้าม
		}

		if err := db.Create(&review).Error; err != nil {
			return err
		}
	}

	return nil
}
