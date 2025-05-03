-- สร้างตาราง tenants เก็บข้อมูลผู้เช่า (SaaS Tenant)
CREATE TABLE IF NOT EXISTS tenants (
    id          SERIAL PRIMARY KEY,            -- รหัสผู้เช่าแบบ Auto-increment
    name        TEXT       NOT NULL,            -- ชื่อธุรกิจหรือร้านค้า
    domain      TEXT       NOT NULL UNIQUE,     -- โดเมนสำหรับเข้าระบบ (ไม่ซ้ำ)
    is_active   BOOLEAN    NOT NULL DEFAULT TRUE, -- สถานะการใช้งาน (true=active, false=disabled)
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่บันทึกครั้งแรก
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่อัปเดตล่าสุด
    deleted_at  TIMESTAMPTZ NULL                   -- วันที่ลบ record (soft delete)
);


