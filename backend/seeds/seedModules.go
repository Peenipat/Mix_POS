// seeds/seedModules.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedModules ต้องรันก่อน SeedRoles
func SeedModules(db *gorm.DB) error {
    now := time.Now()
    modules := []coreModels.Modules{
        {Key: "CORE",        Description: "Core system",        CreatedAt: now, UpdatedAt: now},
        {Key: "BOOKING",     Description: "Appointment booking",CreatedAt: now, UpdatedAt: now},
        // เพิ่มโมดูลอื่น ๆ ตามต้องการ
    }

    for _, m := range modules {
        record := coreModels.Modules{Key: m.Key}
        attrs  := coreModels.Modules{
            Description: m.Description,
            UpdatedAt:   now,
        }
        if err := db.
            Where("key = ?", m.Key).
            Assign(attrs).
            FirstOrCreate(&record, record).
            Error; err != nil {
            return err
        }
    }
    return nil
}
