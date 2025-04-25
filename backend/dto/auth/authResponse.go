package authDto
type UserInfoResponse struct {
	ID       uint   `json:"id" example:"1"`
	Username string `json:"username" example:"admin"`
	Email    string `json:"email" example:"admin@example.com"`
	Role     string `json:"role" example:"SUPER_ADMIN"`
}

type LoginResponse struct {
	Token string           `json:"token"`
	User  UserInfoResponse `json:"user"`
}
