// seeds/seedRoles.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedRoles ต้องรันหลัง SeedModules
func SeedRoles(db *gorm.DB) error {
    // 1) หา module CORE ก่อน
    var coreMod coreModels.Modules
    if err := db.Where("key = ?", "CORE").First(&coreMod).Error; err != nil {
        return err
    }

    // 2) รายชื่อ role ที่จะ seed
    now := time.Now()
    roles := []coreModels.Role{
        {ModuleID: coreMod.ID, Name: coreModels.RoleNameSaaSSuperAdmin,   Description: "ควบคุมระบบ SaaS ทั้งหมด",        CreatedAt: now, UpdatedAt: now},
        {ModuleID: coreMod.ID, Name: coreModels.RoleNameTenantAdmin,    Description: "หัวหน้าผู้เช่า ดูข้อมูลทุกสาขาของตน", CreatedAt: now, UpdatedAt: now},
        {ModuleID: coreMod.ID, Name: coreModels.RoleNameBranchAdmin,    Description: "หัวหน้าสาขา แต่ละร้าน",             CreatedAt: now, UpdatedAt: now},
        {ModuleID: coreMod.ID, Name: coreModels.RoleNameAssistantManager,Description: "รองหัวหน้าสาขา",                 CreatedAt: now, UpdatedAt: now},
        {ModuleID: coreMod.ID, Name: coreModels.RoleNameStaff,          Description: "พนักงานสาขา",                  CreatedAt: now, UpdatedAt: now},
        {ModuleID: coreMod.ID, Name: coreModels.RoleNameUser,           Description: "ผู้ใช้ทั่วไป (ยังไม่เช่า)",        CreatedAt: now, UpdatedAt: now},
    }

    // 3) Loop สร้างหรืออัปเดต role ตาม ModuleID + Name
    for _, r := range roles {
        // ใช้ FirstOrCreate แบบ composite key (module_id + name)
        record := coreModels.Role{
            ModuleID: r.ModuleID,
            Name:     r.Name,
        }
        attrs := coreModels.Role{
            Description: r.Description,
            UpdatedAt:   now,
        }
        if err := db.
            Where("module_id = ? AND name = ?", r.ModuleID, r.Name).
            Assign(attrs).
            FirstOrCreate(&record, record).
            Error; err != nil {
            return err
        }
    }
    return nil
}
