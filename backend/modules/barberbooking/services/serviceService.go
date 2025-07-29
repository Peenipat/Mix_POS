package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	aws "myapp/cmd/worker"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"time"

	"gorm.io/gorm"
)

type ServiceService struct {
	DB *gorm.DB
}

func NewServiceService(db *gorm.DB) barberBookingPort.IServiceService {
	return &ServiceService{DB: db}
}

// func validateServiceInput(svc *barberBookingModels.Service) error {
// 	if len(svc.Name) < 2 || len(svc.Name) > 100 {
// 		return errors.New("name must be 2 - 100 characters")
// 	}
// 	if svc.Duration <= 0 || svc.Duration > 240 {
// 		return errors.New("duration must be 1 - 240 minutes")
// 	}
// 	if svc.Price < 0 || svc.Price > 100000 {
// 		return errors.New("price must be between 0 and 100,000")
// 	}
// 	return nil
// }

func (s *ServiceService) GetAllServices(tenantID uint, branchID uint) ([]barberBookingModels.Service, error) {
    var services []barberBookingModels.Service

    if err := s.DB.
        Where("tenant_id = ? AND branch_id = ?", tenantID, branchID).Order("id asc").
        Find(&services).Error; err != nil {
        return nil, err
    }
    return services, nil
}

func (s *ServiceService) GetServiceByID(id uint) (*barberBookingModels.Service, error) {
	var service barberBookingModels.Service
	if err := s.DB.First(&service, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &service, nil
}

func (s *ServiceService) CreateService(
	ctx context.Context,
	tenantID uint,
	branchID uint,
	payload *barberBookingPort.CreateServiceRequest, 
	file *multipart.FileHeader, 
) (*barberBookingModels.Service, error) {
	// 1. Validate ข้อมูลเบื้องต้น
	if payload.Name == "" || payload.Duration <= 0 || payload.Price < 0 {
		return nil, fmt.Errorf("invalid service data")
	}
	if tenantID == 0 || branchID == 0 {
		return nil, fmt.Errorf("tenant_id and branch_id cannot be zero")
	}

	service := &barberBookingModels.Service{
		TenantID:    tenantID,
		BranchID:    branchID,
		Name:        payload.Name,
		Description: payload.Description,
		Duration:    payload.Duration,
		Price:       payload.Price,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// 3. ถ้ามีไฟล์ → อัปโหลดขึ้น S3
	if file != nil {
		imgPath, imgName, err := aws.UploadToS3(file, "services") // ✅ กำหนด folder ใน S3
		if err != nil {
			return nil, fmt.Errorf("failed to upload image to S3: %w", err)
		}
		service.Img_path = imgPath
		service.Img_name = imgName
	}

	// 4. Save service
	if err := s.DB.WithContext(ctx).Create(service).Error; err != nil {
		return nil, fmt.Errorf("failed to create service: %w", err)
	}

	return service, nil
}

func (s *ServiceService) UpdateService(
	ctx context.Context,
	serviceID uint,
	payload *barberBookingPort.UpdateServiceRequest,
	file *multipart.FileHeader, // ✅ รูปภาพใหม่ (optional)
) (*barberBookingModels.Service, error) {
	// 1. ตรวจสอบค่าเบื้องต้น
	if payload.Name == "" || payload.Duration <= 0 || payload.Price < 0 {
		return nil, fmt.Errorf("invalid service input")
	}

	// 2. ดึงข้อมูลเดิม
	var service barberBookingModels.Service
	if err := s.DB.WithContext(ctx).First(&service, serviceID).Error; err != nil {
		return nil, fmt.Errorf("service not found")
	}

	// 3. ถ้ามีไฟล์ใหม่ → อัปโหลดขึ้น S3
	if file != nil {
		imgPath, imgName, err := aws.UploadToS3(file, "services")
		if err != nil {
			return nil, fmt.Errorf("failed to upload image: %w", err)
		}
		service.Img_path = imgPath
		service.Img_name = imgName
	}

	// 4. อัปเดตฟิลด์ที่เปลี่ยน
	service.Name = payload.Name
	service.Description = payload.Description
	service.Duration = payload.Duration
	service.Price = payload.Price
	service.UpdatedAt = time.Now()

	// 5. Save
	if err := s.DB.WithContext(ctx).Save(&service).Error; err != nil {
		return nil, fmt.Errorf("failed to update service: %w", err)
	}

	return &service, nil
}


func (s *ServiceService) DeleteService(id uint) error {
	var service barberBookingModels.Service
	if err := s.DB.First(&service, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("service not found")
		}
		return err
	}

	return s.DB.Delete(&service).Error
}


