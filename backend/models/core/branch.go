package coreModels

import (
	"time"

	"gorm.io/gorm"
)

// Branch represents a physical location of a Tenant
// Supports soft delete via DeletedAt field
type Branch struct {
	ID        uint           `gorm:"primaryKey" json:"branch_id"`
	TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
	Tenant    Tenant         `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Name      string         `gorm:"type:text;not null" json:"name"`
	Address   string         `gorm:"type:text" json:"address"`
	Timezone  string         `gorm:"type:text" json:"timezone"`
	Contact   string         `gorm:"type:text" json:"contact"`
	IsActive  bool           `gorm:"default:true" json:"is_active"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}