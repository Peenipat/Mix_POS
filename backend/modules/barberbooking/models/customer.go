package barberBookingModels

import "time"


type Customer struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	TenantID  uint      `gorm:"not null;index:idx_tenant_email,priority:1"` // Composite Index
	Name      string    `gorm:"not null"`
	Phone     string    `gorm:"type:text"`    // optional
	Email     string    `gorm:"type:text;index:idx_tenant_email,priority:2"` // Composite Index
	CreatedAt time.Time `gorm:"autoCreateTime"`
}