CREATE TABLE IF NOT EXISTS barber_workload (
  id SERIAL PRIMARY KEY,
  barber_id INT NOT NULL,
  date DATE NOT NULL,
  total_appointments INT NOT NULL DEFAULT 0,
  total_hours INT NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ DEFAULT now(),

  CONSTRAINT fk_barber_workload_barber FOREIGN KEY (barber_id)
    REFERENCES barbers(id)
    ON DELETE CASCADE
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_barber_workload_unique
  ON barber_workload(barber_id, date);
