package coreRoutes

import (
	"github.com/gofiber/fiber/v2"
	Core_controllers "myapp/modules/core/controllers"
)


func RegisterTenantUserRoutes(router fiber.Router, ctrl *Core_controllers.TenantUserController) {
	tenantuser := router.Group("/tenant-user")
	tenantuser.Post("/tenants/:tenant_id/users/:user_id", ctrl.AddUserToTenant)
	tenantuser.Delete("/tenants/:tenant_id/users/:user_id", ctrl.RemoveUserFromTenant)

}

