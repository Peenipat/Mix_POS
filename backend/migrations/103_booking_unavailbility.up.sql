CREATE TABLE IF NOT EXISTS unavailabilities (
  id          SERIAL PRIMARY KEY,
  barber_id   INT NULL,               -- ref: booking.barbers.id
  branch_id   INT NULL,               -- ref: core.branches.id
  date        DATE NOT NULL,          -- วันหยุด
  reason      TEXT,

  -- ❌ ตัด FK ออกเพื่อไม่ผูกตายตัวกับโมดูลอื่น
  -- CONSTRAINT fk_unavailabilities_barber FOREIGN KEY (barber_id) REFERENCES barbers(id) ON DELETE SET NULL,
  -- CONSTRAINT fk_unavailabilities_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_unavailabilities_date   ON unavailabilities(date);
CREATE INDEX IF NOT EXISTS idx_unavailabilities_barber ON unavailabilities(barber_id);
CREATE INDEX IF NOT EXISTS idx_unavailabilities_branch ON unavailabilities(branch_id);
