// pkg/seeds/seed_barbers.go
package seeds

import (
    "errors"

    "gorm.io/gorm"
    coreModels    "myapp/modules/core/models"
    bookingModels "myapp/modules/barberbooking/models"
)

// SeedBarbers สร้างรายการช่าง (Barber) ให้กับ Default Branch
func SeedBarbers(db *gorm.DB) error {
    // 1) หา Default Branch
    var branch coreModels.Branch
    if err := db.Where("name = ?", "Branch 1").First(&branch).Error; err != nil {
        return errors.New("default branch not found: " + err.Error())
    }

    // 2) หา Users ที่เราต้องการให้เป็นช่าง (ตัวอย่าง: assistant_mgr + staff_user)
    emails := []string{
        "assistant@gmail.com",
        "staff@gmail.com",
    }
    var users []coreModels.User
    if err := db.Where("email IN ?", emails).Find(&users).Error; err != nil {
        return err
    }
    if len(users) == 0 {
        return errors.New("no users found for barber seeding")
    }

    // 3) สร้างหรืออัปเดต Barber record
    for _, u := range users {
        record := bookingModels.Barber{
            BranchID: branch.ID,
            UserID:   u.ID,
            TenantID: branch.TenantID,
            Description: "ช่าง",
        }
        if err := db.FirstOrCreate(&record, record).Error; err != nil {
            return err
        }
    }
    return nil
}
