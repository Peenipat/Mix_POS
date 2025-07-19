package barberBookingModels

import (
	"time"
	"gorm.io/gorm"
	helperFunc "myapp/modules/barberbooking"
)

type WorkingDayOverride struct {
	ID        uint      			`gorm:"primaryKey" json:"id"`
	BranchID  uint      			`gorm:"not null" json:"branch_id"`
	WorkDate  time.Time 			`gorm:"type:date;not null" json:"work_date"`
	StartTime helperFunc.TimeOnly 	`gorm:"type:time;not null" json:"start_time"`
	EndTime   helperFunc.TimeOnly 	`gorm:"type:time;not null" json:"end_time"`
	IsClosed  bool           		`gorm:"not null json:is_closed"` 
	CreatedAt time.Time 			`gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time 			`gorm:"autoUpdateTime" json:"updated_at"`
   	DeletedAt gorm.DeletedAt 		`gorm:"index" json:"deleted_at,omitempty"`
}
