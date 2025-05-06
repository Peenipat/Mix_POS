CREATE TABLE IF NOT EXISTS customers (
  id         SERIAL PRIMARY KEY,
  tenant_id  INT NOT NULL REFERENCES tenants(id),
  name       TEXT NOT NULL,
  phone      TEXT,
  email      TEXT,
  created_at TIMESTAMP DEFAULT now(),

  CONSTRAINT uq_customer_email UNIQUE (tenant_id, email)
);

-- Composite Index (tenant_id, email)
CREATE INDEX IF NOT EXISTS idx_customers_tenant_email ON customers(tenant_id, email);
