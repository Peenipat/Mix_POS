package Core_authDto
type UserInfoResponse struct {
    ID       uint   `json:"id"`
    Username string `json:"username"`
    Email    string `json:"email"`
    RoleID   uint   `json:"role_id"` // เพิ่มตรงนี้
    Role     string `json:"role"`    // และปรับชื่อ field ให้ตรงกับ JSON
}

type LoginResponse struct {
	Token string           `json:"token"`
	User  UserInfoResponse `json:"user"`
}
