// modules/core/models/module.go
package coreModels

import (
    "time"

    "gorm.io/gorm"
)

// Module เก็บชื่อโมดูล/ฟีเจอร์หลักของระบบ
type Modules struct {
    ID          uint           `gorm:"primaryKey" json:"id"`
    Key         string         `gorm:"size:50;not null;uniqueIndex" json:"key"`       // เช่น "CORE","BOOKING","POS_RESTAURANT"
    Description string         `gorm:"type:text;not null" json:"description"`         // คำอธิบายโมดูล
    CreatedAt   time.Time      `json:"created_at"`
    UpdatedAt   time.Time      `json:"updated_at"`
    DeletedAt   gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
