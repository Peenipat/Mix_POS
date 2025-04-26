package models

import "gorm.io/gorm"

type AuditLog struct {
	gorm.Model
	Action   string
	ID   uint
	OldValue string
	NewValue string
}
