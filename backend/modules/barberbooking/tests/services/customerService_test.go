package barberbookingServiceTest

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	bookingModels "myapp/modules/barberbooking/models"
	bookingServices "myapp/modules/barberbooking/services"
)

func setupTestCustomerDB(t *testing.T) *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	assert.NoError(t, err)

	err = db.AutoMigrate(&bookingModels.Customer{})
	assert.NoError(t, err)

	return db
}

func TestCustomerService_CRUD(t *testing.T) {
	ctx := context.Background()
	db := setupTestCustomerDB(t)
	svc := bookingServices.NewCustomerService(db)

	tenantID := uint(1)

	t.Run("CreateCustomer", func(t *testing.T) {
		customer := &bookingModels.Customer{
			TenantID:  tenantID,
			Name:      "Alice",
			Email:     "alice@example.com",
			Phone:     "0801234567",
			CreatedAt: time.Now(),
		}
		err := svc.CreateCustomer(ctx, customer)
		assert.NoError(t, err)
		assert.NotZero(t, customer.ID)
	})

	t.Run("CreateCustomer_DuplicateEmail", func(t *testing.T) {
		cust := &bookingModels.Customer{
			TenantID: tenantID,
			Name:     "Alice Clone",
			Email:    "alice@example.com", // ซ้ำ
		}
		err := svc.CreateCustomer(ctx, cust)
		assert.Error(t, err)
	})

	// t.Run("GetAllCustomers", func(t *testing.T) {
	// 	customers, err := svc.GetAllCustomers(ctx, tenantID)
	// 	assert.NoError(t, err)
	// 	assert.GreaterOrEqual(t, len(customers), 1)
	// })

	t.Run("GetCustomerByID", func(t *testing.T) {
		customer := &bookingModels.Customer{
			TenantID: tenantID,
			Name:     "Bob",
			Email:    "bob@example.com",
		}
		_ = svc.CreateCustomer(ctx, customer)

		found, err := svc.GetCustomerByID(ctx, tenantID, customer.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Bob", found.Name)
	})

	t.Run("UpdateCustomer", func(t *testing.T) {
		customer := &bookingModels.Customer{
			TenantID: tenantID,
			Name:     "Charlie",
			Email:    "charlie@example.com",
		}
		_ = svc.CreateCustomer(ctx, customer)

		customerNew := &bookingModels.Customer{Name: "Charlie Updated"}
		CustomerUpdate, err := svc.UpdateCustomer(ctx, tenantID, customer.ID, customerNew)
		assert.NoError(t, err)
		assert.Equal(t, "Charlie Updated", CustomerUpdate.Name)
	})

	t.Run("DeleteCustomer", func(t *testing.T) {
		customer := &bookingModels.Customer{
			TenantID: tenantID,
			Name:     "Delete Me",
			Email:    "deleteme@example.com",
		}
		_ = svc.CreateCustomer(ctx, customer)

		err := svc.DeleteCustomer(ctx, tenantID, customer.ID)
		assert.NoError(t, err)

		found, err := svc.GetCustomerByID(ctx, tenantID, customer.ID)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	// t.Run("CreateCustomer_MissingEmail", func(t *testing.T) {
	// 	cust := &bookingModels.Customer{
	// 		TenantID: tenantID,
	// 		Name:     "No Email",
	// 	}
	// 	err := svc.CreateCustomer(ctx, cust)
	// 	assert.Error(t, err)
	// })

	t.Run("CreateCustomer_MissingTenant", func(t *testing.T) {
		cust := &bookingModels.Customer{
			Name:  "No Tenant",
			Email: "notenant@example.com",
		}
		err := svc.CreateCustomer(ctx, cust)
		assert.Error(t, err)
	})

	t.Run("GetCustomerByID_NotFound", func(t *testing.T) {
		found, err := svc.GetCustomerByID(ctx, tenantID, 999999)
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

	t.Run("UpdateCustomer_NotFound", func(t *testing.T) {
		cust := &bookingModels.Customer{
			Name: "No Name",
		}
		CustomerUpdate, err := svc.UpdateCustomer(ctx, tenantID, 999999, cust)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
		assert.Nil(t, CustomerUpdate)
		assert.Error(t, err)
	})

	t.Run("UpdateCustomer_EmptyData", func(t *testing.T) {
		customer := &bookingModels.Customer{
			TenantID: tenantID,
			Name:     "Target Empty Update",
			Email:    "emptyupdate@example.com",
		}
		_ = svc.CreateCustomer(ctx, customer)
		customerEmptyData:= &bookingModels.Customer{}
		cus,err := svc.UpdateCustomer(ctx, tenantID, customer.ID,customerEmptyData)
		assert.Error(t, err) // GORM updates nothing, no error
		assert.Nil(t,cus)
		assert.Contains(t, err.Error(), "name must be 2 - 100 characters")
	})

	t.Run("DeleteCustomer_NotFound", func(t *testing.T) {
		err := svc.DeleteCustomer(ctx, tenantID, 999999)
		assert.ErrorIs(t, err, gorm.ErrRecordNotFound)
	})

	t.Run("FindCustomerByEmail_Found", func(t *testing.T) {
		cust := &bookingModels.Customer{
			TenantID: tenantID,
			Name:     "Email Lookup",
			Email:    "lookup@example.com",
		}
		_ = svc.CreateCustomer(ctx, cust)

		found, err := svc.FindCustomerByEmail(ctx, tenantID, "lookup@example.com")
		assert.NoError(t, err)
		assert.NotNil(t, found)
		assert.Equal(t, "Email Lookup", found.Name)
	})

	t.Run("FindCustomerByEmail_NotFound", func(t *testing.T) {
		found, err := svc.FindCustomerByEmail(ctx, tenantID, "notfound@example.com")
		assert.NoError(t, err)
		assert.Nil(t, found)
	})

}
