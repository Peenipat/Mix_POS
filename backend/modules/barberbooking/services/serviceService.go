
package barberBookingService

import (
	"errors"
	barberBookingModels "myapp/modules/barberbooking/models"
	"gorm.io/gorm"
	"time"
)

type ServiceService struct {
	DB *gorm.DB
}

func NewServiceService(db *gorm.DB) *ServiceService {
	return &ServiceService{DB: db}
}

func validateServiceInput(svc *barberBookingModels.Service) error {
	if len(svc.Name) < 2 || len(svc.Name) > 100 {
		return errors.New("name must be 2 - 100 characters")
	}
	if svc.Duration <= 0 || svc.Duration > 240 {
		return errors.New("duration must be 1 - 240 minutes")
	}
	if svc.Price < 0 || svc.Price > 100000 {
		return errors.New("price must be between 0 and 100,000")
	}
	return nil
}

func (s *ServiceService) GetAllServices() ([]barberBookingModels.Service, error) {
	var services []barberBookingModels.Service
	if err := s.DB.Find(&services).Error; err != nil {
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

func (s *ServiceService) CreateService(service *barberBookingModels.Service) error {
	if service.Name == "" || service.Duration <= 0 || service.Price < 0 {
		return errors.New("invalid service data")
	}

	if err := validateServiceInput(service) ; err != nil{
		return err
	}
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	
	return s.DB.Create(service).Error
}

func (s *ServiceService) UpdateService(id uint, updated *barberBookingModels.Service) (*barberBookingModels.Service, error) {
	if err := validateServiceInput(updated); err != nil {
		return nil, err
	}

	var existing barberBookingModels.Service
	if err := s.DB.First(&existing, id).Error; err != nil {
		return nil, err
	}

	existing.Name = updated.Name
	existing.Price = updated.Price
	existing.Duration = updated.Duration
	existing.UpdatedAt = time.Now()

	if err := s.DB.Save(&existing).Error; err != nil {
		return nil, err
	}
	return &existing, nil
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


