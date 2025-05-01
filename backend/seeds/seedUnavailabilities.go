package seeds

import (
    "time"

    "gorm.io/gorm"
    coreModels "myapp/modules/core/models"
    bookingModels "myapp/modules/barberbooking/models"
)

// SeedUnavailabilities สร้างตัวอย่างวันหยุด/ไม่ว่างสำหรับสาขาและช่าง
// ต้องรันหลังจาก SeedBranches() และ SeedBarbers()
func SeedUnavailabilities(db *gorm.DB) error {
    // 1) โหลด Default Branch
    var branch coreModels.Branch
    if err := db.Where("name = ?", "Default Branch").First(&branch).Error; err != nil {
        return err
    }

    // 2) โหลดช่างทั้งหมดใน Default Branch
    var barbers []bookingModels.Barber
    if err := db.Where("branch_id = ?", branch.ID).Find(&barbers).Error; err != nil {
        return err
    }

    // 3) กำหนดวันที่ตัวอย่าง (วันพรุ่งนี้)
    today := time.Now()
    unavailDate := time.Date(
        today.Year(), today.Month(), today.Day()+1,
        0, 0, 0, 0, today.Location(),
    )

    // 4) สร้างวันหยุดระดับสาขา (เช่น ปิดสาขาเพื่อ maintenance)
    branchRecord := bookingModels.Unavailability{
        BranchID: &branch.ID,
        BarberID: nil,
        Date:     unavailDate,
    }
    branchAttrs := bookingModels.Unavailability{
        Reason: "Branch maintenance",
    }
    if err := db.Where(branchRecord).
        Assign(branchAttrs).
        FirstOrCreate(&branchRecord).Error; err != nil {
        return err
    }

    // 5) สร้างวันไม่ว่างสำหรับแต่ละช่าง (เช่น ฝึกอบรมช่าง)
    for _, b := range barbers {
        barberRecord := bookingModels.Unavailability{
            BranchID: &branch.ID,
            BarberID: &b.ID,
            Date:     unavailDate,
        }
        barberAttrs := bookingModels.Unavailability{
            Reason: "Barber training",
        }
        if err := db.Where(barberRecord).
            Assign(barberAttrs).
            FirstOrCreate(&barberRecord).Error; err != nil {
            return err
        }
    }

    return nil
}
