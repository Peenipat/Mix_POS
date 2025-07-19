-- สร้างตาราง branches เก็บสาขาของแต่ละ tenant
CREATE TABLE IF NOT EXISTS branches (
    id          SERIAL PRIMARY KEY,                             -- Branch ID
    tenant_id   INT NOT NULL,                                   -- FK to tenants.id
    name        TEXT NOT NULL,                                  -- Branch name
    address     TEXT,                                           -- Optional address
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),             -- Created timestamp
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),             -- Updated timestamp
    deleted_at  TIMESTAMPTZ,                                    -- Soft delete timestamp

    CONSTRAINT fk_branches_tenant
        FOREIGN KEY (tenant_id)
        REFERENCES tenants(id)
        ON DELETE CASCADE                                       -- Cascade delete if tenant is deleted
);

