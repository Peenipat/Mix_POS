package models

import (
	"time"

	"gorm.io/datatypes"
)

// SystemLog represents an entry in the system_logs table
// It captures events and metadata for auditing, debugging, and analytics.
type SystemLog struct {
	LogID         uint           `gorm:"primaryKey;column:log_id"`
	CreatedAt     time.Time      `gorm:"column:created_at;autoCreateTime"`

	UserID        *uint          `gorm:"column:user_id"`          // FK to users.id
	UserRole      *string        `gorm:"type:varchar(20);column:user_role"`

	Action        string         `gorm:"type:varchar(50);column:action;not null"`
	Resource      string         `gorm:"type:varchar(50);column:resource;not null"`
	Status        string         `gorm:"type:varchar(20);column:status;not null"`

	IPAddress 	  *string 		 `gorm:"type:inet;column:ip_address"`

	HTTPMethod    string         `gorm:"type:varchar(10);column:http_method;default:'GET';not null"`
	Endpoint      string         `gorm:"type:varchar(255);column:endpoint;not null"`
	StatusCode    *int           `gorm:"column:status_code"`

	XForwardedFor *string        `gorm:"type:varchar(100);column:x_forwarded_for"`
	UserAgent     *string        `gorm:"type:text;column:user_agent"`
	Referer       *string        `gorm:"type:text;column:referer"`
	Origin        *string        `gorm:"type:text;column:origin"`

	ClientApp     *string        `gorm:"type:varchar(50);column:client_app"`
	BranchID      *uint          `gorm:"column:branch_id"`        // FK to branches.id

	Details       datatypes.JSON `gorm:"type:jsonb;column:details"`
	Metadata      datatypes.JSON `gorm:"type:jsonb;column:metadata"`
}

// TableName overrides the default table name
func (SystemLog) TableName() string {
	return "system_logs"
}
