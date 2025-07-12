CREATE TABLE IF NOT EXISTS customers (
  id         SERIAL PRIMARY KEY,
  tenant_id  INT NOT NULL,
  branch_id   INT NOT NULL, 
  name       TEXT NOT NULL,
  phone      TEXT,
  password    TEXT,
  email      TEXT,
  created_at TIMESTAMP DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ,

  CONSTRAINT uq_customer_email UNIQUE (tenant_id, email)
);

-- Composite Index (tenant_id, email)
CREATE INDEX IF NOT EXISTS idx_customers_tenant_email ON customers(tenant_id, email);
CREATE INDEX idx_customers_lower_email ON customers(LOWER(email));
