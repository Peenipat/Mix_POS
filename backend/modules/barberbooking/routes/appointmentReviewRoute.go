package routes

import (
	middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAppointmentReviewRoute(router fiber.Router, ctrl *barberBookingController.AppointmentReviewController) {

	group := router.Group("/tenants/:tenant_id/appointments")
	group.Post("/:appointment_id/reviews",barberbookingMiddlewares.RequireTenant(), ctrl.CreateReview)
	group.Put("/reviews/:review_id",barberbookingMiddlewares.RequireTenant(), ctrl.UpdateReview)

	group.Use(middlewares.RequireAuth())
    group.Get("/tenants/:tenant_id/reviews/:review_id", barberbookingMiddlewares.RequireTenant(),ctrl.GetReviewByID) //แก้เรื่อง role หน่อยนะ 
}
