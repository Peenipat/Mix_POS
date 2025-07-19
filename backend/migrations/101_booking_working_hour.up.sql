-- สร้างตาราง working_hours เก็บเวลาทำการของแต่ละสาขา
CREATE TABLE IF NOT EXISTS working_hours (
  id          SERIAL         PRIMARY KEY,               -- รหัสอัตโนมัติ
  branch_id   INT            NOT NULL,   
  tenant_id   INT            NOT NULL,               
  weekday     INT            NOT NULL,                  -- 0=Sunday…6=Saturday
  start_time  TIME           NOT NULL,                  -- เวลาเริ่มต้นวันทำการ
  end_time    TIME           NOT NULL,                  -- เวลาสิ้นสุดวันทำการ
  is_closed   BOOLEAN        NOT NULL DEFAULT FALSE,
  created_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่สร้างแถว
  updated_at  TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่แก้ไขล่าสุด
  deleted_at  TIMESTAMPTZ    NULL,                      -- soft‐delete timestamp

  CONSTRAINT uq_working_hours_branch_weekday 
    UNIQUE (branch_id, weekday),
  

  CONSTRAINT fk_working_hours_branch
    FOREIGN KEY (branch_id)
    REFERENCES branches(id)
    ON DELETE CASCADE
);