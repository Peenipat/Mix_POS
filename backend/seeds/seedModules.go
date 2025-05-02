package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

func SeedModules(db *gorm.DB) error {
    modules := []coreModels.Module{
        {Key: "CORE",           Description: "Core system (auth, users, roles)"},
        {Key: "BOOKING",        Description: "Booking module"},
        {Key: "POS_RESTAURANT", Description: "Restaurant POS module"},
        // ... เพิ่มโมดูลอื่นๆ ตามต้องการ
    }

    now := time.Now()
    for _, m := range modules {
        // ถ้ายังไม่มี key นี้ให้สร้าง (set timestamps)
        record := coreModels.Module{Key: m.Key}
        attrs  := coreModels.Module{
            Description: m.Description,
            CreatedAt:   now,
            UpdatedAt:   now,
        }
        if err := db.
            Where(&record).
            Assign(attrs).
            FirstOrCreate(&record).
            Error; err != nil {
            return err
        }
    }
    return nil
}

