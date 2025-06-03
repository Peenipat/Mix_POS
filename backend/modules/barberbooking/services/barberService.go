package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"time"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"

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


	// ลบ record เดิมที่ถูก soft-delete ไปแล้ว (ถ้ามี user_id เดิม)
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
func (s *BarberService) GetBarberByID(ctx context.Context, id uint) (*barberBookingModels.Barber, error) {
	var barber barberBookingModels.Barber
	if err := s.DB.WithContext(ctx).First(&barber, id).Error; err != nil {
		return nil, err
	}
	return &barber, nil
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
			barbers.phone_number,
            users.username,
            users.email,
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

// UpdateBarber updates barber info
func (s *BarberService) UpdateBarber(
    ctx context.Context,
    barberID uint,
    updated *barberBookingModels.Barber,
    updatedUsername string,
    updatedEmail string,
) (*barberBookingModels.Barber, error) {
    // 1. ดึงข้อมูล barber ปัจจุบัน
    var barber barberBookingModels.Barber
    if err := s.DB.WithContext(ctx).First(&barber, barberID).Error; err != nil {
        return nil, err
    }

    // 2. ถ้าอยากอัปเดตเบอร์โทรศัพท์ ให้เซ็ตเข้าไป
    barber.PhoneNumber = updated.PhoneNumber

    // 3. (ถ้าต้องการเปลี่ยนสาขา หรือเปลี่ยนผู้ใช้ที่ผูกอยู่ ก็ใส่ตรงนี้)
    barber.BranchID = updated.BranchID
    barber.UserID = updated.UserID

    // 4. อัปเดตวันที่แก้ไข
    barber.UpdatedAt = time.Now()

    // 5. บันทึกลง table barbers
    if err := s.DB.WithContext(ctx).Save(&barber).Error; err != nil {
        return nil, err
    }

    // 6. ถ้าต้องการอัปเดตข้อมูลในตาราง users (username / email) ก็ทำแยกอีกที
    if updatedUsername != "" || updatedEmail != "" {
        // ดึง User record เดิมตาม barber.UserID
        var user coreModels.User
        if err := s.DB.WithContext(ctx).First(&user, barber.UserID).Error; err != nil {
            return nil, err
        }

        // เซ็ตค่าถ้ามีส่งมา
        if updatedUsername != "" {
            user.Username = updatedUsername
        }
        if updatedEmail != "" {
            user.Email = updatedEmail
        }
        user.UpdatedAt = time.Now()

        if err := s.DB.WithContext(ctx).Save(&user).Error; err != nil {
            return nil, err
        }
    }

    // 7. รีเทิร์น barber object กลับไป (ถ้าต้องการ preload ข้อมูล user/branch เพิ่ม ให้ใช้ Preload ก่อน Save)
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






