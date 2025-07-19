package coreRoutes

import (
	"github.com/gofiber/fiber/v2"
	Core_controllers "myapp/modules/core/controllers"
	middlewares "myapp/middlewares"
)


func RegisterTenantUserRoutes(router fiber.Router, ctrl *Core_controllers.TenantUserController) {
	tenantuser := router.Group("/tenant-user")
	tenantuser.Use(middlewares.RequireAuth())
	tenantuser.Post("/tenants/:tenant_id/users/:user_id", ctrl.AddUserToTenant)
	tenantuser.Delete("/tenants/:tenant_id/users/:user_id", ctrl.RemoveUserFromTenant)
	tenantuser.Get("/user/:user_id",ctrl.ListTenantsByUser)

	tenantuser.Get("/tenants/:tenant_id",ctrl.ListUsersForTenant)

}

