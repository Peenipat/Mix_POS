package seeds

import (
    "errors"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedTenantUsers สร้างความสัมพันธ์ many-to-many ระหว่าง Tenant กับ User
// ต้องรันหลังจาก SeedTenants() และ SeedUsers() เรียบร้อยแล้ว
func SeedTenantUsers(db *gorm.DB) error {
    // 1) ดึง Default Tenant
    var tenant coreModels.Tenant
    if err := db.Where("domain = ?", "default.example.com").
        First(&tenant).Error; err != nil {
        return errors.New("default tenant not found: " + err.Error())
    }

    // 2) ดึงผู้ใช้ทุก Role ที่ควรแมปกับ tenant นี้ (ยกตัวอย่างเอาทุกคนที่มี branch_id == nil)
    var users []coreModels.User
    if err := db.Where("branch_id = 1").Find(&users).Error; err != nil {
        return err
    }

    // 3) สร้างหรืออัปเดตแถวใน join-table tenant_users
    for _, u := range users {
        tu := coreModels.TenantUser{
            TenantID: tenant.ID,
            UserID:   u.ID,
        }
        if err := db.FirstOrCreate(&tu, tu).Error; err != nil {
            return err
        }
    }
    return nil
}
