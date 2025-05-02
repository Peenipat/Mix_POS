-- สร้างตาราง branches เก็บสาขาของแต่ละ tenant
CREATE TABLE IF NOT EXISTS branches (
    id          SERIAL PRIMARY KEY,            -- รหัสสาขาแบบ Auto-increment
    tenant_id   INT        NOT NULL,            -- FK ไปยัง tenants.id (เจ้าของสาขา)
    name        TEXT       NOT NULL,            -- ชื่อสาขา
    address     TEXT       NULL,                -- ที่อยู่สาขา
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่บันทึก
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่อัปเดต
    deleted_at  TIMESTAMPTZ NULL,                -- วันที่ลบ (soft delete)
    CONSTRAINT fk_branches_tenant FOREIGN KEY (tenant_id) REFERENCES tenants(id) ON DELETE CASCADE -- หากลบ tenant ให้ลบ branches ที่เกี่ยวข้องอัตโนมัติ
);