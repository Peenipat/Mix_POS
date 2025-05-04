
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
	service.CreatedAt = time.Now()
	service.UpdatedAt = time.Now()
	return s.DB.Create(service).Error
}

func (s *ServiceService) UpdateService(id uint, updated *barberBookingModels.Service) (*barberBookingModels.Service, error) {
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
	if err := s.DB.Delete(&barberBookingModels.Service{}, id).Error; err != nil {
		return err
	}
	return nil
}
