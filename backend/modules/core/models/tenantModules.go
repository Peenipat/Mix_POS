package coreModels

type TenantModule struct {
	TenantID uint `gorm:"primaryKey" json:"tenant_id"`
	ModuleID uint `gorm:"primaryKey" json:"module_id"`

	// Optional relations
	Tenant *Tenant `gorm:"foreignKey:TenantID" json:"tenant,omitempty"`
	Module *Module `gorm:"foreignKey:ModuleID" json:"module,omitempty"`
}
