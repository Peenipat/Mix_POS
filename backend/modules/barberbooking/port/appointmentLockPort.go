package barberBookingPort

import (
	"context"
	"time"
	barberBookingModels "myapp/modules/barberbooking/models"
)

type AppointmentLockInput struct {
	TenantID    uint      `json:"tenant_id"`
	BranchID    uint      `json:"branch_id"`
	BarberID    uint      `json:"barber_id"`
	CustomerID  uint      `json:"customer_id"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
}

type IAppointmentLock interface {
	// 🔒 สร้าง lock ชั่วคราว
	CreateAppointmentLock(ctx context.Context, input AppointmentLockInput) (*barberBookingModels.AppointmentLock, error)

	// 🧹 ปล่อย lock (เช่น ปิด modal)
	ReleaseAppointmentLock(ctx context.Context, lockID uint) error

	// ✅ เช็คว่า slot ยังว่างหรือไม่ (รวมทั้ง confirmed และ locked)
	IsSlotAvailable(ctx context.Context, tenantID, branchID, barberID uint, start, end time.Time) (bool, error)

	// 📦 ดึง lock ทั้งหมดของวัน
	GetAppointmentLocks(ctx context.Context, branchID, barberID uint, date time.Time) ([]barberBookingModels.AppointmentLock, error)

	// ⏰ ล้าง lock ที่หมดอายุ (ใช้ใน Cron job)
	CleanupExpiredLocks(ctx context.Context) error
}
