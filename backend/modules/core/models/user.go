package coreModels

import (
	"time"
		"gorm.io/gorm"
)
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"type:text;not null" json:"username"`
	Email       string         `gorm:"uniqueIndex" json:"email"`
	Password    string         `gorm:"not null" json:"-"` // ซ่อนไม่ให้แสดงออกไป

	RoleID      uint           `gorm:"not null" json:"role_id"`
	Role        Role           `gorm:"foreignKey:RoleID" json:"role"`

	BranchID    *uint          `gorm:"index" json:"branch_id,omitempty"`
	Branch      *Branch        `gorm:"foreignKey:BranchID" json:"branch,omitempty"`

	TenantUsers []TenantUser   `gorm:"foreignKey:UserID" json:"tenant_users,omitempty"`

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
