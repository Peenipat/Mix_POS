CREATE TABLE IF NOT EXISTS tenant_modules (
  tenant_id INT NOT NULL,
  module_id INT NOT NULL,

  CONSTRAINT fk_tenant_modules_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
  CONSTRAINT fk_tenant_modules_module FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE,

  CONSTRAINT pk_tenant_modules PRIMARY KEY (tenant_id, module_id)
);