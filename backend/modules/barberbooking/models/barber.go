package barberBookingModels

import (
	"time"

	"gorm.io/gorm"
)

type Barber struct {
	ID uint `gorm:"primaryKey" json:"id"`

	BranchID 		uint 			`gorm:"not null;index" json:"branch_id"`
	UserID   		uint 			`gorm:"not null;uniqueIndex" json:"user_id"`
	TenantID 		uint 			`gorm:"not null;" json:"tenant_id" `
	PhoneNumber 	string      	`gorm:"type:varchar(20);not null" json:"phone_number"`
	Description     string      	`gorm:"type:varchar(100);not null" json:"description"`

	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// ❗ DO NOT preload User/Branch struct (loose coupling)
	// User   coreModels.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// Branch coreModels.Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
