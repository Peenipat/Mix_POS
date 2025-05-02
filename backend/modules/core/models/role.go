// modules/core/models/role.go
package coreModels

import (
    "time"

    "gorm.io/gorm"
)

// RoleName ประกาศค่าคงที่สำหรับชื่อ role
type RoleName string

const (
    RoleNameSaaSSuperAdmin   RoleName = "SAAS_SUPER_ADMIN"
    RoleNameTenantAdmin      RoleName = "TENANT_ADMIN"
    RoleNameBranchAdmin      RoleName = "BRANCH_ADMIN"
    RoleNameAssistantManager RoleName = "ASSISTANT_MANAGER"
    RoleNameStaff            RoleName = "STAFF"
    RoleNameUser             RoleName = "USER"
)

// Role ผูกกับ Module ผ่าน ModuleID
type Role struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    ModuleID    uint           `gorm:"not null;index" json:"module_id"`               // FK → modules(id)
    Module      Module         `gorm:"foreignKey:ModuleID" json:"module,omitempty"`   // preload โมดูลได้
    Name        RoleName       `gorm:"size:50;not null" json:"name"`                  // ชื่อเช่น SUPER_ADMIN
    Description string         `gorm:"type:text" json:"description"`                  // คำอธิบายบทบาท
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
