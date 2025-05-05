CREATE TABLE IF NOT EXISTS modules (
  id           SERIAL PRIMARY KEY,
  name         VARCHAR(100) NOT NULL UNIQUE,     -- ชื่อโมดูล เช่น "barber_booking"
  description  TEXT,                             -- คำอธิบายเพิ่มเติม
  created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
