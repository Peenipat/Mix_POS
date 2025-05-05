-- Join table mapping tenants to user accounts (many-to-many)
CREATE TABLE IF NOT EXISTS tenant_users (
    tenant_id INT NOT NULL,               -- FK to tenants.id
    user_id   INT NOT NULL,               -- FK to users.id
    PRIMARY KEY (tenant_id, user_id),     -- Composite PK
    CONSTRAINT fk_tenant_users_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_tenant_users_user   FOREIGN KEY (user_id)   REFERENCES users(id)   ON DELETE CASCADE
);