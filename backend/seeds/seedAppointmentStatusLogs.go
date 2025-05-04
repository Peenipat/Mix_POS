package seeds

import (
	"errors"
	"time"

	"gorm.io/gorm"

	bookingModels "myapp/modules/barberbooking/models"
)

func SeedAppointmentStatusLogs(db *gorm.DB) error {
	// 1. ดึง appointments ที่อยู่ในสถานะ 'PENDING' หรือ 'CONFIRMED'
	var appointments []bookingModels.Appointment
	if err := db.Where("status IN ?", []string{"PENDING", "CONFIRMED"}).
		Limit(5).Find(&appointments).Error; err != nil {
		return errors.New("failed to load appointments: " + err.Error())
	}

	if len(appointments) == 0 {
		return errors.New("no appointments found for status logging")
	}

	now := time.Now()

	for i, appt := range appointments {
		// เปลี่ยนเป็น COMPLETED เพื่อให้ seed review ได้ด้วย
		newStatus := bookingModels.StatusComplete

		// สร้าง log
		log := bookingModels.AppointmentStatusLog{
			AppointmentID:        appt.ID,
			OldStatus:            string(appt.Status),
			NewStatus:            string(newStatus),
			ChangedAt:            now.Add(-time.Duration(i) * time.Hour),
			ChangedByCustomerID:  &appt.CustomerID,
			Notes:                "เปลี่ยนสถานะทดสอบโดย seed",
		}

		// ข้ามถ้า log มีอยู่แล้ว
		var exists bookingModels.AppointmentStatusLog
		if err := db.Where("appointment_id = ?", appt.ID).First(&exists).Error; err == nil {
			continue
		}

		// 2. สร้าง log ใหม่
		if err := db.Create(&log).Error; err != nil {
			return err
		}

		// 3. อัปเดตสถานะในตาราง appointments ให้ตรงกัน
		if err := db.Model(&appt).Update("status", newStatus).Error; err != nil {
			return err
		}
	}

	return nil
}

