// pkg/seeds/seed_working_hours.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    "gorm.io/gorm/clause"

    coreModels    "myapp/modules/core/models"
    bookingModels "myapp/modules/barberbooking/models"
)

// SeedWorkingHours สร้างตารางเวลาทำงาน (WorkingHour) ให้ Default Branch
func SeedWorkingHours(db *gorm.DB) error {
    // โหลด Default Branch
    var branch coreModels.Branch
    if err := db.Where("name = ?", "Branch 1").First(&branch).Error; err != nil {
        return err
    }

    // กำหนดวันจันทร์–ศุกร์
    weekdays := []int{1, 2, 3, 4, 5}

    // เวลา 9:00–17:00 (ใช้เฉพาะ Time component)
    now := time.Now()
    startTime := time.Date(now.Year(), now.Month(), now.Day(), 9, 0, 0, 0, now.Location())
    endTime   := time.Date(now.Year(), now.Month(), now.Day(), 17, 0, 0, 0, now.Location())

    for _, wd := range weekdays {
        wh := bookingModels.WorkingHour{
            BranchID:  branch.ID,
            Weekday:   wd,
            StartTime: startTime,
            EndTime:   endTime,
        }
        // OnConflict: ถ้า (branch_id, weekday) ซ้ำ ให้อัปเดต start_time, end_time
        if err := db.Clauses(clause.OnConflict{
            Columns:   []clause.Column{{Name: "branch_id"}, {Name: "weekday"}},
            DoUpdates: clause.AssignmentColumns([]string{"start_time", "end_time"}),
        }).Create(&wh).Error; err != nil {
            return err
        }
    }
    return nil
}
