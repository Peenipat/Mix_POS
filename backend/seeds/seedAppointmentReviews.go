package seeds

import (
    "errors"
    "time"

    "gorm.io/gorm"
    models "myapp/modules/barberbooking/models"
)

func SeedAppointmentReviews(db *gorm.DB) error {
    reviews := []models.AppointmentReview{
        // … ตัวอย่าง reviews ที่จะ seed
    }

    for _, r := range reviews {
        // ตรวจสอบก่อนว่า appointment นั้นมีอยู่จริงหรือไม่
        var appt models.Appointment
        if err := db.First(&appt, r.AppointmentID).Error; err != nil {
            if errors.Is(err, gorm.ErrRecordNotFound) {
                // ไม่มี appointment → ข้ามเคสนี้
                continue
            }
            return err
        }

        // ลองดูว่ามี review เดิมหรือยัง
        var existing models.AppointmentReview
        err := db.Where("appointment_id = ?", r.AppointmentID).
            First(&existing).Error

        now := time.Now()
        r.CreatedAt = now
        r.UpdatedAt = now

        if errors.Is(err, gorm.ErrRecordNotFound) {
            // ยังไม่มี → insert ใหม่
            if err := db.Create(&r).Error; err != nil {
                return err
            }
        } else if err != nil {
            return err
        } else {
            // มีอยู่แล้ว → update fields แล้ว save
            existing.Rating = r.Rating
            existing.Comment = r.Comment
            existing.UpdatedAt = now
            if err := db.Save(&existing).Error; err != nil {
                return err
            }
        }
    }
    return nil
}
