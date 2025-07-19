// modules/core/models/role.go
package coreModels
import (
    "time"

    "gorm.io/gorm"
)

type RoleName string

const (
	RoleNameSaaSSuperAdmin   RoleName = "SAAS_SUPER_ADMIN" // admin ฝั่ง SaaS ดูแลระบบทั้งหมด
	RoleNameTenant      	 RoleName = "TENANT"
	RoleNameTenantAdmin      RoleName = "TENANT_ADMIN"
	RoleNameBranchAdmin      RoleName = "BRANCH_ADMIN"
	RoleNameAssistantManager RoleName = "ASSISTANT_MANAGER"
	RoleNameStaff            RoleName = "STAFF"
	RoleNameUser             RoleName = "USER"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`

	TenantID    *uint          `gorm:"index;uniqueIndex:uq_roles_scope,priority:1" json:"tenant_id,omitempty"`
	Tenant      *Tenant        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tenant,omitempty"`

	ModuleID    *uint          `gorm:"index;uniqueIndex:uq_roles_scope,priority:2" json:"module_id,omitempty"`
	Module      *Module        `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"module,omitempty"`

	Name        string         `gorm:"type:varchar(50);not null;uniqueIndex:uq_roles_scope,priority:3" json:"name"`
	Description string         `gorm:"type:text" json:"description"`

	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

