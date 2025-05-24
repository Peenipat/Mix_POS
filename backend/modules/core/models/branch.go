package coreModels

import (
	"time"

	"gorm.io/gorm"
)

// Branch represents a physical location of a Tenant
// Supports soft delete via DeletedAt field
type Branch struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	TenantID  uint           `gorm:"not null;uniqueIndex:idx_tenant_name" json:"tenant_id"`
	Name      string         `gorm:"type:text;not null;uniqueIndex:idx_tenant_name" json:"name"` // composite unique with TenantID
	Address   *string        `gorm:"type:text" json:"address,omitempty"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	Tenant Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Users  []User `gorm:"foreignKey:BranchID" json:"users,omitempty"`
}
