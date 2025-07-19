package corePort
import (
	"context"
	coreModels "myapp/modules/core/models"
)

type CreateTenantInput struct {
    Name   string 
    Domain string 
}

type UpdateTenantInput struct {
    ID       uint    
    Name     *string 
    Domain   *string 
    IsActive *bool  
}


type ITenant interface {
    CreateTenant(ctx context.Context, input CreateTenantInput) (*coreModels.Tenant, error)
    GetTenantByID(ctx context.Context, id uint) (*coreModels.Tenant, error)
    UpdateTenant(ctx context.Context, input UpdateTenantInput) error
    DeleteTenant(ctx context.Context, id uint) error
    ListTenants(ctx context.Context, onlyActive bool) ([]coreModels.Tenant, error)
    // ActivateTenant(ctx context.Context, id uint) error      
    // DeactivateTenant(ctx context.Context, id uint) error    
}
