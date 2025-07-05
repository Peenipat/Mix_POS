package routes

import (
	// middlewares "myapp/middlewares"
	barberBookingControllerMix "myapp/modules/barberbooking/controllers"
	// barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterCalendarRoute(router fiber.Router, ctrl *barberBookingControllerMix.CalendarController) {
	group := router.Group("/tenants/:tenant_id/branches/:branch_id")
	// group.Use(middlewares.RequireAuth())
	// group.Use(barberbookingMiddlewares.RequireTenant())
	group.Get("/available-slots", ctrl.GetAvailableSlots)
}
