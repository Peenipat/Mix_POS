package barberBookingController

import (
	"fmt"
	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
	"net/http"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
)

type CustomerController struct {
	CustomerService barberBookingPort.ICustomer
}

func NewCustomerController(scv barberBookingPort.ICustomer) *CustomerController {
	return &CustomerController{
		CustomerService: scv,
	}
}

var RolesCanManageCustomer = []coreModels.RoleName{
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
	coreModels.RoleNameBranchAdmin,
}

// GetAllCustomers godoc
// @Summary      ดึงรายชื่อลูกค้า
// @Description  คืนรายการ Customer ทั้งหมดของ Tenant ที่ระบุ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                               true  "รหัส Tenant"
// @Success      200        {object}  map[string]interface{}            "คืนค่า status success, message และ array ของ Customer ใน key `data`"
// @Failure      400        {object}  map[string]string                 "Invalid tenant ID หรือ Failed to fetch customer"
// @Failure      403        {object}  map[string]string                 "Permission denied"
// @Router       /tenants/:tenant_id/branch/:branch_id/customers [get]
// @Security     ApiKeyAuth
func (ctrl *CustomerController) GetAllCustomers(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	fmt.Println("roleStr:", roleStr)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageCustomer) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	tenantId, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant ID",
			"error":   err.Error(),
		})
	}

	branchId, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch ID",
			"error":   err.Error(),
		})
	}
	customerList, err := ctrl.CustomerService.GetAllCustomers(c.Context(), uint(tenantId),uint(branchId))
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch customer",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Customer retrieved",
		"data":    customerList,
	})
}

// GetCustomerByID godoc
// @Summary      ดึงข้อมูลลูกค้าตาม ID
// @Description  คืนข้อมูล Customer ตามรหัสที่ระบุ ภายใต้ Tenant ที่ระบุ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                             true  "รหัส Tenant"
// @Param        cus_id     path      uint                             true  "รหัส Customer"
// @Success      200        {object}  map[string]interface{}          "คืนค่า status success, message และข้อมูล Customer ใน key `data`"
// @Failure      400        {object}  map[string]string               "Invalid tenant_id หรือ invalid cus_id"
// @Failure      403        {object}  map[string]string               "Permission denied"
// @Failure      404        {object}  map[string]string               "Customer not found"
// @Failure      500        {object}  map[string]string               "Failed to fetch customer"
// @Router      /tenants/:tenant_id/customers/:cus_id [get]
// @Security     ApiKeyAuth
func (ctrl *CustomerController) GetCustomerByID(c *fiber.Ctx) error {
	// 1. เช็คสิทธิ์
	roleStr, ok := c.Locals("role").(string)
	fmt.Println("roleStr:", roleStr)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageCustomer) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// 2. Parse tenant_id
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// 3. Parse and validate cus_id (single pass)
	cusParam := c.Params("cus_id")
	id64, err := strconv.ParseUint(cusParam, 10, 64)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "invalid cus_id"})
	}
	customerID := uint(id64)

	// 4. เรียก Service
	customer, err := ctrl.CustomerService.GetCustomerByID(c.Context(), tenantID, customerID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch customer",
			"error":   "Internal server error",
		})
	}

	// 5. Not found
	if customer == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Customer not found",
		})
	}

	// 6. Success
	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Customer retrieved",
		"data":    customer,
	})
}

// CreateCustomer godoc
// @Summary      สร้างลูกค้าใหม่
// @Description  เพิ่ม Customer ใหม่ภายใต้ Tenant ที่ระบุ (กรอกชื่อและเบอร์โทร)  
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                          true  "รหัส Tenant"
// @Param        body       body      barberBookingModels.Customer  true  "Payload สำหรับสร้าง Customer (Name, Phone)"
// @Success      201        {object}  map[string]string            "คืนค่า status success และข้อความยืนยันการสร้าง"
// @Failure      400        {object}  map[string]string            "Invalid tenant ID หรือ Invalid request body หรือ Invalid Customer input"
// @Failure      500        {object}  map[string]string            "Failed to create customer"
// @Router       /tenants/:tenant_id/customers [post]
// @Security     ApiKeyAuth
func (ctrl *CustomerController) CreateCustomer(c *fiber.Ctx) error {

	var payload barberBookingModels.Customer
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant ID",
		})
	}
	payload.TenantID = tenantID
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" || len(payload.Name) >= 100 || len(payload.Phone) != 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Customer input",
		})
	}

	if err := ctrl.CustomerService.CreateCustomer(c.Context(), &payload); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to create customer",
			"error":   "Can't create Customer",
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer created",
	})

}

