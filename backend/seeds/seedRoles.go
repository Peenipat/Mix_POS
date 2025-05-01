package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedRoles สร้างรายการบทบาท (roles) ตั้งต้น ใช้ FirstOrCreate เพื่อป้องกัน insert ซ้ำ
func SeedRoles(db *gorm.DB) error {
    // กำหนดค่าบทบาทที่ต้องการ seed
    roles := []coreModels.Role{
        {Name: string(coreModels.RoleNameSaaSSuperAdmin),   Description: "ควบคุมระบบ SaaS ทั้งหมด"},
        {Name: string(coreModels.RoleNameTenantAdmin),      Description: "หัวหน้าผู้เช่า ดูข้อมูลทุกสาขาของตน"},
        {Name: string(coreModels.RoleNameBranchAdmin),      Description: "หัวหน้าสาขา แต่ละร้าน"},
        {Name: string(coreModels.RoleNameAssistantManager), Description: "รองหัวหน้าสาขา"},
        {Name: string(coreModels.RoleNameStaff),            Description: "พนักงานสาขา"},
        {Name: string(coreModels.RoleNameUser),             Description: "ผู้ใช้ทั่วไป (ยังไม่เช่า)"},
    }

    now := time.Now()
    for _, r := range roles {
        // FirstOrCreate: ถ้ายังไม่มีแถวที่ตรง with Role.Name ก็สร้างใหม่ (set timestamps)
        record := coreModels.Role{Name: r.Name}
        attrs  := coreModels.Role{
            Description: r.Description,
            CreatedAt:   now,
            UpdatedAt:   now,
        }
        if err := db.Where(record).
            Assign(attrs).
            FirstOrCreate(&record).Error; err != nil {
            return err
        }
    }
    return nil
}
