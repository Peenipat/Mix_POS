CREATE TABLE IF NOT EXISTS appointments (
  id            SERIAL PRIMARY KEY,

  tenant_id     INT NOT NULL,                    -- ✅ เพิ่ม tenant สำหรับ multi-tenant
  branch_id     INT NOT NULL,                    -- ❌ ไม่มี FK เพราะอยู่อีก module
  service_id    INT NOT NULL,                    -- ✅ FK → services.id (internal)
  barber_id     INT,                             -- ✅ FK → barbers.id (optional)
  customer_id   INT NOT NULL,                    -- ✅ FK → customers.id (internal)
  user_id       INT,                             -- ❌ ไม่มี FK เพราะ user อยู่อีก module

  start_time    TIMESTAMPTZ NOT NULL,
  end_time      TIMESTAMPTZ NOT NULL,

  status        VARCHAR(20) NOT NULL DEFAULT 'PENDING',
  notes         TEXT,

  created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at    TIMESTAMPTZ,

  -- 🔗 Constraints
  CONSTRAINT fk_appointments_service
    FOREIGN KEY (service_id)
    REFERENCES services(id)
    ON DELETE RESTRICT,

  CONSTRAINT fk_appointments_barber
    FOREIGN KEY (barber_id)
    REFERENCES barbers(id)
    ON DELETE SET NULL,

  CONSTRAINT fk_appointments_customer
    FOREIGN KEY (customer_id)
    REFERENCES customers(id)
    ON DELETE CASCADE
);

-- 🔍 Composite index ช่วยให้ query ดีขึ้น
CREATE INDEX IF NOT EXISTS idx_appointments_tenant_branch_time
  ON appointments(tenant_id, branch_id, start_time, end_time);

CREATE INDEX IF NOT EXISTS idx_appointments_branch   ON appointments(branch_id);
CREATE INDEX IF NOT EXISTS idx_appointments_service  ON appointments(service_id);
CREATE INDEX IF NOT EXISTS idx_appointments_barber   ON appointments(barber_id);
CREATE INDEX IF NOT EXISTS idx_appointments_customer ON appointments(customer_id);
CREATE INDEX IF NOT EXISTS idx_appointments_start    ON appointments(start_time);
