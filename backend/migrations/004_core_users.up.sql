-- สร้างตาราง users เก็บข้อมูลบัญชีผู้ใช้ของระบบ
CREATE TABLE IF NOT EXISTS users (
    id          SERIAL PRIMARY KEY,            -- รหัสผู้ใช้ (Auto-increment)
    username    TEXT       NOT NULL,            -- ชื่อแสดงของผู้ใช้ (จะแสดงบน UI)
    email       TEXT       NOT NULL UNIQUE,     -- อีเมลเข้าสู่ระบบ & คีย์ไม่ซ้ำ
    password    TEXT       NOT NULL,            -- รหัสผ่านที่เข้ารหัส (hashed)
    role        VARCHAR(20) NOT NULL,           -- บทบาทของผู้ใช้ (ENUM: SUPER_ADMIN, TENANT_ADMIN, BRANCH_ADMIN, ASSISTANT_MANAGER, STAFF, USER)
    branch_id   INT        NULL,                -- FK: สาขาที่ผู้ใช้สังกัด (null=SuperAdmin ไม่มีสาขา)
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่สร้าง record
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(), -- วันที่ปรับปรุงล่าสุด
    deleted_at  TIMESTAMPTZ NULL,               -- วันที่ลบ (soft delete)
    CONSTRAINT fk_users_branch FOREIGN KEY (branch_id) REFERENCES branches(id) ON DELETE SET NULL -- หากลบ branch ให้ตั้ง branch_id ของผู้ใช้เป็น NULL
);