// UpdateCustomer godoc
// @Summary      แก้ไขข้อมูลลูกค้า
// @Description  อัปเดตข้อมูล Customer ตามรหัสที่ระบุ ภายใต้ Tenant ที่ระบุ
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                           true  "รหัส Tenant"
// @Param        cus_id     path      uint                           true  "รหัส Customer"
// @Param        body       body      barberBookingModels.Customer   true  "Payload สำหรับอัปเดต Customer (Name, Phone, Email)"
// @Success      200        {object}  barberBookingModels.Customer   "คืนค่า status success และข้อมูล Customer ที่อัปเดตใน key `data`"
// @Failure      400        {object}  map[string]string              "Invalid tenant_id, invalid cus_id, Invalid request body หรือ Invalid Customer input"
// @Failure      404        {object}  map[string]string              "Customer not found"`
// @Failure      500        {object}  map[string]string              "Failed to update customer"
// @Router       /tenants/:tenant_id/customers/:cus_id  [put]
// @Security     ApiKeyAuth
func (ctrl *CustomerController) UpdateCustomer(c *fiber.Ctx) error {
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	customerID, err := helperFunc.ParseUintParam(c, "cus_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var payload barberBookingModels.Customer
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	payload.Name = strings.TrimSpace(payload.Name)
	if payload.Name == "" || len(payload.Name) >= 100 || len(payload.Phone) != 10 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid Customer input",
		})
	}

	existingCustomer, err := ctrl.CustomerService.GetCustomerByID(c.Context(), tenantID, customerID)
	if err != nil || existingCustomer == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Customer not found",
		})
	}

	existingCustomer.Name = payload.Name
	existingCustomer.Phone = payload.Phone
	existingCustomer.Email = payload.Email // optional

	updatedCustomer, err := ctrl.CustomerService.UpdateCustomer(c.Context(), tenantID, customerID, existingCustomer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to update customer",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "Customer Updated",
		"data":    updatedCustomer,
	})
}

// DeleteCustomer godoc
// @Summary      ลบลูกค้า
// @Description  ลบ Customer ตามรหัสที่ระบุ ภายใต้ Tenant ที่ระบุ (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint               true  "รหัส Tenant"
// @Param        cus_id     path      uint               true  "รหัส Customer"
// @Success      200        {object}  map[string]string  "คืนค่า status success และข้อความยืนยันการลบ"
// @Failure      400        {object}  map[string]string  "Invalid tenant_id หรือ invalid cus_id"
// @Failure      403        {object}  map[string]string  "Permission denied"
// @Failure      500        {object}  map[string]string  "Failed to delete customer"
// @Router       /tenants/:tenant_id/customers/:cus_id [delete]
// @Security     ApiKeyAuth
func (ctrl *CustomerController) DeleteCustomer(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageCustomer) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	customerID, err := helperFunc.ParseUintParam(c, "cus_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	if err := ctrl.CustomerService.DeleteCustomer(c.Context(), tenantID, customerID); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete customer",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Customer delete successfully",
	})

}


// FindCustomerByEmail godoc
// @Summary      ค้นหาลูกค้าตามอีเมล
// @Description  ค้นหา Customer ภายใต้ Tenant ที่ระบุ โดยใช้ Email ใน body (ต้องมีสิทธิ์ SaaSSuperAdmin, Tenant, TenantAdmin หรือ BranchAdmin)
// @Tags         Customer
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint                          true  "รหัส Tenant"
// @Param        body       body      barberBookingModels.Customer true  "Payload ที่ประกอบด้วยฟิลด์ Email (json: email)"
// @Success      200        {object}  map[string]interface{}       "คืนค่า status success, message และข้อมูล Customer ใน key `data`"
// @Failure      400        {object}  map[string]string            "Invalid tenant_id หรือ Invalid request body"
// @Failure      403        {object}  map[string]string            "Permission denied"
// @Failure      404        {object}  map[string]string            "Customer not found"
// @Failure      500        {object}  map[string]string            "Failed to find customer by email"
// @Router       /tenants/:tenant_id/customers/find-by-email [post]
// @Security     ApiKeyAuth
func (ctrl *CustomerController) FindCustomerByEmail(c *fiber.Ctx) error {
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageCustomer) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	var payload barberBookingModels.Customer
	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	email := strings.ToLower(strings.TrimSpace(payload.Email))

	customer, err := ctrl.CustomerService.FindCustomerByEmail(c.Context(), tenantID, email)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to delete Customer",
			"error":   err.Error(),
		})
	}

	if customer == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"status":  "error",
			"message": "Customer not found",
		})
	}

	return c.JSON(fiber.Map{
		"status":  "success",
		"message": "Customer retrieved",
		"data":    customer,
	})
}


func (ctrl *CustomerController) GetPendingAndCancelledByCustomer(c *fiber.Ctx) error {
    // 1. ดึงค่า tenant_id และ branch_id จาก URL parameter
    // Fiber ใช้ c.Params("name") หรือ c.Params("name", "defaultValue") ได้
    tenantIDParam := c.Params("tenant_id")
    branchIDParam := c.Params("branch_id")
	customerIDParam := c.Params("cus_id")
	

    tenantID64, err := strconv.ParseUint(tenantIDParam, 10, 0)
    if err != nil {
        // คืน HTTP 400 พร้อม JSON error
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid tenant_id parameter",
        })
    }
    branchID64, err := strconv.ParseUint(branchIDParam, 10, 0)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid branch_id parameter",
        })
    }
	customerID64, err := strconv.ParseUint(customerIDParam, 10, 0)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "invalid cus_id parameter",
        })
    }
    tenantID := uint(tenantID64)
    branchID := uint(branchID64)
	customerID := uint(customerID64)

    // 2. เรียก service ให้ทำงาน
    results, err := ctrl.CustomerService.GetPendingAndCancelledCount(
        c.Context(),
        tenantID,
        branchID,
		customerID,
    )
    if err != nil {
        // 500 Internal Server Error
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": err.Error(),
        })
    }

    // 3. ส่งผลลัพธ์กลับ (status success + data)
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   results,
    })
}

