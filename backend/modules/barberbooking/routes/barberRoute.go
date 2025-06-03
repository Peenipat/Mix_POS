package routes

import (
	"github.com/gofiber/fiber/v2"
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
)

func RegisterBarberRoutes(router fiber.Router, ctrl *barberBookingController.BarberController) {
	// ===== PUBLIC ROUTES (ไม่ต้อง Auth) =====
	
	router.Get("/barbers/:barber_id", ctrl.GetBarberByID)
	router.Get("/tenants/:tenant_id/barbers", barberbookingMiddlewares.RequireTenant(), ctrl.ListBarbersByTenant)

	// ===== PROTECTED ROUTES (ต้อง Auth + Tenant) =====
	group := router.Group("/tenants/:tenant_id/barbers")
	group.Use(middlewares.RequireAuth())
	group.Get("/branches/:branch_id/barbers", ctrl.ListBarbersByBranch)
	group.Get("/branches/:branch_id/user",ctrl.ListUserNotBarber)
	group.Post("/branches/:branch_id", ctrl.CreateBarber)
	group.Get("/users/:user_id/barber", barberbookingMiddlewares.RequireTenant(),ctrl.GetBarberByUser)
	group.Put("/:barber_id", barberbookingMiddlewares.RequireTenant(), ctrl.UpdateBarber)
	group.Delete("/:barber_id", barberbookingMiddlewares.RequireTenant(), ctrl.DeleteBarber)
	

	
}
