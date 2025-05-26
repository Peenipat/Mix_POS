package corePort

import "context"

// "context"
// coreModels "myapp/modules/core/models"

type RegisterInput struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

type CreateUserInput struct {
	Username string `json:"username" validate:"required"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Role     string `json:"role" validate:"required"`
}

type ChangeRoleInput struct {
	ID   uint   `json:"id" validate:"required"`
	Role string `json:"role" validate:"required"`
}

type UserInfoResponse struct {
	ID       uint   `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	RoleID   uint   `json:"role_id"`
	Role     string `json:"role"`
}

type MeDTO struct {
    ID        uint   `json:"id"`
    Username  string `json:"username"`
    Email     string `json:"email"`
    BranchID  *uint  `json:"branch_id"`
    TenantIDs []uint `json:"tenant_ids"`
}

type LoginResponse struct {
	Token string           `json:"token"`
	User  UserInfoResponse `json:"user"`
}

type IUser interface {
	CreateUserFromRegister(input RegisterInput) error
	CreateUserFromAdmin(input CreateUserInput) error
	ChangeRoleFromAdmin(input ChangeRoleInput) error
	GetAllUsers(limit int, offset int) ([]UserInfoResponse, error)
	FilterUsersByRole(role string) ([]UserInfoResponse, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error 
	Me(ctx context.Context, userID uint) (*MeDTO, error)

	// GetUserByID(ctx context.Context, id uint) (*coreModels.User, error)
	// GetUserByEmail(ctx context.Context, email string) (*coreModels.User, error)

	// UpdateUser(ctx context.Context, u *coreModels.User) error
	// ChangePassword(ctx context.Context, userID uint, newPassword string) error
	// DeleteUser(ctx context.Context, id uint) error

	//Authenticate(ctx context.Context, email, password string) (*coreModels.User, error)
	// ListUsers(ctx context.Context, filter UserFilter) ([]coreModels.User, error)
}


