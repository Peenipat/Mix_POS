package Core_authDto

import "time"

type CreateLogRequest struct {
    UserID      *string   `json:"user_id,omitempty"`
    Action      string    `json:"action" binding:"required"`
    Resource    *string   `json:"resource,omitempty"`
    Status      string    `json:"status" binding:"required"`
    HTTPMethod  string    `json:"http_method" binding:"required"`
    Endpoint    string    `json:"endpoint" binding:"required"`
    Details     any       `json:"details,omitempty"`    // auto-marshaled to JSONB
    BranchID    *string   `json:"branch_id,omitempty"`
}

type CreateLogResponse struct {
    LogID     uint      `json:"log_id"`
    CreatedAt time.Time `json:"created_at"`
}
