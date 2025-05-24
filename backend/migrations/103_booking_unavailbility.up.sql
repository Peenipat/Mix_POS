CREATE TABLE IF NOT EXISTS unavailabilities (
  id          SERIAL PRIMARY KEY,
  barber_id   INT NULL,
  branch_id   INT NULL,
  date        DATE NOT NULL,
  reason      TEXT,
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at  TIMESTAMPTZ                  -- soft‐delete timestamp

  -- ❌ ตัด FK ออกเพื่อไม่ผูกตายตัวกับโมดูลอื่น
  -- CONSTRAINT fk_unavailabilities_barber FOREIGN KEY (barber_id) REFERENCES barbers(id) ON DELETE SET NULL,
  -- CONSTRAINT fk_unavailabilities_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_unavailabilities_date   ON unavailabilities(date);
CREATE INDEX IF NOT EXISTS idx_unavailabilities_barber ON unavailabilities(barber_id);
CREATE INDEX IF NOT EXISTS idx_unavailabilities_branch ON unavailabilities(branch_id);
CREATE UNIQUE INDEX uq_unavailability ON unavailabilities (date, barber_id, branch_id);
