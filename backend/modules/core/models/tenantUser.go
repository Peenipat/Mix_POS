package coreModels
type TenantUser struct {
	TenantID uint   `gorm:"primaryKey;index" json:"tenant_id"`
	UserID   uint   `gorm:"primaryKey;index" json:"user_id"`
	Tenant   Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	User     User   `gorm:"foreignKey:UserID" json:"user,omitempty"`
}