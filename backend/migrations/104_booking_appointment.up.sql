-- สร้างตาราง appointments เก็บข้อมูลการจองคิว
CREATE TABLE IF NOT EXISTS appointments (
  id            SERIAL         PRIMARY KEY,               -- รหัสอัตโนมัติ
  branch_id     INT            NOT NULL,                  -- FK ไปยัง branches.id
  service_id    INT            NOT NULL,                  -- FK ไปยัง services.id
  barber_id     INT            NULL,                      -- FK ไปยัง barbers.id (ถ้าเลือกช่างได้)
  customer_id   INT            NOT NULL,                  -- FK ไปยัง users.id (ลูกค้า)
  start_time    TIMESTAMPTZ    NOT NULL,                  -- เวลาเริ่มคิว
  end_time      TIMESTAMPTZ    NOT NULL,                  -- เวลาสิ้นสุดคิว
  status        VARCHAR(20)    NOT NULL DEFAULT 'PENDING',-- สถานะ (PENDING, CONFIRMED, CANCELLED ฯลฯ)
  notes         TEXT           NULL,                      -- ข้อความเพิ่มเติม
  created_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่สร้าง
  updated_at    TIMESTAMPTZ    NOT NULL DEFAULT now(),    -- วันที่อัปเดตล่าสุด
  deleted_at    TIMESTAMPTZ    NULL,                      -- soft‐delete timestamp

  CONSTRAINT fk_appointments_branch   FOREIGN KEY (branch_id)   REFERENCES branches(id)   ON DELETE CASCADE,
  CONSTRAINT fk_appointments_service  FOREIGN KEY (service_id)  REFERENCES services(id)   ON DELETE RESTRICT,
  CONSTRAINT fk_appointments_barber   FOREIGN KEY (barber_id)   REFERENCES barbers(id)    ON DELETE SET NULL,
  CONSTRAINT fk_appointments_customer FOREIGN KEY (customer_id) REFERENCES users(id)      ON DELETE CASCADE
);

-- ดัชนีช่วยค้นหา
CREATE INDEX IF NOT EXISTS idx_appointments_branch   ON appointments(branch_id);
CREATE INDEX IF NOT EXISTS idx_appointments_service  ON appointments(service_id);
CREATE INDEX IF NOT EXISTS idx_appointments_barber   ON appointments(barber_id);
CREATE INDEX IF NOT EXISTS idx_appointments_customer ON appointments(customer_id);
CREATE INDEX IF NOT EXISTS idx_appointments_start    ON appointments(start_time);