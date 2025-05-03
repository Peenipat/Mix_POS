CREATE TABLE IF NOT EXISTS appointment_reviews (
  id             SERIAL PRIMARY KEY,
  appointment_id INT NOT NULL UNIQUE, -- 1 รีวิวต่อ 1 appointment
  customer_id    INT,                 -- FK → customers.id
  rating         INT NOT NULL CHECK (rating >= 1 AND rating <= 5),
  comment        TEXT,
  created_at     TIMESTAMPTZ DEFAULT now(),

  CONSTRAINT fk_review_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
  CONSTRAINT fk_review_customer    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_reviews_customer_id ON appointment_reviews(customer_id);
