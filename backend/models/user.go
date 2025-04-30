package models
import (
	"time"
  )
type StoreRole string
type User struct {
	ID        	  uint      	`json:"id" example:"1"`
  	CreatedAt 	  time.Time 	`json:"created_at"`
  	UpdatedAt 	  time.Time 	`json:"updated_at"`
  	DeletedAt 	  *time.Time 	`json:"deleted_at,omitempty"`
	Username      string
	Email         string     	`gorm:"uniqueIndex"`
	Password  	  string
	Role          RoleName       `gorm:"type:VARCHAR(20);not null"`     // SUPER_ADMIN, BRANCH_ADMIN
// 	StoreRole     StoreRole  `gorm:"type:VARCHAR(30);default:'EMPLOYEE'"` // OWNER, MANAGER, etc.
// 	StoreID       *uint      // null = super admin
// 	Store         *Store
}
const (
	StoreOwner        StoreRole = "OWNER"
	StoreManager      StoreRole = "MANAGER" 
	StoreViceManager  StoreRole = "VICE_MANAGER"
	StoreAssistant    StoreRole = "ASSISTANT_MANAGER"
	StoreEmployee     StoreRole = "EMPLOYEE"
	StorePartTime     StoreRole = "PART_TIME"
	StoreIntern       StoreRole = "INTERN"
)




