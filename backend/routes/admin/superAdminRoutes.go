package admin

import (
	"myapp/controllers"
	"myapp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin", middlewares.RequireAuth(), middlewares.RequireSuperAdmin())
	adminGroup.Post("/create_users", controllers.CreateUserFromAdmin)
	adminGroup.Put("/change_role", controllers.ChangeUserRole)
	adminGroup.Get("/users",controllers.GetAllUsers)
	adminGroup.Get("/user-by-role",controllers.FilterUsersByRole)
	// adminGroup.Get("/system_logs", controllers.GetSystemLogs)
	adminGroup.Post("/system_logs",           controllers.CreateLog)
    // ดึงรายการ log พร้อมกรอง และ pagination
    adminGroup.Get("/system_logs",            controllers.GetSystemLogs)
    // ดูรายละเอียด log ทีละรายการ ตาม log_id
    adminGroup.Get("/system_logs/:log_id",    controllers.GetSystemLogByID)
}


