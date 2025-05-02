package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

func SeedRoles(db *gorm.DB) error {
    // 1) หา module CORE
    var coreMod coreModels.Module
    if err := db.Where("key = ?", "CORE").First(&coreMod).Error; err != nil {
        return err
    }

    // 2) รายชื่อ role ที่จะ seed
    roles := []coreModels.Role{
        {
            ModuleID:    coreMod.ID,
            Name:        coreModels.RoleNameSaaSSuperAdmin,
            Description: "ควบคุมระบบ SaaS ทั้งหมด",
        },
        {
            ModuleID:    coreMod.ID,
            Name:        coreModels.RoleNameTenantAdmin,
            Description: "หัวหน้าผู้เช่า ดูข้อมูลทุกสาขาของตน",
        },
        {
            ModuleID:    coreMod.ID,
            Name:        coreModels.RoleNameBranchAdmin,
            Description: "หัวหน้าสาขา แต่ละร้าน",
        },
        {
            ModuleID:    coreMod.ID,
            Name:        coreModels.RoleNameAssistantManager,
            Description: "รองหัวหน้าสาขา",
        },
        {
            ModuleID:    coreMod.ID,
            Name:        coreModels.RoleNameStaff,
            Description: "พนักงานสาขา",
        },
        {
            ModuleID:    coreMod.ID,
            Name:        coreModels.RoleNameUser,
            Description: "ผู้ใช้ทั่วไป (ยังไม่เช่า)",
        },
    }

    now := time.Now()
    // 3) Loop สร้างหรืออัปเดต role ตาม key + name
    for _, r := range roles {
        record := coreModels.Role{
            ModuleID: r.ModuleID,
            Name:     r.Name,
        }
        attrs := coreModels.Role{
            Description: r.Description,
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