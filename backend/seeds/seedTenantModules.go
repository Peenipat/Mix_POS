package seeds

import (
	"errors"
	"gorm.io/gorm"

	coreModels "myapp/modules/core/models"
)

func SeedTenantModules(db *gorm.DB) error {
	// 1) หา tenant ที่ต้องการ (เช่น Default Tenant)
	var tenant coreModels.Tenant
	if err := db.Where("domain = ?", "default.example.com").First(&tenant).Error; err != nil {
		return errors.New("default tenant not found: " + err.Error())
	}

	// 2) หา modules ที่เราจะผูกกับ tenant นี้
	var modules []coreModels.Module
	if err := db.Where("name IN ?", []string{"barber_booking", "pos"}).Find(&modules).Error; err != nil {
		return errors.New("failed to load modules: " + err.Error())
	}
	if len(modules) == 0 {
		return errors.New("no modules found to assign")
	}

	// 3) สร้าง tenant_modules
	for _, m := range modules {
		record := coreModels.TenantModule{
			TenantID: tenant.ID,
			ModuleID: m.ID,
		}
		if err := db.FirstOrCreate(&record, record).Error; err != nil {
			return err
		}
	}

	return nil
}
