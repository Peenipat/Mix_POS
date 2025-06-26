package routes

import (
	"github.com/gofiber/fiber/v2"
	// middlewares "myapp/middlewares"
	barberBookingController "myapp/modules/barberbooking/controllers"
	// barberbookingMiddlewares "myapp/modules/barberbooking/middlewares"
)

func RegisterWorkingDayOverrideRoutes(router fiber.Router, ctrl *barberBookingController.WorkingDayOverrideController) {
	router.Post("/working-day-overrides", ctrl.Create)
	router.Put("/working-day-overrides/:id", ctrl.Update)
	router.Get("/working-day-overrides/:id", ctrl.GetByID)
	router.Delete("/working-day-overrides/:id", ctrl.DeleteWorkingDayOverride)

}