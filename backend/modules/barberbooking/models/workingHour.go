package bookingModels
import (
	"time"
)
type WorkingHour struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BranchID  uint           `gorm:"not null;index" json:"branch_id"`
	Weekday   int            `gorm:"not null" json:"weekday"`      // 0=Sunday…6=Saturday
	StartTime time.Time      `gorm:"not null" json:"start_time"`
	EndTime   time.Time      `gorm:"not null" json:"end_time"`
  }
  