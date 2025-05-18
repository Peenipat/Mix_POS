package routes


import (
	// middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	// barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAppointmentStatusLogRoute(router fiber.Router ,ctrl *barberBookingController.AppointmentStatusLogController ){
	group := router.Group("/tenants/:tenant_id/appointments")
	group.Get("/:appointment_id/logs", ctrl.GetAppointmentLogs) 
	

}