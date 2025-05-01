-- สร้างตารางเก็บข้อมูล Service แต่ละรายการ
CREATE TABLE IF NOT EXISTS services (
  id          SERIAL         PRIMARY KEY,              -- รหัสอัตโนมัติ
  name        VARCHAR(100)   NOT NULL,                 -- ชื่อบริการ
  duration    INT            NOT NULL,                 -- ระยะเวลาโดยประมาณ
  price       NUMERIC        NOT NULL,                 -- ราคาบริการ
  created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),   -- เวลาสร้างแถว
  updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),   -- เวลาแก้ไขล่าสุด
  deleted_at  TIMESTAMPTZ    NULL                      -- soft‐delete timestamp
);

-- optional: index ช่วยค้นหาเฉพาะรายการที่ยังไม่ลบ
CREATE INDEX IF NOT EXISTS idx_services_not_deleted
  ON services (deleted_at)
  WHERE deleted_at IS NULL;