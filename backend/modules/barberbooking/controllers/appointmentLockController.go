package barberBookingController

import (
	"context"
	"log"
	"time"

	barberBookingPort "myapp/modules/barberbooking/port"
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type AppointmentLockController struct {
	Service barberBookingPort.IAppointmentLock
}

func NewAppointmentLockController(service barberBookingPort.IAppointmentLock) *AppointmentLockController {
	return &AppointmentLockController{
		Service: service,
	}
}

type CreateAppointmentLockRequest struct {
	TenantID   uint      `json:"tenant_id" validate:"required"`
	BranchID   uint      `json:"branch_id" validate:"required"`
	BarberID   uint      `json:"barber_id" validate:"required"`
	CustomerID uint      `json:"customer_id" validate:"required"`
	StartTime  time.Time `json:"start_time" validate:"required"`
	EndTime    time.Time `json:"end_time" validate:"required"`
}

// POST /appointment-locks
func (ctl *AppointmentLockController) CreateAppointmentLock(c *fiber.Ctx) error {
	var req CreateAppointmentLockRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	log.Println("parsed:", req)

	lock, err := ctl.Service.CreateAppointmentLock(
		context.Background(),
		barberBookingPort.AppointmentLockInput{
			TenantID:   req.TenantID,
			BranchID:   req.BranchID,
			BarberID:   req.BarberID,
			CustomerID: req.CustomerID,
			StartTime:  req.StartTime,
			EndTime:    req.EndTime,
		},
	)
	if err != nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(lock)
}

func (ctl *AppointmentLockController) ReleaseAppointmentLock(c *fiber.Ctx) error {
	idParam := c.Params("lock_id")
	lockID, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid lock ID"})
	}

	err = ctl.Service.ReleaseAppointmentLock(context.Background(), uint(lockID))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// GET /appointment-locks?branch_id=1&barber_id=2&date=2025-07-13
func (ctl *AppointmentLockController) GetAppointmentLocks(c *fiber.Ctx) error {
	branchID, _ := strconv.Atoi(c.Query("branch_id"))
	barberID, _ := strconv.Atoi(c.Query("barber_id"))
	dateStr := c.Query("date")

	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "invalid date format"})
	}

	locks, err := ctl.Service.GetAppointmentLocks(
		context.Background(),
		uint(branchID),
		uint(barberID),
		date,
	)
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(locks)
}

