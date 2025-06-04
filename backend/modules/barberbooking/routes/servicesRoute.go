// barberbooking/routes/service_routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
	middlewares "myapp/middlewares"
)

func RegisterServiceRoutes(router fiber.Router, ctrl *barberBookingController.ServiceController) {

	router.Delete("/services/:service_id",middlewares.RequireAuth(),ctrl.DeleteService)

	group := router.Group("/tenants/:tenant_id/branch/:branch_id/services")

	group.Get("/", ctrl.GetAllServices)    //  public
	group.Get("/:service_id", ctrl.GetServiceByID) //  public

	group.Use(middlewares.RequireAuth())
	group.Post("/", ctrl.CreateService)
	group.Put("/:service_id", barberbookingMiddlewares.RequireTenant(), ctrl.UpdateService)
	
}
