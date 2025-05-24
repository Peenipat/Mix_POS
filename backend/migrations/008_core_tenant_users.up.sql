CREATE TABLE IF NOT EXISTS tenant_users (
  tenant_id INT NOT NULL,
  user_id   INT NOT NULL,

  CONSTRAINT pk_tenant_users PRIMARY KEY (tenant_id, user_id),
  CONSTRAINT fk_tenant_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
  CONSTRAINT fk_tenant_users_user   FOREIGN KEY (user_id)   REFERENCES users(id)   ON DELETE CASCADE
);
