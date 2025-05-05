CREATE TABLE IF NOT EXISTS roles (
    id           SERIAL PRIMARY KEY,
    tenant_id    INT,                              -- nullable เพื่อรองรับ global role
    module_id    INT,                              -- nullable เพื่อรองรับ global หรือ cross-module role
    name         VARCHAR(50) NOT NULL,             -- ชื่อ role (เช่น STAFF, USER)
    description  TEXT,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at   TIMESTAMPTZ,

    -- Composite unique constraint: ห้ามซ้ำใน scope เดียวกัน
    CONSTRAINT uq_roles_scope UNIQUE (tenant_id, module_id, name),

    -- Foreign keys
    CONSTRAINT fk_roles_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE,
    CONSTRAINT fk_roles_module FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE SET NULL
);

-- Optional indexes (เพิ่ม performance)
CREATE INDEX IF NOT EXISTS idx_roles_tenant ON roles(tenant_id);
CREATE INDEX IF NOT EXISTS idx_roles_module ON roles(module_id);
