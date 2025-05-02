-- สร้างตาราง roles เก็บข้อมูลบทบาทระบบ
CREATE TABLE IF NOT EXISTS roles (
  id          SERIAL        PRIMARY KEY,                -- รหัสบทบาทแบบ Auto-increment
  module_id   INT           NOT NULL,                   -- ตัวเชื่อมกับ modules
  name        VARCHAR(50)   NOT NULL,                   -- ตัวเชื่อมกับ modules
  description TEXT          NULL,                       -- ตัวเชื่อมกับ modules
  created_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),     -- วันที่บันทึก
  updated_at  TIMESTAMPTZ   NOT NULL DEFAULT now(),     -- วันที่อัปเดต
  deleted_at  TIMESTAMPTZ   NULL,                       -- วันที่ลบ (soft delete)
  CONSTRAINT fk_roles_module FOREIGN KEY (module_id) REFERENCES modules(id) ON DELETE CASCADE, -- หากลบ tenant ให้ลบ branches ที่เกี่ยวข้องอัตโนมัติ
  CONSTRAINT uq_roles_module_name UNIQUE (module_id, name)
);