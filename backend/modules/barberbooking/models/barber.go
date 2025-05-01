package bookingModels

import (
	"time"

	"gorm.io/gorm"
	"myapp/modules/core/models"
)
// Barber แทนช่างในแต่ละสาขา
// ผูกกับ User ผ่าน models.User.ID
// ผูกกับ Branch ผ่าน models.Branch.ID
type Barber struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BranchID  uint           `gorm:"not null;index" json:"branch_id"`
	Branch    coreModels.Branch  `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	UserID    uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	User      coreModels.User    `gorm:"foreignKey:UserID" json:"user,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
