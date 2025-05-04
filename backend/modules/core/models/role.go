package coreModels

import (
	"time"

	"gorm.io/gorm"
)

type RoleName string

const (
	RoleNameSaaSSuperAdmin   RoleName = "SAAS_SUPER_ADMIN" // admin ฝั่ง SaaS ดูแลระบบทั้งหมด
	RoleNameTenantAdmin      RoleName = "TENANT_ADMIN"
	RoleNameBranchAdmin      RoleName = "BRANCH_ADMIN"
	RoleNameAssistantManager RoleName = "ASSISTANT_MANAGER"
	RoleNameStaff            RoleName = "STAFF"
	RoleNameUser             RoleName = "USER"
)

type Role struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	TenantID    *uint          `gorm:"index;uniqueIndex:uq_roles_scope,priority:1" json:"tenant_id,omitempty"`     // รองรับ Global Role ถ้า null
	ModuleName  *string        `gorm:"type:varchar(50);index;uniqueIndex:uq_roles_scope,priority:2" json:"module_name,omitempty"`
	Name        string         `gorm:"type:varchar(50);not null;uniqueIndex:uq_roles_scope,priority:3" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
