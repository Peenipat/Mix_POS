package barberBookingModels

import (
	"time"

	"gorm.io/gorm"
)

// Service แทนบริการต่างๆ (เช่น ตัดผม สระผม ไดร์)
type Service struct {
	ID          	uint           `gorm:"primaryKey;autoIncrement" json:"id"` 
	TenantID    	uint           `gorm:"not null;index" json:"tenant_id"`
	BranchID 		uint 		   `gorm:"not null;index" json:"branch_id"`

	Name        	string         `gorm:"type:varchar(100);not null" json:"name"`
	Description     string      	`gorm:"type:varchar(100);not null" json:"description"`
	Duration    	int            `gorm:"not null" json:"duration"`   
	Price       	float64        `gorm:"not null" json:"price"`  
	Img_path  		string 			`gorm:"column:img_path" json:"Img_path,omitempty"`
    Img_name 		string 			`gorm:"column:img_name" json:"Img_name,omitempty"`

	CreatedAt   	time.Time      `json:"created_at"`
	UpdatedAt   	time.Time      `json:"updated_at"`
	DeletedAt   	gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

