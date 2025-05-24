-- สร้างตารางเก็บ Log แต่ละรายการ
CREATE TABLE IF NOT EXISTS system_logs (
  log_id            BIGSERIAL      PRIMARY KEY,               -- PK อัตโนมัติ
  created_at        TIMESTAMPTZ     NOT NULL DEFAULT now(),   -- เวลาที่บันทึก
  user_id           INT             NULL,                     -- FK ไป users.id (ถ้ามี)
  user_role         VARCHAR(20)     NULL,                     -- role ของผู้ใช้ตอนนั้น
  action            VARCHAR(50)     NOT NULL,                 -- เช่น "LOGIN_SUCCESS"
  resource          VARCHAR(50)     NOT NULL,                 -- เช่น "Auth" หรือ "Order"
  status            VARCHAR(20)     NOT NULL,                 -- "success" / "failure"
  ip_address        INET            NULL,                     -- IP ต้นทาง
  http_method       VARCHAR(10)     NOT NULL DEFAULT 'GET',   -- GET/POST/PUT...
  endpoint          TEXT            NOT NULL,                 -- path ของ API
  status_code       INT             NULL,                     -- HTTP status code ตอบกลับ
  x_forwarded_for   VARCHAR(100)    NULL,                     -- Header X-Forwarded-For
  user_agent        TEXT            NULL,                     -- Header User-Agent
  referer           TEXT            NULL,                     -- Header Referer
  origin            TEXT            NULL,                     -- Header Origin
  client_app        VARCHAR(50)     NULL,                     -- เช่น mobile/desktop
  branch_id         INT             NULL,                     -- FK ไป branches.id
  details           JSONB           NULL,                     -- payload เสริม, e.g. {"email":"..."}
  metadata          JSONB           NULL                      -- คีย์–ค่าขยายอื่นๆ
);

-- สร้างดัชนีช่วยค้นหา
CREATE INDEX IF NOT EXISTS idx_system_logs__user    ON system_logs(user_id, created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_logs__action  ON system_logs(action,    created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_logs__date    ON system_logs(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_system_logs__gin     ON system_logs USING GIN(details);
