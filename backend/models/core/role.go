package coreModels
import (
	"time"

	"gorm.io/gorm"
)
type RoleName string
const (
    RoleNameSaaSSuperAdmin   RoleName = "SAAS_SUPER_ADMIN"
    RoleNameTenantAdmin      RoleName = "TENANT_ADMIN"
    RoleNameBranchAdmin      RoleName = "BRANCH_ADMIN"
    RoleNameAssistantManager RoleName = "ASSISTANT_MANAGER"
    RoleNameStaff            RoleName = "STAFF"
    RoleNameUser        	 RoleName = "USER"
)

// Role defines a system role for RBAC
// Supports soft delete via DeletedAt field
type Role struct {
	ID          uint           `gorm:"primaryKey" json:"role_id"`
	Name        string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"name"`
	Description string         `gorm:"type:text" json:"description"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}