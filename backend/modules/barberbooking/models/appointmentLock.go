package barberBookingModels
import (
	"time"
)

type AppointmentLock struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	TenantID   uint      `gorm:"not null;index" json:"tenant_id"`
	BranchID   uint      `gorm:"not null;index" json:"branch_id"`
	BarberID   uint      `gorm:"not null;index" json:"barber_id"`
	CustomerID uint      `gorm:"not null;index" json:"customer_id"`

	StartTime  time.Time `gorm:"not null;index" json:"start_time"`
	EndTime    time.Time `gorm:"not null;index" json:"end_time"`

	ExpiresAt  time.Time `gorm:"index" json:"expires_at"`

	IsActive   bool      `gorm:"default:true" json:"is_active"`

	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}
