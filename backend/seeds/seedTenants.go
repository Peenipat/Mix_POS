// pkg/seeds/seedTenants.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedTenants สร้างข้อมูล Tenant ตั้งต้น
// ตัวอย่างนี้เรา seed “Default Tenant” คุณสามารถเพิ่มรายการอื่นได้ตามต้องการ
func SeedTenants(db *gorm.DB) error {
    tenants := []coreModels.Tenant{
        {
            Name:     "Default Tenant",
            Domain:   "default.example.com",
            IsActive: true,
        },
        // ถ้าต้องการ seed tenant เพิ่มเติมให้เพิ่มเข้ามาที่นี่
        // { Name: "Another Tenant", Domain: "another.example.com", IsActive: true },
    }

    now := time.Now()
    for _, t := range tenants {
        record := coreModels.Tenant{Domain: t.Domain}
        attrs  := coreModels.Tenant{
            Name:      t.Name,
            IsActive:  t.IsActive,
            CreatedAt: now,
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
