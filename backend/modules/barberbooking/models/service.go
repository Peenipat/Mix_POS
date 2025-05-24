package barberBookingModels

import (
	"time"

	"gorm.io/gorm"
)

// Service แทนบริการต่างๆ (เช่น ตัดผม สระผม ไดร์)
type Service struct {
	ID          uint           `gorm:"primaryKey;autoIncrement"`
	Name        string         `gorm:"type:varchar(100);not null" json:"name"`
	TenantID    uint           `gorm:"not null;index" json:"tenant_id"`
	Duration    int            `gorm:"not null" json:"duration"`    // ระยะเวลาโดยประมาณ
	Price       float64        `gorm:"not null" json:"price"`       // ราคาบริการ
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

