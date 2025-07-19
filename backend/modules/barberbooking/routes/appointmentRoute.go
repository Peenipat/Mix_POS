package routes

import (
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAppointmentRoute(router fiber.Router ,ctrl *barberBookingController.AppointmentController){
	router.Get("/branches/:branch_id/appointments",ctrl.GetAppointmentsByBranch)
	group := router.Group("/tenants/:tenant_id/appointments")
	group.Get("/", ctrl.ListAppointments) //รอเช็คเรื่อง not_found //
	group.Get("/barbers/:barber_id/availability", ctrl.CheckBarberAvailability) //ยังไม่ผ่านใน dev
	group.Get("/branches/:branch_id/available-barbers",ctrl.GetAvailableBarbers,)
	group.Get("/:appointment_id", ctrl.GetAppointmentByID)//

	group.Post("/", ctrl.CreateAppointment)//
	group.Post("/:appointment_id/cancel", ctrl.CancelAppointment) //ขาดการเช็คเรื่อง not_found cus_id , user_id //
	group.Post("/:appointment_id/reschedule", ctrl.RescheduleAppointment) //

	group.Put("/:appointment_id", ctrl.UpdateAppointment) // เหมือน log จะไม่ update //
	
	group.Use(middlewares.RequireAuth())
	group.Delete("/:appointment_id", barberbookingMiddlewares.RequireTenant(), ctrl.DeleteAppointment)//

}