package barberBookingModels

import (
	"time"
	 coreModels "myapp/modules/core/models"
	"gorm.io/gorm"
)

type Barber struct {
	ID uint `gorm:"primaryKey" json:"id"`

	BranchID 		uint 			`gorm:"not null;index" json:"branch_id"`
	UserID   		uint 			`gorm:"not null;uniqueIndex" json:"user_id"`
	User   			coreModels.User `json:"user"`
	TenantID 		uint 			`gorm:"not null;" json:"tenant_id" `
	RoleUser     	string      	`gorm:"type:varchar(100);" json:"role_user"`
	Description     string      	`gorm:"type:varchar(100);not null" json:"description"`

	CreatedAt 		time.Time      `json:"created_at"`
	UpdatedAt 		time.Time      `json:"updated_at"`
	DeletedAt 		gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// ❗ DO NOT preload User/Branch struct (loose coupling)
	// User   coreModels.User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
	// Branch coreModels.Branch `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
}
