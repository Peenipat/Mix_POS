package barberBookingController

import (
    helperFunc "myapp/modules/barberbooking"
    barberBookingPort "myapp/modules/barberbooking/port"

    "github.com/gofiber/fiber/v2"
)

type AppointmentStatusLogController struct {
    AppointmentStatusLogService barberBookingPort.IAppointmentStatusLogService
}

func NewAppointmentStatusLogController(
    svc barberBookingPort.IAppointmentStatusLogService,
) *AppointmentStatusLogController {
    return &AppointmentStatusLogController{
        AppointmentStatusLogService: svc,
    }
}

// GetAppointmentLogs godoc
// @Summary      ดึงประวัติสถานะของนัดหมาย
// @Description  คืนรายการ AppointmentStatusLog ของนัดหมายที่ระบุ ภายใต้ Tenant เพื่อความสอดคล้องของ URL
// @Tags         AppointmentLog
// @Accept       json
// @Produce      json
// @Param        tenant_id       path      uint                                           true  "รหัส Tenant"
// @Param        appointment_id  path      uint                                           true  "รหัส Appointment"
// @Success      200             {object}  map[string][]barberBookingModels.AppointmentStatusLog  "คืนค่า status success และ array ของ logs ใน key `data`"
// @Failure      400             {object}  map[string]string                             "Invalid tenant_id หรือ appointment_id"
// @Failure      500             {object}  map[string]string                             "Failed to fetch logs"
// @Router       /tenants/:tenant_id/appointments/:appointment_id/logs [get]
// @Security     ApiKeyAuth
func (ctrl *AppointmentStatusLogController) GetAppointmentLogs(c *fiber.Ctx) error {
    // 1. Parse tenant_id (เพื่อ consistency)
    if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant_id",
        })
    }

    // 2. Parse appointment_id
    apptID, err := helperFunc.ParseUintParam(c, "appointment_id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid appointment_id",
        })
    }

    // 3. Call service
    logs, err := ctrl.AppointmentStatusLogService.GetLogsForAppointment(c.Context(), apptID)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to fetch logs",
            "error":   err.Error(),
        })
    }

    // 4. Return success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   logs,
    })
}
