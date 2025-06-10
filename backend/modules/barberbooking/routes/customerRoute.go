package routes

import (
	"github.com/gofiber/fiber/v2"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
	middlewares "myapp/middlewares"
	
)

func RegisterCustomerRoutes(router fiber.Router, ctrl *barberBookingController.CustomerController) {
	
	group := router.Group("/tenants/:tenant_id/branch/:branch_id/customers")
	//  PUBLIC route – ลูกค้าสมัครเองได้
	group.Post("/", ctrl.CreateCustomer)
	group.Use(middlewares.RequireAuth())

	
	group.Get("/:cus_id/appointments", ctrl.GetPendingAndCancelledByCustomer)

	// PROTECTED route – ต้องมี role (admin/manager) ผ่าน middleware
	group.Get("/", barberbookingMiddlewares.RequireTenant(), ctrl.GetAllCustomers)
	group.Get("/:cus_id", barberbookingMiddlewares.RequireTenant(), ctrl.GetCustomerByID)
	group.Put("/:cus_id", barberbookingMiddlewares.RequireTenant(), ctrl.UpdateCustomer)
	group.Delete("/:cus_id", barberbookingMiddlewares.RequireTenant(), ctrl.DeleteCustomer)
	

	//  PROTECTED: หา customer จาก email
	group.Post("/find-by-email", barberbookingMiddlewares.RequireTenant(), ctrl.FindCustomerByEmail)
}
