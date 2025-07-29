package routes

import (
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	// barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)


func RegisterWorkingHourRoute(router fiber.Router ,ctrl barberBookingController.WorkingHourController){
	group := router.Group("/tenants/:tenant_id/workinghour")
	
	group.Get("/branches/:branch_id",ctrl.GetWorkingHours)
	group.Get("/branches/:branch_id/slots",ctrl.GetAvailableSlots)

	group.Use(middlewares.RequireAuth())
	group.Post("branches/:branch_id",ctrl.CreateWorkingHours)
	group.Put("branches/:branch_id",ctrl.UpdateWorkingHours)
}