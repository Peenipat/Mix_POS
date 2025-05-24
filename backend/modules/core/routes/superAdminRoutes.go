package coreRoutes

import (
	Core_controllers "myapp/modules/core/controllers"
	"myapp/middlewares"

	"github.com/gofiber/fiber/v2"
)

func RegisterAdminRoutes(router fiber.Router, ctrl *Core_controllers.UserController) {
	adminGroup := router.Group("/admin", middlewares.RequireAuth(), middlewares.RequireSuperAdmin())
	adminGroup.Post("/create_users", ctrl.CreateUserFromAdmin)
	adminGroup.Put("/change_role", ctrl.ChangeUserRole)
	adminGroup.Get("/users",ctrl.GetAllUsers)
	adminGroup.Get("/user-by-role",ctrl.FilterUsersByRole)
	
	// adminGroup.Get("/system_logs", controllers.GetSystemLogs)
	adminGroup.Post("/system_logs",           Core_controllers.CreateLog)
    // ดึงรายการ log พร้อมกรอง และ pagination
    adminGroup.Get("/system_logs",            Core_controllers.GetSystemLogs)
    // ดูรายละเอียด log ทีละรายการ ตาม log_id
    adminGroup.Get("/system_logs/:log_id",    Core_controllers.GetSystemLogByID)
}


