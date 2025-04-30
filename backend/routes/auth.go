package routes

import (
	"github.com/gofiber/fiber/v2"
	Core_controllers "myapp/modules/core/controllers"
)


func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", Core_controllers.CreateUserFromRegister)
	auth.Post("/login", Core_controllers.LoginHandler)

}
