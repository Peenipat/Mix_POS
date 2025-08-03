CREATE TABLE IF NOT EXISTS working_day_overrides (
  id          SERIAL PRIMARY KEY,
  branch_id   INT NOT NULL,
  work_date   DATE NOT NULL,      
  start_time  TIME NOT NULL,
  end_time    TIME NOT NULL,
  is_closed   BOOLEAN        NOT NULL DEFAULT FALSE,
  reason      TEXT,  
  created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
  deleted_at  TIMESTAMPTZ    NULL
);