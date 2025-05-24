package coreRoutes

import (
	"github.com/gofiber/fiber/v2"
	Core_controllers "myapp/modules/core/controllers"
	middlewares "myapp/middlewares"
)


func RegisterTenantRoutes(router fiber.Router, ctrl *Core_controllers.TenantController) {
	tenant := router.Group("/tenant-route")
	tenant.Use(middlewares.RequireAuth())
	tenant.Post("/create",ctrl.CreateTenant)
	tenant.Get("/:id", ctrl.GetTenantByID)
	tenant.Get("/",ctrl.ListTenants)
	tenant.Put("/:id",ctrl.UpdateTenant)
	tenant.Delete("/:id",ctrl.DeleteTenant)
}