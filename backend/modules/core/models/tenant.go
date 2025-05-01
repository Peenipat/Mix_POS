package coreModels
import (
	"time"

	"gorm.io/gorm"
)

// Tenant represents a SaaS tenant (a subscribing business)
type Tenant struct {
	ID        uint           `gorm:"primaryKey" json:"tenant_id"`
	Name      string         `gorm:"type:text;not null" json:"name"`
	Domain    string         `gorm:"type:text;uniqueIndex;not null" json:"domain"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}