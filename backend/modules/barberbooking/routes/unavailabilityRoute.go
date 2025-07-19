package routes

import (
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)
func RegisterUnavailabilityRoute(router fiber.Router ,ctrl *barberBookingController.UnavailabilityController){
	group := router.Group("/tenants/:tenant_id/unavailability")

	group.Get("/branches/:branch_id", ctrl.GetUnavailabilitiesByBranch) // ลูกค้าดูวันหยุดของสาขา
	group.Get("/barbers/:barber_id", ctrl.GetUnavailabilitiesByBarber)  // ลูกค้าดูวันหยุดของช่าง
	group.Use(middlewares.RequireAuth())

	group.Post("/",barberbookingMiddlewares.RequireTenant(),ctrl.CreateUnavailability)
	group.Put("/:id",barberbookingMiddlewares.RequireTenant(),ctrl.UpdateUnavailability)
	group.Delete("/:id",barberbookingMiddlewares.RequireTenant(),ctrl.DeleteUnavailability)
	
}