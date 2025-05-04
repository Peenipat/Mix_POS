package barberBookingModels
import (
	"time"
)
type BarberWorkloads struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	BarberID         uint      `gorm:"not null;index:idx_barber_date,unique" json:"barber_id"`
	Barber           Barber    `gorm:"foreignKey:BarberID" json:"barber,omitempty"`
	Date             time.Time `gorm:"not null;index:idx_barber_date,unique" json:"date"`
	TotalAppointments int      `gorm:"not null;default:0" json:"total_appointments"`
	TotalHours       int       `gorm:"not null;default:0" json:"total_hours"`
	CreatedAt        time.Time `gorm:"default:now()" json:"created_at"`
}
