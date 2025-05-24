package coreRoutes

import (
	"github.com/gofiber/fiber/v2"
	Core_controllers "myapp/modules/core/controllers"
)


func RegisterUserRoutes(router fiber.Router, ctrl *Core_controllers.UserController) {
	user := router.Group("/user")
	user.Put("/change-password/:id", ctrl.ChangePassword)


}
