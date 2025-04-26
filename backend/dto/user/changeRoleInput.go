package userDto

type ChangeRoleInput struct {
	ID uint   `json:"id" validate:"required"`
	Role   string `json:"role" validate:"required"`
}