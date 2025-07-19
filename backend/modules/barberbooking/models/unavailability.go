package barberBookingModels
import (
	"time"
	"gorm.io/gorm"
)


  type Unavailability struct {
	ID        uint           `gorm:"primaryKey" json:"id"`

	BarberID  *uint          `gorm:"index" json:"barber_id,omitempty"`  // → booking.barbers.id
	BranchID  *uint          `gorm:"index" json:"branch_id,omitempty"`  // → core.branches.id
	Date      time.Time      `gorm:"not null;index" json:"date"`
	Reason    string         `gorm:"type:text" json:"reason,omitempty"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// ❗ DO NOT preload Barber / Branch (avoid tight coupling)
}
