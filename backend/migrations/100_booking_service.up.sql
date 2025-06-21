DROP TABLE IF EXISTS services;

CREATE TABLE services (
  id         SERIAL PRIMARY KEY,
  tenant_id  INT NOT NULL,
  branch_id  INT NOT NULL,
  name       VARCHAR(100) NOT NULL,
  description VARCHAR(100)   NOT NULL,
  duration   INT NOT NULL,
  price      NUMERIC NOT NULL,
  img_path    TEXT           NULL,
  img_name    TEXT           NULL,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at TIMESTAMPTZ
);

CREATE INDEX IF NOT EXISTS idx_services_tenant_not_deleted
  ON services (tenant_id, deleted_at)
  WHERE deleted_at IS NULL;
