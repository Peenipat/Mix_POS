// barberbooking/routes/service_routes.go
package routes

import (
	"github.com/gofiber/fiber/v2"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
)

func RegisterServiceRoutes(router fiber.Router, ctrl *barberBookingController.ServiceController) {

	group := router.Group("/services")

	group.Get("/", ctrl.GetAllServices)    //  public
	group.Get("/:id", ctrl.GetServiceByID) //  public

	group.Post("/", barberbookingMiddlewares.RequireTenant(), ctrl.CreateService)
	group.Put("/:id", barberbookingMiddlewares.RequireTenant(), ctrl.UpdateService)
	group.Delete("/:id", barberbookingMiddlewares.RequireTenant(), ctrl.DeleteService)
}
