package barberBookingModels

import "time"

type Customer struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"type:text;not null" json:"name"`
	Phone     string    `gorm:"type:text" json:"phone,omitempty"`
	Email     string    `gorm:"type:text" json:"email,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}