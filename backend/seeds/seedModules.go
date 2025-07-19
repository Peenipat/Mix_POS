package seeds

import (
	"time"

	"gorm.io/gorm"
	coreModels "myapp/modules/core/models"
)

func SeedModules(db *gorm.DB) error {
	modules := []coreModels.Module{
		{Name: "barber_booking", Description: "ระบบจองคิวตัดผม"},
		{Name: "pos", Description: "ระบบขายหน้าร้าน"},
		{Name: "inventory", Description: "ระบบจัดการสต๊อก"},
	}

	now := time.Now()

	for _, m := range modules {
		record := coreModels.Module{Name: m.Name}
		attrs := coreModels.Module{
			Description: m.Description,
			CreatedAt:   now,
		}
		if err := db.Where(record).Assign(attrs).FirstOrCreate(&record).Error; err != nil {
			return err
		}
	}

	return nil
}
