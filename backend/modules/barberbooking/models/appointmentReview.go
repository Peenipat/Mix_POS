package barberBookingModels

import (
	"gorm.io/gorm"
	"time"
)

type AppointmentReview struct {
	ID            uint        `gorm:"primaryKey" json:"id"`
	AppointmentID uint        `gorm:"not null;uniqueIndex" json:"appointment_id"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`

	CustomerID *uint     `gorm:"index" json:"customer_id,omitempty"`
	Customer   *Customer `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`

	Rating    int            `gorm:"not null;check:rating >= 1 AND rating <= 5" json:"rating"`
	Comment   string         `gorm:"type:text" json:"comment,omitempty"`
	CreatedAt time.Time      `gorm:"default:now()" json:"created_at"`
	UpdatedAt time.Time      `gorm:"autoUpdateTime" json:"updated_at,omitempty"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
