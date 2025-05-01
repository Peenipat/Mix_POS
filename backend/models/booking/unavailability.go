package bookingModels
import (
	"time"
)
type Unavailability struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	BarberID  *uint          `gorm:"index" json:"barber_id,omitempty"`
	BranchID  *uint          `gorm:"index" json:"branch_id,omitempty"`
	Date      time.Time      `gorm:"not null;index" json:"date"`
	Reason    string         `gorm:"type:text" json:"reason"`
  }
  