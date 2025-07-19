package seeds

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	bookingModels "myapp/modules/barberbooking/models"
)

// SeedBarberWorkloads สร้างข้อมูล workload ของ barber ย้อนหลัง 3 วัน
func SeedBarberWorkloads(db *gorm.DB) error {
	var barbers []bookingModels.Barber
	if err := db.Find(&barbers).Error; err != nil {
		return err
	}
	if len(barbers) == 0 {
		return errors.New("no barbers found for workload seeding")
	}

	today := time.Now().Truncate(24 * time.Hour)
	for _, b := range barbers {
		for i := 0; i < 3; i++ { // ย้อนหลัง 3 วัน
			date := today.AddDate(0, 0, -i)
			workload := bookingModels.BarberWorkload{
				BarberID:          b.ID,
				Date:              date,
				TotalAppointments: 3 + i,   // จำลองจำนวนนัด
				TotalHours:        2 + i,   // จำลองจำนวนชั่วโมง
				CreatedAt:         time.Now(),
			}
			if err := db.Clauses(clause.OnConflict{
				Columns:   []clause.Column{{Name: "barber_id"}, {Name: "date"}},
				DoNothing: true,
			}).Create(&workload).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
