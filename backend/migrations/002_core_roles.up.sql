-- สร้างตาราง roles เก็บข้อมูลบทบาทระบบ
CREATE TABLE IF NOT EXISTS roles (
    id          SERIAL PRIMARY KEY,            -- รหัสบทบาทแบบ Auto-increment
    name        VARCHAR(50) NOT NULL UNIQUE,    -- ชื่อบทบาท (เช่น SUPER_ADMIN, BRANCH_ADMIN)
    description TEXT       NULL,                -- คำอธิบายบทบาท
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่บันทึก
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่อัปเดต
    deleted_at  TIMESTAMPTZ NULL                   -- วันที่ลบ (soft delete)
);