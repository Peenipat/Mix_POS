ALTER TABLE working_hours
  -- เอาเฉพาะส่วนเวลา (TIME) บวกกับวันนี้ ให้กลายเป็น TIMESTAMPTZ
  ALTER COLUMN start_time TYPE TIMESTAMP WITH TIME ZONE
    USING (date_trunc('day', now()) + start_time::interval),
  ALTER COLUMN end_time TYPE TIMESTAMP WITH TIME ZONE
    USING (date_trunc('day', now()) + end_time::interval);
