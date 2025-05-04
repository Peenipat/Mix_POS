package coreModels

import (
	"time"

	"gorm.io/gorm"
)

// Branch represents a physical location of a Tenant
// Supports soft delete via DeletedAt field
type Branch struct {
    ID        uint           `gorm:"primaryKey" json:"id"`
    TenantID  uint           `gorm:"not null;index" json:"tenant_id"`
    Name      string         `gorm:"type:text;not null" json:"name"`
    Address   *string        `gorm:"type:text" json:"address,omitempty"`
    CreatedAt time.Time      `gorm:"autoCreateTime" json:"created_at"`
    UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

    // Relations (optional preload)
    Tenant    *Tenant        `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
    Users     []User         `gorm:"foreignKey:BranchID" json:"users,omitempty"`
}