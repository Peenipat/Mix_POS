-- 007_core_modules.up.sql
-- สร้างตารางเก็บชื่อ feature/โมดูลต่างๆ ในระบบ
CREATE TABLE IF NOT EXISTS modules (
  id            SERIAL        PRIMARY KEY,             -- รหัสโมดูล
  key           VARCHAR(50)   NOT NULL UNIQUE,         -- ชื่อโมดูล (เช่น "BOOKING", "POS_RESTAURANT", "INVENTORY")\
  description   TEXT          NOT NULL,
  created_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),  -- เมื่อสร้าง
  updated_at    TIMESTAMPTZ   NOT NULL DEFAULT now(),  -- เมื่ออัปเดตล่าสุด
  deleted_at    TIMESTAMPTZ   NULL                     -- soft-delete
);

-- ช่วยค้นหาโมดูลตามชื่อ
CREATE INDEX IF NOT EXISTS idx_modules_name ON modules(key);
