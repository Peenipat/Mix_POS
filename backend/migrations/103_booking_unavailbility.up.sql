-- สร้างตาราง unavailabilities เก็บวันหยุดพักของช่างหรือสาขา
CREATE TABLE IF NOT EXISTS unavailabilities (
  id          SERIAL      PRIMARY KEY,        -- รหัสอัตโนมัติ
  barber_id   INT         NULL,               -- FK ไปยัง barbers.id (ถ้ามี)
  branch_id   INT         NULL,               -- FK ไปยัง branches.id (ถ้ามี)
  date        DATE        NOT NULL,           -- วันที่หยุดให้บริการ
  reason      TEXT        NULL,               -- เหตุผลหรือคำอธิบาย
  CONSTRAINT fk_unavailabilities_barber FOREIGN KEY (barber_id)
    REFERENCES barbers(id) ON DELETE SET NULL,
  CONSTRAINT fk_unavailabilities_branch FOREIGN KEY (branch_id)
    REFERENCES branches(id) ON DELETE SET NULL
);

-- ดัชนีช่วยค้นหา
CREATE INDEX IF NOT EXISTS idx_unavailabilities_date   ON unavailabilities(date);
CREATE INDEX IF NOT EXISTS idx_unavailabilities_barber ON unavailabilities(barber_id);
CREATE INDEX IF NOT EXISTS idx_unavailabilities_branch ON unavailabilities(branch_id);