package coreModels
import (
	"time"
  )
type User struct {
	ID        	  uint      	`json:"id" example:"1"`
  	CreatedAt 	  time.Time 	`json:"created_at"`
  	UpdatedAt 	  time.Time 	`json:"updated_at"`
  	DeletedAt 	  *time.Time 	`json:"deleted_at,omitempty"`
	Username      string
	Email         string     	`gorm:"uniqueIndex"`
	Password  	  string
	Role          RoleName       `gorm:"type:VARCHAR(20);not null"`     // SUPER_ADMIN, BRANCH_ADMIN
	BranchID      *uint          `gorm:"index" json:"branch_id,omitempty"`
	Branch        *Branch        `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	TenantUsers  []TenantUser   `gorm:"foreignKey:UserID" json:"tenant_users,omitempty"` // join table entries
}





