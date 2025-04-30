package admin

import (
	Core_controllers "myapp/modules/core/controllers"
	"myapp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func SetupAdminRoutes(app *fiber.App) {
	adminGroup := app.Group("/admin", middlewares.RequireAuth(), middlewares.RequireSuperAdmin())
	adminGroup.Post("/create_users", Core_controllers.CreateUserFromAdmin)
	adminGroup.Put("/change_role", Core_controllers.ChangeUserRole)
	adminGroup.Get("/users",Core_controllers.GetAllUsers)
	adminGroup.Get("/user-by-role",Core_controllers.FilterUsersByRole)
	
	// adminGroup.Get("/system_logs", controllers.GetSystemLogs)
	adminGroup.Post("/system_logs",           Core_controllers.CreateLog)
    // ดึงรายการ log พร้อมกรอง และ pagination
    adminGroup.Get("/system_logs",            Core_controllers.GetSystemLogs)
    // ดูรายละเอียด log ทีละรายการ ตาม log_id
    adminGroup.Get("/system_logs/:log_id",    Core_controllers.GetSystemLogByID)
}


