package barberbookingtestService

import (
	"testing"
	"time"

	barberBookingServices "myapp/modules/barberbooking/services"
	barberBookingModels "myapp/modules/barberbooking/models"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestDB(t *testing.T) (*gorm.DB, *barberBookingServices.ServiceService) {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	assert.NoError(t, err)

	err = db.AutoMigrate(&barberBookingModels.Service{})
	assert.NoError(t, err)

	serviceService := barberBookingServices.NewServiceService(db)
	return db, serviceService
}

// test สร้าง service พร้อมกับ getbyID พร้อมกัน
func TestCreateAndGetService(t *testing.T) {
	_, svc := setupTestDB(t)

	// Create
	newSvc := &barberBookingModels.Service{
		Name:     "Haircut",
		Duration: 30,
		Price:    200,
	}
	err := svc.CreateService(newSvc)
	assert.NoError(t, err)
	assert.NotZero(t, newSvc.ID)

	// Get by ID
	result, err := svc.GetServiceByID(newSvc.ID)
	assert.NoError(t, err)
	assert.Equal(t, "Haircut", result.Name)
	assert.Equal(t, 30, result.Duration)
}

//ค้นหา ID ไม่เจอ
func TestGetServiceByID_NotFound(t *testing.T) {
	_, svc := setupTestDB(t)

	// ดึง ID ที่ไม่มีในระบบ เช่น 999
	result, err := svc.GetServiceByID(999)

	assert.NoError(t, err)
	assert.Nil(t, result)
}

//ใส่ข้อมูลติดลบ 
func TestCreateService_InvalidInput(t *testing.T) {
	_, svc := setupTestDB(t)

	// duration = 0, price = -100 → invalid
	newSvc := &barberBookingModels.Service{
		Name:     "Invalid Service",
		Duration: 0,
		Price:    -100,
	}
	err := svc.CreateService(newSvc)

	// ต้อง error เพราะ field invalid
	assert.Error(t, err)
}

//ดึงข้อมูล service ทั้งหมด
func TestGetAllServices(t *testing.T) {
	_, svc := setupTestDB(t)

	// Seed 2 records
	svc.CreateService(&barberBookingModels.Service{Name: "Shampoo", Duration: 15, Price: 100})
	svc.CreateService(&barberBookingModels.Service{Name: "Beard Trim", Duration: 20, Price: 150})

	services, err := svc.GetAllServices()
	assert.NoError(t, err)
	assert.Len(t, services, 2)
}

//ดึงข้อมูล service ทั้งหมดแล้วไม่มีข้อมูล
func TestGetAllServices_Empty(t *testing.T) {
	_, svc := setupTestDB(t)

	services, err := svc.GetAllServices()
	assert.NoError(t, err)
	assert.Len(t, services, 0)
}

// ทดสอบว่าบริการที่ถูกลบจะไม่แสดง
func TestGetAllServices_SkipDeleted(t *testing.T) {
	db, svc := setupTestDB(t)

	// สร้างบริการ 2 รายการ
	svc.CreateService(&barberBookingModels.Service{Name: "Haircut", Duration: 30, Price: 200})
	svc.CreateService(&barberBookingModels.Service{Name: "Facial", Duration: 20, Price: 300})

	// ลบอันหนึ่ง
	_ = db.Delete(&barberBookingModels.Service{}, 1).Error

	services, err := svc.GetAllServices()
	assert.NoError(t, err)
	assert.Len(t, services, 1)
	assert.Equal(t, "Facial", services[0].Name)
}

//update ข้อมูล ปกติ
func TestUpdateService(t *testing.T) {
	_, svc := setupTestDB(t)

	svc.CreateService(&barberBookingModels.Service{Name: "Facial", Duration: 25, Price: 300})

	update := &barberBookingModels.Service{Name: "Facial Deluxe", Duration: 30, Price: 400}
	updated, err := svc.UpdateService(1, update)
	assert.NoError(t, err)
	assert.Equal(t, "Facial Deluxe", updated.Name)
	assert.Equal(t, 400.0, updated.Price)
}

//update ข้อมูลโดย ID ไม่ถูก
func TestUpdateService_NotFound(t *testing.T) {
	_, svc := setupTestDB(t)

	update := &barberBookingModels.Service{Name: "Ghost Service", Duration: 10, Price: 99}
	result, err := svc.UpdateService(999, update)

	assert.Nil(t, result)
	assert.Error(t, err)
}

//update ข้อมูลที่ผิด
func TestUpdateService_InvalidInput(t *testing.T) {
	_, svc := setupTestDB(t)

	// สร้างก่อน
	svc.CreateService(&barberBookingModels.Service{Name: "Massage", Duration: 60, Price: 500})

	// พยายามอัปเดตด้วย duration = 0
	invalid := &barberBookingModels.Service{Name: "Massage", Duration: 0, Price: 500}
	result, err := svc.UpdateService(1, invalid)

	assert.Nil(t, result)
	assert.Error(t, err)
}

// update แล้วต้องมี UpdatedAt ใหม่
func TestUpdateService_TimestampUpdated(t *testing.T) {
	_, svc := setupTestDB(t)

	svc.CreateService(&barberBookingModels.Service{Name: "Shave", Duration: 10, Price: 100})
	beforeUpdate, _ := svc.GetServiceByID(1)
	time.Sleep(1 * time.Second)

	// อัปเดตชื่อ
	update := &barberBookingModels.Service{Name: "Shave Premium", Duration: 10, Price: 100}
	afterUpdate, err := svc.UpdateService(1, update)

	assert.NoError(t, err)
	assert.True(t, afterUpdate.UpdatedAt.After(beforeUpdate.UpdatedAt))
}

//delete ปกติ
func TestDeleteService(t *testing.T) {
	_, svc := setupTestDB(t)

	svc.CreateService(&barberBookingModels.Service{Name: "Nose Wax", Duration: 10, Price: 120})

	err := svc.DeleteService(1)
	assert.NoError(t, err)

	svcAfterDelete, err := svc.GetServiceByID(1)
	assert.NoError(t, err)
	assert.Nil(t, svcAfterDelete)
}

//delete ID ไม่มี
func TestDeleteService_NotFound(t *testing.T) {
	_, svc := setupTestDB(t)

	err := svc.DeleteService(999) // ID ที่ไม่มี
	assert.Error(t, err)
}

//ลบข้อมูลเดิมซ้ำ
func TestDeleteService_Twice(t *testing.T) {
	_, svc := setupTestDB(t)

	svc.CreateService(&barberBookingModels.Service{Name: "Combo Pack", Duration: 45, Price: 700})
	err := svc.DeleteService(1)
	assert.NoError(t, err)

	// ลบซ้ำอีกครั้ง
	err = svc.DeleteService(1)
	assert.Error(t, err)
}

