CREATE TABLE IF NOT EXISTS appointment_locks (
  id           SERIAL PRIMARY KEY,

  tenant_id    INT NOT NULL,                    -- สำหรับ multi-tenant
  branch_id    INT NOT NULL,                    -- ไม่ต้องมี FK ถ้าอยู่อีกโมดูล
  barber_id    INT NOT NULL,                    -- FK → barbers.id (internal)
  customer_id  INT NOT NULL,                    -- FK → customers.id (internal)

  start_time   TIMESTAMPTZ NOT NULL,
  end_time     TIMESTAMPTZ NOT NULL,
  expires_at   TIMESTAMPTZ NOT NULL,            -- กำหนดเวลาหมดอายุ lock

  is_active    BOOLEAN NOT NULL DEFAULT TRUE,   -- ใช้ระบุว่ายังใช้งานอยู่หรือไม่

  created_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at   TIMESTAMPTZ NOT NULL DEFAULT now(),

  -- 🔗 Constraints
  CONSTRAINT fk_locks_barber
    FOREIGN KEY (barber_id)
    REFERENCES barbers(id)
    ON DELETE CASCADE,

  CONSTRAINT fk_locks_customer
    FOREIGN KEY (customer_id)
    REFERENCES customers(id)
    ON DELETE CASCADE
);


-- ใช้เวลา + barber_id เพื่อเช็ค slot ซ้อนกันเร็วขึ้น
CREATE INDEX IF NOT EXISTS idx_locks_barber_time
  ON appointment_locks(barber_id, start_time, end_time);

-- ใช้ tenant + branch + barber ร่วมกับเวลา
CREATE INDEX IF NOT EXISTS idx_locks_tenant_branch_barber_time
  ON appointment_locks(tenant_id, branch_id, barber_id, start_time, end_time);

-- ใช้สำหรับ query active locks
CREATE INDEX IF NOT EXISTS idx_locks_is_active_expires_at
  ON appointment_locks(is_active, expires_at);