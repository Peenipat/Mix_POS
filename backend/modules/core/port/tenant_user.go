package corePort
import (
	"context"
	coreModels "myapp/modules/core/models"
)
type ITenantUser interface {
    AddUserToTenant(ctx context.Context, tenantID, userID uint) error
    RemoveUserFromTenant(ctx context.Context, tenantID, userID uint) error
    ListUsersByTenant(ctx context.Context, tenantID uint) ([]coreModels.User, error)
    ListTenantsByUser(ctx context.Context, userID uint) ([]coreModels.Tenant, error)
    IsUserInTenant(ctx context.Context, tenantID, userID uint) (bool, error)
}