package coreModels

import (
	"time"

	"gorm.io/gorm"
)

// Branch represents a physical location of a Tenant
// Supports soft delete via DeletedAt field
type Branch struct {
    ID        uint           `gorm:"primaryKey" json:"id"`             // รหัสสาขา
    TenantID  uint           `gorm:"not null;index" json:"tenant_id"`  // FK ไปยัง tenants.id
    Name      string         `gorm:"type:text;not null" json:"name"`   // ชื่อสาขา
    Address   *string        `gorm:"type:text" json:"address,omitempty"` // ที่อยู่สาขา
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

    // ความสัมพันธ์
    Tenant    Tenant         `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
    Users     []User         `gorm:"foreignKey:BranchID" json:"users,omitempty"`
    // ถ้ามีข้อมูล WorkingHour, Unavailability ก็ใส่ relation เพิ่มได้
}