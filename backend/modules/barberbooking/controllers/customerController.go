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
	customerList, err := ctrl.CustomerService.GetAllCustomers(c.Context(), uint(tenantId))
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
