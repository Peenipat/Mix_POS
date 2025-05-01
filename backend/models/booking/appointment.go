

package bookingModels
import (
	"time"
	"myapp/models/core"
	"gorm.io/gorm"
)

// AppointmentStatus แทนสถานะการจองคิว

type AppointmentStatus string

const (
	StatusPending   AppointmentStatus = "PENDING"
	StatusConfirmed AppointmentStatus = "CONFIRMED"
	StatusCancelled AppointmentStatus = "CANCELLED"
	StatusComplete  AppointmentStatus = "COMPLETED"
)

// Appointment แทนข้อมูลการจองคิว
// เชื่อม Branch, Service, Barber (optional), User
// ระบุเวลาเริ่มและเวลาสิ้นสุดของคิว
type Appointment struct {
	ID           uint                `gorm:"primaryKey" json:"id"`
	BranchID     uint                `gorm:"not null;index" json:"branch_id"`
	Branch       coreModels.Branch       `gorm:"foreignKey:BranchID" json:"branch,omitempty"`
	ServiceID    uint                `gorm:"not null;index" json:"service_id"`
	Service      Service             `gorm:"foreignKey:ServiceID" json:"service,omitempty"`
	BarberID     *uint               `gorm:"index" json:"barber_id,omitempty"`
	Barber       *Barber             `gorm:"foreignKey:BarberID" json:"barber,omitempty"`
	CustomerID   uint                `gorm:"not null;index" json:"customer_id"`
	Customer     coreModels.User         `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	StartTime    time.Time           `gorm:"not null;index" json:"start_time"`
	EndTime      time.Time           `gorm:"not null" json:"end_time"`
	Status       AppointmentStatus   `gorm:"type:varchar(20);not null;default:'PENDING'" json:"status"`
	Notes        string              `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt    time.Time           `json:"created_at"`
	UpdatedAt    time.Time           `json:"updated_at"`
	DeletedAt    gorm.DeletedAt      `gorm:"index" json:"deleted_at,omitempty"`
}
