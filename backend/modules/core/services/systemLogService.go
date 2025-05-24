package coreServices

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"myapp/modules/core/models"
)

// LogFilter defines parameters for querying system logs
// - UserID: กรองตามผู้ใช้
// - Action: กรองตามชื่อเหตุการณ์ 
// - Endpoint: กรองตาม API endpoint 
// - From, To: กรองช่วงเวลาของ created_at
// - Status: กรองตามสถานะ success/failure
// - BranchID: กรองตามสาขา 
// - Page, Limit: pagination (เริ่ม page=1)
type LogFilter struct {
	UserID   *uuid.UUID  
	Action   *string     
	Endpoint *string     
	From     *time.Time  // จากช่วงเวลา
	To       *time.Time  //ถึงช่วงเวลา
	Status   *string     
	BranchID *uuid.UUID  
	Page     int         
	Limit    int         
}

// SystemLogService defines available methods for system logging
// - Create: สร้าง log ใหม่
// - Query: ค้นหา log ตามเงื่อนไข พร้อม pagination
// - GetByID: ดึง log รายการเดียวตาม ID
//สร้าง interface สำหรับแต่ละ function
type SystemLogService interface {
	Create(ctx context.Context, entry *coreModels.SystemLog) error
	Query(ctx context.Context, filter LogFilter) ([]coreModels.SystemLog, int64, error)
	GetByID(ctx context.Context, id uint) (coreModels.SystemLog, error)
}
type systemLogService struct {
	db *gorm.DB
}

// NewSystemLogService creates a new SystemLogService
func NewSystemLogService(db *gorm.DB) SystemLogService {
	return &systemLogService{db: db}
}

// สร้าง log และ save
func (s *systemLogService) Create(ctx context.Context, entry *coreModels.SystemLog) error {
	return s.db.WithContext(ctx).Create(entry).Error
}

// ดึง log ตามเงื่อนไข
func (s *systemLogService) Query(ctx context.Context, filter LogFilter) ([]coreModels.SystemLog, int64, error) {
	tx := s.db.WithContext(ctx).Model(&coreModels.SystemLog{})

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

	// ดึงข้อมูลและเรียงจากอันที่สร้างล่าสุดไปหาเก่าที่สุด
	var logs []coreModels.SystemLog
	err := tx.Order("created_at DESC").Offset(offset).Limit(filter.Limit).Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// ดึง log ตาม ID
func (s *systemLogService) GetByID(ctx context.Context, id uint) (coreModels.SystemLog, error) {
	var entry coreModels.SystemLog
	err := s.db.WithContext(ctx).First(&entry, "log_id = ?", id).Error
	return entry, err
}
