package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"time"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	aws "myapp/cmd/worker"

	"gorm.io/gorm"
)

type BarberService struct {
	DB *gorm.DB
}

func NewBarberService(db *gorm.DB) barberBookingPort.IBarber {
	return &BarberService{DB: db}
}

// CreateBarber creates a new barber
func (s *BarberService) CreateBarber(ctx context.Context, barber *barberBookingModels.Barber) error {
	// Validation ID
	if barber.BranchID == 0 {
		return fmt.Errorf("branch_id is required")
	}
	if barber.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}

	var existing barberBookingModels.Barber
	err := s.DB.WithContext(ctx).
		Unscoped(). // (return DeleteAt != nil)
		Where("user_id = ?", barber.UserID).
		First(&existing).Error

	if err == nil && existing.DeletedAt.Valid {
		// ถ้ามีและถูก soft-delete → ลบทิ้งจริงก่อน (hard delete)
		if err := s.DB.WithContext(ctx).Unscoped().Delete(&existing).Error; err != nil {
			return fmt.Errorf("failed to purge existing deleted barber: %w", err)
		}
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return fmt.Errorf("failed to check existing barber: %w", err)
	}

	// สร้างใหม่
	barber.CreatedAt = time.Now()
	barber.UpdatedAt = time.Now()
	return s.DB.WithContext(ctx).Create(barber).Error
}

// GetBarberByID fetches a single barber by ID
func (s *BarberService) GetBarberByID(ctx context.Context, id uint) (*barberBookingPort.BarberDetailResponse, error) {
	var barber barberBookingModels.Barber

	if err := s.DB.WithContext(ctx).
		Preload("User").
		First(&barber, id).Error; err != nil {
		return nil, err
	}

	resp := &barberBookingPort.BarberDetailResponse{
		ID:          barber.ID,
		BranchID:    barber.BranchID,
		TenantID:    barber.TenantID,
		UserID:      barber.UserID,
		RoleUser:    barber.RoleUser,
		Description: barber.Description,
		User: struct {
			ID          uint   `json:"id"`
			Username    string `json:"username"`
			Email       string `json:"email"`
			PhoneNumber string `json:"phone_number"`
			BranchID    uint   `json:"branch_id"`
			ImgPath     string `json:"Img_path"`
			ImgName     string `json:"Img_name"`
		}{
			ID:          barber.User.ID,
			Username:    barber.User.Username,
			Email:       barber.User.Email,
			PhoneNumber: barber.User.PhoneNumber,
			BranchID:    *barber.User.BranchID,
			ImgPath:     barber.User.Img_path,
			ImgName:     barber.User.Img_name,
		},
	}

	return resp, nil
}

// ListBarbers optionally filters by branch_id
func (s *BarberService) ListBarbersByBranch(ctx context.Context, branchID *uint) ([]barberBookingPort.BarberWithUser, error) {
	// Make a slice of the port’s DTO type
	var rows []barberBookingPort.BarberWithUser

	q := s.DB.WithContext(ctx).
		Model(&barberBookingModels.Barber{}).
		Select(`
            barbers.id,
            barbers.branch_id,
            barbers.user_id,
			users.phone_number,
            users.username,
            users.email,
            users.img_path,
            users.img_name,
            barbers.description,
            barbers.role_user,
            barbers.created_at,
            barbers.updated_at
        `).
		Joins(`LEFT JOIN users ON users.id = barbers.user_id`)

	if branchID != nil {
		q = q.Where("barbers.branch_id = ?", *branchID)
	}

	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}

