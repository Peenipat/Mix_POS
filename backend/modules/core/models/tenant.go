package coreModels
import (
	"time"
)

// Tenant represents a SaaS tenant (a subscribing business)
type Tenant struct {
    ID        uint      `gorm:"primaryKey" json:"id"`
    Name      string    `gorm:"type:text;not null" json:"name"`
    Domain    string    `gorm:"type:text;uniqueIndex;not null" json:"domain"`
    IsActive  bool      `gorm:"default:true;not null" json:"is_active"`
    CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt *time.Time `gorm:"index" json:"deleted_at,omitempty"`
}