-- สร้างตาราง barbers แทนช่างผมในแต่ละสาขา
CREATE TABLE IF NOT EXISTS barbers (
  id          SERIAL         PRIMARY KEY,               -- รหัสอัตโนมัติ
  branch_id   INT            NOT NULL,                  -- FK ไปยัง branches.id
  user_id     INT            NOT NULL UNIQUE,           -- FK ไปยัง users.id (ผูก 1:1 กับบัญชีผู้ใช้)
  created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่สร้างแถว
  updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่แก้ไขล่าสุด
  deleted_at  TIMESTAMPTZ    NULL,                      -- soft‐delete timestamp

  CONSTRAINT fk_barbers_branch FOREIGN KEY (branch_id)
    REFERENCES branches(id)
    ON DELETE CASCADE,                              -- หากลบสาขา ให้ลบช่างที่สังกัดด้วย

  CONSTRAINT fk_barbers_user   FOREIGN KEY (user_id)
    REFERENCES users(id)
    ON DELETE CASCADE                               -- หากลบบัญชีผู้ใช้ ให้ลบ record ช่างด้วย
);

-- ดัชนีช่วยค้นหา
CREATE INDEX IF NOT EXISTS idx_barbers_branch ON barbers(branch_id);
CREATE INDEX IF NOT EXISTS idx_barbers_user   ON barbers(user_id);