func (s *BarberService) UpdateBarber(
	ctx context.Context,
	barberID uint,
	payload *barberBookingPort.UpdateBarberRequest,
	file *multipart.FileHeader, // ✅ เพิ่มรับไฟล์
) (*barberBookingModels.Barber, error) {
	// 1. ดึง Barber ปัจจุบัน
	var barber barberBookingModels.Barber
	if err := s.DB.WithContext(ctx).
		Preload("User").
		First(&barber, barberID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("barber not found")
		}
		return nil, fmt.Errorf("failed to fetch barber: %w", err)
	}

	// 2. ถ้ามีการอัปโหลดไฟล์รูปใหม่ → อัปโหลดไป S3
	if file != nil {
		imgPath, imgName, err := aws.UploadToS3(file, "barbers")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to S3: %w", err)
		}
		barber.User.Img_path = imgPath
		barber.User.Img_name = imgName
	}

	// 3. อัปเดต Barber fields
	barber.BranchID = payload.BranchID
	barber.Description = payload.Description
	barber.RoleUser = payload.RoleUser
	barber.UpdatedAt = time.Now()

	// 4. อัปเดต User fields
	if payload.Username != "" {
		barber.User.Username = payload.Username
	}
	if payload.Email != "" {
		barber.User.Email = payload.Email
	}
	if payload.PhoneNumber != "" {
		barber.User.PhoneNumber = payload.PhoneNumber
	}
	// ไม่ต้องอัปเดต img_path/img_name จาก payload แล้ว เพราะถ้าไฟล์ใหม่มาเราจัดการให้แล้ว
	barber.User.UpdatedAt = time.Now()

	// 5. Save ทั้ง barber และ user
	if err := s.DB.WithContext(ctx).Save(&barber).Error; err != nil {
		return nil, fmt.Errorf("failed to save barber: %w", err)
	}
	if err := s.DB.WithContext(ctx).Save(&barber.User).Error; err != nil {
		return nil, fmt.Errorf("failed to save user: %w", err)
	}

	// 6. Reload barber อีกรอบ
	if err := s.DB.WithContext(ctx).
		Preload("User").
		Preload("User.Role").
		First(&barber, barber.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to reload updated barber: %w", err)
	}

	return &barber, nil
}

// DeleteBarber performs soft delete
func (s *BarberService) DeleteBarber(ctx context.Context, id uint) error {
	result := s.DB.WithContext(ctx).Delete(&barberBookingModels.Barber{}, id)
	if result.RowsAffected == 0 {
		return errors.New("barber not found")
	}
	return result.Error
}

func (s *BarberService) GetBarberByUser(ctx context.Context, userID uint) (*barberBookingModels.Barber, error) {
	var barber barberBookingModels.Barber
	err := s.DB.WithContext(ctx).
		Where("user_id = ?", userID).
		First(&barber).Error
	if err != nil {
		return nil, err
	}
	return &barber, nil
}

func (s *BarberService) ListBarbersByTenant(ctx context.Context, tenantID uint) ([]barberBookingModels.Barber, error) {
	var barbers []barberBookingModels.Barber

	err := s.DB.WithContext(ctx).
		Joins("JOIN branches ON branches.id = barbers.branch_id").
		Where("branches.tenant_id = ?", tenantID).
		Where("barbers.deleted_at IS NULL").
		Find(&barbers).Error

	return barbers, err
}

func (s *BarberService) ListUserNotBarber(ctx context.Context, branchID *uint) ([]barberBookingPort.UserNotBarber, error) {
	// 1. ถ้า branchID เป็น nil ให้รีเทิร์น error ทันที
	if branchID == nil {
		return nil, errors.New("branchID is required")
	}

	rows := []barberBookingPort.UserNotBarber{}

	q := s.DB.WithContext(ctx).
		Model(&coreModels.User{}).
		Select(`
        users.id           AS user_id,
        users.username     AS username,
        users.email        AS email,
        users.created_at   AS created_at,
        users.updated_at   AS updated_at
    `).
		Joins(`
        LEFT JOIN barbers 
          ON barbers.user_id = users.id `).
		Where("barbers.id IS NULL AND users.branch_id = ?", *branchID)

	if err := q.Scan(&rows).Error; err != nil {
		return nil, err
	}
	return rows, nil
}
