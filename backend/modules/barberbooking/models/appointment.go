

package barberBookingModels
import (
	"time"
	"gorm.io/gorm"
	coreModels "myapp/modules/core/models"
)

// AppointmentStatus แทนสถานะการจองคิว

type AppointmentStatus string

const (
	StatusPending     AppointmentStatus = "PENDING" //รอดำเนินการ
	StatusConfirmed   AppointmentStatus = "CONFIRMED" //รับ
	StatusCancelled   AppointmentStatus = "CANCELLED" //ยกเลิก
	StatusComplete    AppointmentStatus = "COMPLETED" //จบงาน
	StatusNoShow      AppointmentStatus = "NO_SHOW"
	StatusRescheduled AppointmentStatus = "RESCHEDULED" //เปลี่ยนเวลา
)

type Appointment struct {
	ID         uint              `gorm:"primaryKey" json:"id"`
	BranchID   uint              `gorm:"index" json:"branch_id"`      // ไม่ preload branch
	
	ServiceID  uint              `gorm:"not null" json:"service_id"`
	Service    Service           `gorm:"foreignKey:ServiceID" json:"service,omitempty"`

	BarberID   uint             `gorm:"index" json:"barber_id,omitempty"`
	Barber     coreModels.User      `gorm:"foreignKey:BarberID;references:ID" json:"barber,omitempty"`

	CustomerID uint              `gorm:"not null;index" json:"customer_id"`
	Customer   Customer          `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`

	UserID     *uint             `gorm:"index" json:"user_id,omitempty"` // อ้าง user ที่สร้างคิว (ไม่มี FK)
	TenantID   uint 			 `gorm:"not null;index" json:"tenant_id"`
	StartTime  time.Time         `gorm:"not null" json:"start_time"`
	EndTime    time.Time         `gorm:"not null" json:"end_time"`
	Status     AppointmentStatus `gorm:"type:varchar(20);default:'PENDING'" json:"status"`
	Notes      string            `gorm:"type:text" json:"notes,omitempty"`
	CreatedAt  time.Time         `json:"created_at"`
	UpdatedAt  time.Time         `json:"updated_at"`
	DeletedAt  gorm.DeletedAt    `gorm:"index" json:"deleted_at,omitempty"`
}
