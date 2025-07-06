package routes

import (
	"github.com/gofiber/fiber/v2"
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	// barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
)

func RegisterBarberRoutes(router fiber.Router, ctrl *barberBookingController.BarberController) {
	
	router.Get("/barbers/:barber_id", ctrl.GetBarberByID)
	router.Get("/branches/:branch_id/barbers", ctrl.ListBarbersByBranch)
	router.Get("/tenants/:tenant_id/barbers" ,ctrl.ListBarbersByTenant)

	router.Put("/tenants/:tenant_id/barbers/:barber_id/update-barber",middlewares.RequireAuth(),ctrl.UpdateBarber,)
	router.Delete("/barbers/:barber_id",middlewares.RequireAuth(), ctrl.DeleteBarber)
	router.Get("/users/:user_id/barber",middlewares.RequireAuth(),ctrl.GetBarberByUser)

	group := router.Group("/tenants/:tenant_id/branches/:branch_id")
	group.Use(middlewares.RequireAuth())
	
	group.Get("/users",ctrl.ListUserNotBarber) //ปัญหาคือยังไม่ได้คิดเรื่อง tenant id ไว้
	group.Post("/create-barber", ctrl.CreateBarber)
	group.Get("/users/:user_id/barber",ctrl.GetBarberByUser)


	

	
}
