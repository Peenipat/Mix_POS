package seeds

import (
    "time"

    "gorm.io/gorm"
    bookingModels "myapp/modules/barberbooking/models"
)

// SeedServices สร้างรายการบริการ (services) ตั้งต้น
func SeedServices(db *gorm.DB) error {
    items := []bookingModels.Service{
        {Name: "Haircut",    Duration: 30, Price: 200, TenantID:	1,},
        {Name: "Shampoo",    Duration: 15, Price: 100, TenantID:	1,},
        {Name: "Beard Trim", Duration: 20, Price: 150, TenantID:	1,},
        // เพิ่มรายการบริการอื่นได้ที่นี่…
    }

    now := time.Now()
    for _, svc := range items {
        record := bookingModels.Service{Name: svc.Name}
        attrs  := bookingModels.Service{
            Duration:  svc.Duration,  // นาที
            Price:     svc.Price,
            TenantID: svc.TenantID,
            UpdatedAt: now,

        }
        if err := db.Where(record).
            Assign(attrs).
            FirstOrCreate(&record).Error; err != nil {
            return err
        }
    }
    return nil
}
