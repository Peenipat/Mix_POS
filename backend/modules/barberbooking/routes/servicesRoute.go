// barberbooking/routes/service_routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
	middlewares "myapp/middlewares"
)

func RegisterServiceRoutes(router fiber.Router, ctrl *barberBookingController.ServiceController) {

	group := router.Group("/tenants/:tenant_id/services")

	group.Get("/", ctrl.GetAllServices)    //  public
	group.Get("/:service_id", ctrl.GetServiceByID) //  public

	group.Use(middlewares.RequireAuth())
	group.Post("/", barberbookingMiddlewares.RequireTenant(), ctrl.CreateService)
	group.Put("/:service_id", barberbookingMiddlewares.RequireTenant(), ctrl.UpdateService)
	group.Delete("/:service_id", barberbookingMiddlewares.RequireTenant(), ctrl.DeleteService)
}
