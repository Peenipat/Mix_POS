package seeds

import (
    "time"
    "fmt"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
)

// SeedRoles สร้างรายการบทบาท (roles) ตั้งต้น ใช้ FirstOrCreate เพื่อป้องกัน insert ซ้ำ
func SeedRoles(db *gorm.DB) error {
	// หา tenant ที่ต้องการ
	var defaultTenant coreModels.Tenant
	if err := db.Where("domain = ?", "default.example.com").First(&defaultTenant).Error; err != nil {
		return fmt.Errorf("cannot find default tenant: %w", err)
	}

	moduleName := "barber_booking"
	roles := []coreModels.Role{
        {
            TenantID:    nil,
            ModuleName:  nil,
            Name:        string(coreModels.RoleNameSaaSSuperAdmin),
            Description: "ผู้ดูแลระบบ SaaS ทั้งหมด",
        },
		{
			TenantID:    &defaultTenant.ID,
			ModuleName:  &moduleName,
			Name:        string(coreModels.RoleNameBranchAdmin),
			Description: "หัวหน้าสาขา มีสิทธิ์จัดการข้อมูลในระบบจองคิวตัดผม",
		},
		{
			TenantID:    &defaultTenant.ID,
			ModuleName:  &moduleName,
			Name:        string(coreModels.RoleNameAssistantManager),
			Description: "รองหัวหน้า จัดการคิวและดูรายงาน",
		},
		{
			TenantID:    &defaultTenant.ID,
			ModuleName:  &moduleName,
			Name:        string(coreModels.RoleNameStaff),
			Description: "พนักงานประจำร้าน ดูคิว และแจ้งสถานะ",
		},
		// ตัวอย่าง role ทั่วไป
		{
			TenantID:    &defaultTenant.ID,
			ModuleName:  nil,
			Name:        string(coreModels.RoleNameTenantAdmin),
			Description: "ผู้ดูแลร้านค้า สามารถจัดการผู้ใช้และสาขา",
		},
		{
			TenantID:    &defaultTenant.ID,
			ModuleName:  nil,
			Name:        string(coreModels.RoleNameUser),
			Description: "ผู้ใช้งานทั่วไป",
		},
	}

	now := time.Now()
	for _, r := range roles {
		record := coreModels.Role{
			Name:       r.Name,
			TenantID:   r.TenantID,
			ModuleName: r.ModuleName,
		}
		attrs := coreModels.Role{
			Description: r.Description,
			CreatedAt:   now,
			UpdatedAt:   now,
		}
		if err := db.Where(record).Assign(attrs).FirstOrCreate(&record).Error; err != nil {
			return err
		}
	}

	return nil
}
