package barberBookingModels

import (
	"time"
	"gorm.io/gorm"
)


type Customer struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	TenantID  uint      `gorm:"not null;index:idx_tenant_email,priority:1"` // Composite Index
	Name      string    `gorm:"not null"`
	Phone     string    `gorm:"type:text"`    // optional
	Email     string    `gorm:"type:text;index:idx_tenant_email,priority:2" json:"email"` // Composite Index
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}