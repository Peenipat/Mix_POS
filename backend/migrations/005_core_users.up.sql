--สร้างตาราง users ใหม่ (เก็บ role_id แทน role:string)
CREATE TABLE users (
  id          SERIAL         PRIMARY KEY,               -- PK อัตโนมัติ
  username    TEXT           NOT NULL,                  -- ชื่อแสดงบน UI
  email       TEXT           NOT NULL UNIQUE,           -- อีเมลเข้าสู่ระบบ (unique)
  password    TEXT           NOT NULL,                  -- รหัสผ่าน bcrypt hash
  role_id     INT            NOT NULL,                  -- FK ไปยัง roles.id
  branch_id   INT            NULL,                      -- FK ไปยัง branches.id
  img_path    TEXT           NULL,
  img_name    TEXT           NULL,
  created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่สร้าง
  updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่อัปเดตล่าสุด
  deleted_at  TIMESTAMPTZ    NULL,                      -- soft‐delete

  CONSTRAINT fk_users_role   FOREIGN KEY (role_id)   REFERENCES roles(id)    ON DELETE RESTRICT,
  CONSTRAINT fk_users_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE SET NULL
);

-- 3) สร้างดัชนีช่วยค้นหา
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
CREATE INDEX idx_users_branch     ON users(branch_id);


