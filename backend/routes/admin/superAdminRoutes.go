package admin

import (
	"github.com/gofiber/fiber/v2"
	"myapp/controllers"
	"myapp/middlewares"
	
)

func SetupAdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin", middlewares.RequireAuth(), middlewares.RequireSuperAdmin())
	adminGroup.Post("/create_users", controllers.CreateUserFromAdmin)
}
