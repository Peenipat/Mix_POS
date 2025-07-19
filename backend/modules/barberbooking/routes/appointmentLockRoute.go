package routes

import (
	barberBookingController "myapp/modules/barberbooking/controllers"

	"github.com/gofiber/fiber/v2"
)

func RegisterAppointmentLockRoute(router fiber.Router, ctrl *barberBookingController.AppointmentLockController) {
	group := router.Group("/tenants/:tenant_id/branches/:branch_id/appointments-lock")

	group.Post("/", ctrl.CreateAppointmentLock)
	group.Delete("/:lock_id", ctrl.ReleaseAppointmentLock)
	group.Get("/", ctrl.GetAppointmentLocks)
}