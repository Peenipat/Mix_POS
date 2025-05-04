package seeds

import (
	"time"

	"gorm.io/gorm"
	bookingModels "myapp/modules/barberbooking/models"
)

func SeedCustomers(db *gorm.DB) error {
	customers := []bookingModels.Customer{
		{
			Name:  "สมชาย ใจดี",
			Phone: "0801234567",
			Email: "somchai@example.com",
		},
		{
			Name:  "Jane Doe",
			Phone: "0912345678",
			Email: "jane@example.com",
		},
		{
			Name:  "John Smith",
			Phone: "0899999999",
			Email: "john@example.com",
		},
	}

	now := time.Now()
	for _, c := range customers {
		record := bookingModels.Customer{
			Email: c.Email,
		}
		attrs := bookingModels.Customer{
			Name:      c.Name,
			Phone:     c.Phone,
			CreatedAt: now,
		}
		if err := db.Where(record).Assign(attrs).FirstOrCreate(&record).Error; err != nil {
			return err
		}
	}
	return nil
}
