package routes

import (
    "github.com/gofiber/fiber/v2"

    coreControllers "myapp/modules/core/controllers"
    middlewares "myapp/middlewares"
	coremiddlewares "myapp/modules/core/middlewares"
)

func RegisterBranchRoutes(router fiber.Router, ctrl *coreControllers.BranchController) {
    // ===== PUBLIC ROUTES (no tenant context) =====
    router.Get("/branches", ctrl.GetBranches)           // list all branches (admin/global)
    router.Get("/branches/:id", ctrl.GetBranchByID)     

    // ===== TENANT-SCOPED ROUTES (requires auth + tenant middleware) =====
    tenantGroup := router.Group("/tenants/:tenant_id")
    tenantGroup.Use(middlewares.RequireAuth(), coremiddlewares.RequireTenant())

    tenantGroup.Post("/branches", ctrl.CreateBranch)
    tenantGroup.Get("/branches", ctrl.GetBranchesByTenantID)
    tenantGroup.Put("/branches/:id", ctrl.UpdateBranch)
    tenantGroup.Delete("/branches/:id", ctrl.DeleteBranch)
}
