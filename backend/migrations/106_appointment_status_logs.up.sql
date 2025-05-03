CREATE TABLE IF NOT EXISTS appointment_status_logs (
  id                  SERIAL PRIMARY KEY,
  appointment_id      INT NOT NULL,  -- FK ภายในโมดูล booking
  old_status          VARCHAR(20),
  new_status          VARCHAR(20),
  changed_by_user_id  INT,           -- อ้างถึง core.users → ไม่ใส่ FK
  changed_by_customer_id INT,        --  FK → customers.id
  changed_at          TIMESTAMPTZ DEFAULT now(),
  notes               TEXT,

  CONSTRAINT fk_statuslog_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
  CONSTRAINT fk_statuslog_customer FOREIGN KEY (changed_by_customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_log_appointment_id ON appointment_status_logs(appointment_id);
CREATE INDEX IF NOT EXISTS idx_log_changed_by_user ON appointment_status_logs(changed_by_user_id);
CREATE INDEX IF NOT EXISTS idx_log_changed_by_customer ON appointment_status_logs(changed_by_customer_id);
