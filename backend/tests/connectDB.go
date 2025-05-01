package tests

import (
    "myapp/database"
    coreModels "myapp/models/core"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

// SetupTestDB เปิด DB in-memory แล้ว AutoMigrate โมเดลทั้งหมดที่ test ต้องใช้
func SetupTestDB() *gorm.DB {
    db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
    if err != nil {
        panic("failed to open test db: " + err.Error())
    }
    // override global
    database.DB = db

    // สร้างตารางก่อนรัน test
    if err := db.AutoMigrate(
        &coreModels.Role{},
        &coreModels.User{},
        // ถ้ามี Branch ที่ service test ต้องใช้ ก็เพิ่ม &coreModels.Branch{},
        // ถ้า test ระบบ Booking ก็เพิ่มโมเดล Booking ด้วย:
        // &bookingModels.Service{}, &bookingModels.Barber{}, &bookingModels.Appointment{},
    ); err != nil {
        panic("failed to migrate test db: " + err.Error())
    }

    return db
}
