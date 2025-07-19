// pkg/seeds/seed_appointments.go
package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
    bookingModels "myapp/modules/barberbooking/models"
)

// SeedAppointments สร้างตัวอย่างการนัดหมาย (Appointment)
// เรียกหลังจาก SeedBranches(), SeedServices(), SeedBarbers() และ SeedUsers()
func SeedAppointments(db *gorm.DB) error {
	var (
		branch   coreModels.Branch
		service  bookingModels.Service
		barber   bookingModels.Barber
		customer bookingModels.Customer
	)

	// 1) โหลดข้อมูลที่เกี่ยวข้อง
	if err := db.Where("name = ?", "Branch 1").First(&branch).Error; err != nil {
		return err
	}
	if err := db.Where("name = ?", "Haircut").First(&service).Error; err != nil {
		return err
	}
	if err := db.Where("branch_id = ?", branch.ID).First(&barber).Error; err != nil {
		return err
	}
	if err := db.Where("email = ?", "somchai@example.com").First(&customer).Error; err != nil {
		return err
	}

	// 2) เตรียมเวลานัดหมาย
	tomorrow := time.Now().Add(24 * time.Hour)
	loc := time.Local

	appointments := []bookingModels.Appointment{
		{
			BranchID:   branch.ID,
			ServiceID:  service.ID,
			BarberID:   barber.ID,
			CustomerID: customer.ID,
			TenantID:	1,
			StartTime:  time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 0, 0, 0, loc),
			EndTime:    time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 30, 0, 0, loc),
			Status:     bookingModels.StatusPending,
			Notes:      "Seeded pending appt",
		},
		{
			BranchID:   branch.ID,
			ServiceID:  service.ID,
			BarberID:   barber.ID,
			CustomerID: customer.ID,
			TenantID:	1,
			StartTime:  time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 10, 30, 0, 0, loc),
			EndTime:    time.Date(tomorrow.Year(), tomorrow.Month(), tomorrow.Day(), 11, 0, 0, 0, loc),
			Status:     bookingModels.StatusComplete,
			Notes:      "Seeded completed appt",
		},
	}

	for _, appt := range appointments {
		record := bookingModels.Appointment{
			BranchID:   appt.BranchID,
			ServiceID:  appt.ServiceID,
			BarberID:   appt.BarberID,
			CustomerID: appt.CustomerID,
			TenantID:	appt.TenantID,
			StartTime:  appt.StartTime,
		}
		if err := db.Where(record).Assign(appt).FirstOrCreate(&record).Error; err != nil {
			return err
		}
	}

	return nil
}


