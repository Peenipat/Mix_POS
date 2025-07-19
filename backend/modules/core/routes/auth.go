package coreRoutes

import (
	"github.com/gofiber/fiber/v2"
	Core_controllers "myapp/modules/core/controllers"
)


func SetupAuthRoutes(router fiber.Router, ctrl *Core_controllers.UserController) {
	auth := router.Group("/auth")
	auth.Post("/register", ctrl.CreateUserFromRegister)
	auth.Post("/login", Core_controllers.LoginHandler)

}
