package routes

import (
	"github.com/gofiber/fiber/v2"
	"myapp/controllers"
	"myapp/middlewares"
)

func SetupAuthRoutes(app *fiber.App) {
	auth := app.Group("/auth")
	auth.Post("/register", controllers.CreateUserFromRegister)
	auth.Post("/login", controllers.LoginHandler)
	app.Get("/auth/me", middlewares.RequireAuth(), controllers.GetMe)

}
