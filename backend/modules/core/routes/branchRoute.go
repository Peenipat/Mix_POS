package coreRoutes

import (
    "github.com/gofiber/fiber/v2"

    coreControllers "myapp/modules/core/controllers"
    middlewares "myapp/middlewares"
	coremiddlewares "myapp/modules/core/middlewares"
)

func RegisterBranchRoutes(router fiber.Router, ctrl *coreControllers.BranchController) {
       
    router.Get("/branches/all",middlewares.RequireAuth(),ctrl.GetBranches)  
    router.Get("/branch/:id",middlewares.RequireAuth() ,ctrl.GetBranchByID) 

    tenantGroup := router.Group("/tenants/:tenant_id")
    tenantGroup.Use(middlewares.RequireAuth(), coremiddlewares.RequireTenant())
    tenantGroup.Get("/branches", ctrl.GetBranchesByTenantID)
    tenantGroup.Post("/branches", ctrl.CreateBranch)
    tenantGroup.Put("/branches/:id", ctrl.UpdateBranch)
    tenantGroup.Delete("/branches/:id", ctrl.DeleteBranch)
}
