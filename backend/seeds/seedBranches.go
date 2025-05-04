package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedBranches สร้างข้อมูลสาขา (branches) ตั้งต้น
func SeedBranches(db *gorm.DB) error {
    // สมมติว่าเราต้องการสาขาหลักของ Default Tenant
    var tenant coreModels.Tenant
    if err := db.Where("domain = ?", "default.example.com").First(&tenant).Error; err != nil {
        return err
    }

    branches := []coreModels.Branch{
        {
            Name:     "Default Branch",
            TenantID: tenant.ID,
        },
        // เติมสาขาอื่น ๆ ได้ที่นี่
    }

    now := time.Now()
    for _, b := range branches {
        record := coreModels.Branch{
            TenantID: b.TenantID,
            Name:     b.Name,
        }
        attrs := coreModels.Branch{
            CreatedAt: now,
            UpdatedAt: now,
        }
        
        if err := db.Unscoped().Where(record).Assign(attrs).FirstOrCreate(&record).Error; err != nil {
            return err
        }
        
    }
    return nil
}
