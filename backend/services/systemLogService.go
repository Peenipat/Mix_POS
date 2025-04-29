package services

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"myapp/models"
)

// LogFilter defines parameters for querying system logs
// - UserID: กรองตามผู้ใช้ (UUID)
// - Action: กรองตามชื่อเหตุการณ์ เช่น "LOGIN_SUCCESS"
// - Endpoint: กรองตาม API endpoint เช่น "/orders"
// - From, To: กรองช่วงเวลาของ created_at
// - Status: กรองตามสถานะ success/failure
// - BranchID: กรองตามสาขา (UUID)
// - Page, Limit: pagination (เริ่ม page=1)
type LogFilter struct {
	UserID   *uuid.UUID  
	Action   *string     
	Endpoint *string     
	From     *time.Time  
	To       *time.Time  
	Status   *string     
	BranchID *uuid.UUID  
	Page     int         
	Limit    int         
}

// SystemLogService defines available methods for system logging
// - Create: สร้าง log ใหม่
// - Query: ค้นหา log ตามเงื่อนไข พร้อม pagination
// - GetByID: ดึง log รายการเดียวตาม ID
//
type SystemLogService interface {
	Create(ctx context.Context, entry *models.SystemLog) error
	Query(ctx context.Context, filter LogFilter) ([]models.SystemLog, int64, error)
	GetByID(ctx context.Context, id uint) (models.SystemLog, error)
}

// systemLogService is the default implementation of SystemLogService
type systemLogService struct {
	db *gorm.DB
}

// NewSystemLogService creates a new SystemLogService
func NewSystemLogService(db *gorm.DB) SystemLogService {
	return &systemLogService{db: db}
}

// Create inserts a new log entry into the database
func (s *systemLogService) Create(ctx context.Context, entry *models.SystemLog) error {
	return s.db.WithContext(ctx).Create(entry).Error
}

// Query retrieves logs based on the provided filter, with pagination
func (s *systemLogService) Query(ctx context.Context, filter LogFilter) ([]models.SystemLog, int64, error) {
	tx := s.db.WithContext(ctx).Model(&models.SystemLog{})

	// Apply filters
	if filter.UserID != nil {
		tx = tx.Where("user_id = ?", *filter.UserID)
	}
	if filter.Action != nil {
		tx = tx.Where("action = ?", *filter.Action)
	}
	if filter.Endpoint != nil {
		tx = tx.Where("endpoint = ?", *filter.Endpoint)
	}
	if filter.Status != nil {
		tx = tx.Where("status = ?", *filter.Status)
	}
	if filter.BranchID != nil {
		tx = tx.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.From != nil {
		tx = tx.Where("created_at >= ?", *filter.From)
	}
	if filter.To != nil {
		tx = tx.Where("created_at <= ?", *filter.To)
	}

	// Count total
	var total int64
	if err := tx.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Apply pagination defaults
	if filter.Page < 1 {
		filter.Page = 1
	}
	if filter.Limit <= 0 {
		filter.Limit = 20
	}
	offset := (filter.Page - 1) * filter.Limit

	// Fetch logs ordered by latest first
	var logs []models.SystemLog
	err := tx.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetByID retrieves a single log entry by its primary key
func (s *systemLogService) GetByID(ctx context.Context, id uint) (models.SystemLog, error) {
	var entry models.SystemLog
	err := s.db.WithContext(ctx).First(&entry, "log_id = ?", id).Error
	return entry, err
}
