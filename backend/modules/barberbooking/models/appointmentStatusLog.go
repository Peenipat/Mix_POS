package barberBookingModels
import (
	"time"
)
type AppointmentStatusLog struct {
	ID                  uint      `gorm:"primaryKey"`
	AppointmentID       uint      `gorm:"not null;index"`              // ใช้ struct ได้ถ้าต้องการ preload

	OldStatus           string    `gorm:"type:varchar(20)"`
	NewStatus           string    `gorm:"type:varchar(20)"`

	ChangedByUserID     *uint     `gorm:"index"`                       // Loose FK → ไม่ preload
	ChangedByCustomerID *uint     `gorm:"index"`                       // Loose FK → ไม่ preload
	ChangedAt           time.Time `gorm:"autoCreateTime"`

	Notes               string    `gorm:"type:text"`
}
