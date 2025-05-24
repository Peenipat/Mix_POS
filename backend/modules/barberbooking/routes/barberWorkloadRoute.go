package routes

import (
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterBarberWorkloadRoute(router fiber.Router ,ctrl barberBookingController.BarberWorkloadController){
	group := router.Group("/tenants/:tenant_id/barberworkload")

	group.Get("/barbers/:barber_id", ctrl.GetWorkloadByBarber)
    group.Get("/branches/:branch_id/summary", ctrl.GetWorkloadSummaryByBranch)

    // Any write operation ต้องผ่าน auth + tenant check
    group.Use(middlewares.RequireAuth())
    group.Post("/barbers/:barber_id", barberbookingMiddlewares.RequireTenant(), ctrl.UpsertBarberWorkload)
    group.Put("/barbers/:barber_id", barberbookingMiddlewares.RequireTenant(), ctrl.UpsertBarberWorkload)
	
}