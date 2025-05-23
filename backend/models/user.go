package models
import (
	"gorm.io/gorm"
  )
type Role string
type StoreRole string

const (
	RoleSuperAdmin  Role = "SUPER_ADMIN" // admin กลางที่จะค่อยดูแลระบบทั้งหมด
	RoleBranchAdmin Role = "BRANCH_ADMIN" // admin แต่ละสาขา แต่ละร้านค้า
)

const (
	StoreOwner        StoreRole = "OWNER"
	StoreManager      StoreRole = "MANAGER" 
	StoreViceManager  StoreRole = "VICE_MANAGER"
	StoreAssistant    StoreRole = "ASSISTANT_MANAGER"
	StoreEmployee     StoreRole = "EMPLOYEE"
	StorePartTime     StoreRole = "PART_TIME"
	StoreIntern       StoreRole = "INTERN"
)


type User struct {
	gorm.Model
	Username      string
	Email         string     `gorm:"uniqueIndex"`
	Password  	  string
	Role          Role       `gorm:"type:VARCHAR(20);not null"`     // SUPER_ADMIN, BRANCH_ADMIN
	StoreRole     StoreRole  `gorm:"type:VARCHAR(30);default:'EMPLOYEE'"` // OWNER, MANAGER, etc.
	StoreID       *uint      // null = super admin
	Store         *Store
}
