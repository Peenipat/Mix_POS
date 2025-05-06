package barberBookingService

import (
	"context"
	"errors"

	barberBookingModels "myapp/modules/barberbooking/models"
	"gorm.io/gorm"
)



type customerService struct {
	db *gorm.DB
}

func NewCustomerService(db *gorm.DB) ICustomerService {
	return &customerService{db: db}
}

//  ดึงรายชื่อลูกค้าทั้งหมดของ tenant
func (s *customerService) GetAllCustomers(ctx context.Context, tenantID uint) ([]barberBookingModels.Customer, error) {
	var customers []barberBookingModels.Customer
	if err := s.db.WithContext(ctx).Where("tenant_id = ?", tenantID).Find(&customers).Error; err != nil {
		return nil, err
	}
	return customers, nil
}

// ดึงข้อมูลลูกค้ารายเดียว
func (s *customerService) GetCustomerByID(ctx context.Context, tenantID, customerID uint) (*barberBookingModels.Customer, error) {
	var customer barberBookingModels.Customer
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, customerID).
		First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &customer, nil
}

// เพิ่มลูกค้าใหม่
func (s *customerService) CreateCustomer(ctx context.Context, customer *barberBookingModels.Customer) error {
	// ✅ ตรวจสอบ input เบื้องต้น
	if customer == nil {
		return errors.New("customer data is required")
	}
	if customer.TenantID == 0 {
		return errors.New("tenant_id is required")
	}
	if customer.Email == "" {
		return errors.New("email is required")
	}
	if customer.Name == "" {
		return errors.New("name is required")
	}

	// ❌ ห้าม email ซ้ำใน tenant เดียวกัน
	var existing barberBookingModels.Customer
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", customer.TenantID, customer.Email).
		First(&existing).Error; err == nil {
		return errors.New("customer with this email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	return s.db.WithContext(ctx).Create(customer).Error
}


// แก้ไขข้อมูลลูกค้า
func (s *customerService) UpdateCustomer(ctx context.Context, tenantID, customerID uint, updateData map[string]interface{}) error {
	if len(updateData) == 0 {
		// ไม่ error แต่ไม่อัปเดต
		return nil
	}

	// เช็คว่ามี record อยู่จริงก่อน
	var existing barberBookingModels.Customer
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, customerID).
		First(&existing).Error; err != nil {
		return err
	}

	// ถ้าเจอแล้ว → อัปเดต
	result := s.db.WithContext(ctx).
		Model(&existing).
		Updates(updateData)

	return result.Error
}


// ลบลูกค้า
func (s *customerService) DeleteCustomer(ctx context.Context, tenantID, customerID uint) error {
	result := s.db.WithContext(ctx).
		Where("tenant_id = ? AND id = ?", tenantID, customerID).
		Delete(&barberBookingModels.Customer{})

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *customerService) FindCustomerByEmail(ctx context.Context, tenantID uint, email string) (*barberBookingModels.Customer, error) {
	var customer barberBookingModels.Customer
	if err := s.db.WithContext(ctx).
		Where("tenant_id = ? AND email = ?", tenantID, email).
		First(&customer).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // ไม่เจอ = ไม่ error
		}
		return nil, err
	}
	return &customer, nil
}

