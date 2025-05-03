package barberBookingModels
import (
	"time"
)
type AppointmentStatusLog struct {
	ID                   uint      `gorm:"primaryKey" json:"id"`
	AppointmentID        uint      `gorm:"not null;index" json:"appointment_id"`
	Appointment          Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`

	OldStatus            string    `gorm:"type:varchar(20)" json:"old_status,omitempty"`
	NewStatus            string    `gorm:"type:varchar(20)" json:"new_status,omitempty"`
	
	ChangedByUserID      *uint     `gorm:"index" json:"changed_by_user_id,omitempty"` // ↔ loose link to core.users
	ChangedByCustomerID  *uint     `gorm:"index" json:"changed_by_customer_id,omitempty"`
	ChangedByCustomer    *Customer `gorm:"foreignKey:ChangedByCustomerID" json:"changed_by_customer,omitempty"`

	ChangedAt            time.Time `gorm:"default:now()" json:"changed_at"`
	Notes                string    `gorm:"type:text" json:"notes,omitempty"`
}
