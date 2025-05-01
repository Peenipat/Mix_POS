// pkg/seeds/seed_appointments.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
    bookingModels "myapp/modules/barberbooking/models"
)

// SeedAppointments สร้างตัวอย่างการนัดหมาย (Appointment)
// เรียกหลังจาก SeedBranches(), SeedServices(), SeedBarbers() และ SeedUsers()
func SeedAppointments(db *gorm.DB) error {
    // 1) หา Default Branch, Service, Barber, Customer
    var (
        branch   coreModels.Branch
        service  bookingModels.Service
        barber   bookingModels.Barber
        customer coreModels.User
    )
    if err := db.Where("name = ?", "Default Branch").First(&branch).Error; err != nil {
        return err
    }
    if err := db.Where("name = ?", "Haircut").First(&service).Error; err != nil {
        return err
    }
    if err := db.Where("branch_id = ?", branch.ID).First(&barber).Error; err != nil {
        return err
    }
    // เปลี่ยนมาใช้ email ที่ seed จริง
    if err := db.Where("email = ?", "user@default.example.com").First(&customer).Error; err != nil {
        return err
    }

    // 2) กำหนดเวลานัดสำหรับวันพรุ่งนี้ 10:00–(10:00 + duration นาที)
    tomorrow := time.Now().Add(24 * time.Hour)
    start := time.Date(
        tomorrow.Year(), tomorrow.Month(), tomorrow.Day(),
        10, 0, 0, 0,
        tomorrow.Location(),
    )
    // แปลง service.Duration (นาที) → time.Duration
    end := start.Add(time.Duration(service.Duration) * time.Minute)

    // 3) สร้างหรืออัปเดต Appointment
    record := bookingModels.Appointment{
        BranchID:   branch.ID,
        ServiceID:  service.ID,
        BarberID:   &barber.ID,
        CustomerID: customer.ID,
        StartTime:  start,
    }
    attrs := bookingModels.Appointment{
        EndTime: end,
        Status:  "PENDING",
        Notes:   "Seeded test appt",
    }
    return db.Where(record).
        Assign(attrs).
        FirstOrCreate(&record).Error
}
