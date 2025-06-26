package barberBookingModels
import (
	"time"
	"gorm.io/gorm"
)
type WorkingHour struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BranchID  uint           `gorm:"not null;uniqueIndex:idx_wh_branch_weekday"`
	TenantID  uint           `gorm:"not null;uniqueIndex:idx_wh_branch_weekday"`
    Weekday   int            `gorm:"not null;uniqueIndex:idx_wh_branch_weekday"`
	StartTime time.Time      `gorm:"not null" json:"start_time"`
	EndTime   time.Time      `gorm:"not null" json:"end_time"`
	IsClosed  bool           `gorm:"not null json:is_closed"` 
	CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
  }
  