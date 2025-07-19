// pkg/seeds/seed_working_hours.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    coreModels    "myapp/modules/core/models"
    bookingModels "myapp/modules/barberbooking/models"
)
func SeedWorkingHours(db *gorm.DB) error {
	// โหลด Branch ตัวอย่าง
	var branch coreModels.Branch
	if err := db.Where("name = ?", "Branch 1").First(&branch).Error; err != nil {
		return err
	}

	weekdays := []int{0, 1, 2, 3, 4, 5, 6}
	now := time.Now()
	loc := now.Location()
	startTime := time.Date(0, 1, 1, 9, 0, 0, 0, loc)
	endTime := time.Date(0, 1, 1, 17, 0, 0, 0, loc)

	for _, wd := range weekdays {
		wh := bookingModels.WorkingHour{
			BranchID: branch.ID,
			TenantID: 1,
			Weekday:  wd,
			IsClosed: wd == 0 || wd == 6, // 
		}

		if !wh.IsClosed {
			wh.StartTime = startTime
			wh.EndTime = endTime
		}

		if err := db.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "branch_id"}, {Name: "weekday"}},
			DoUpdates: clause.AssignmentColumns([]string{"start_time", "end_time", "is_closed"}),
		}).Create(&wh).Error; err != nil {
			return err
		}
	}

	return nil
}
