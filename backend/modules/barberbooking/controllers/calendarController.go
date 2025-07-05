package barberBookingController

import (
	"time"

	"github.com/gofiber/fiber/v2"
	barberBookingPortMix "myapp/modules/barberbooking/port"
	barberBookingDto "myapp/modules/barberbooking/dto"
)

type CalendarController struct {
	calendarService barberBookingPortMix.ICalendarService
}

func NewCalendarController(calendarService barberBookingPortMix.ICalendarService) *CalendarController {
	return &CalendarController{calendarService: calendarService}
}


// GetAvailableSlots godoc
// @Summary Get available time slots
// @Description ดึงช่วงเวลาที่สามารถนัดหมายได้ของสาขาในช่วงวันที่กำหนด
// @Tags Calendar
// @Param tenant_id path int true "Tenant ID"
// @Param branch_id path int true "Branch ID"
// @Param start query string true "Start date (format: YYYY-MM-DD)"
// @Param end query string true "End date (format: YYYY-MM-DD)"
// @Produce json
// @Success 200 {array} barberBookingDto.CalendarSlot
// @Failure 400 {object} map[string]string "Invalid input parameters (e.g. missing or invalid tenant_id, branch_id, start or end)"
// @Failure 404 {object} map[string]string "Data not found"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /tenants/{tenant_id}/branches/{branch_id}/available-slots [get]
// @Security     ApiKeyAuth
func (h *CalendarController) GetAvailableSlots(c *fiber.Ctx) error {
	branchID, err := c.ParamsInt("branch_id")
	if err != nil || branchID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid branch ID")
	}

	start := c.Query("start")
	end := c.Query("end")

	if start == "" || end == "" {
		return fiber.NewError(fiber.StatusBadRequest, "start and end date are required")
	}

	startDate, err := time.Parse("2006-01-02", start)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid start date format")
	}
	endDate, err := time.Parse("2006-01-02", end)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, "invalid end date format")
	}

	tenantID, err := c.ParamsInt("tenant_id")
	if err != nil || branchID <= 0 {
		return fiber.NewError(fiber.StatusBadRequest, "invalid branch ID")
	}

	slots, err := h.calendarService.GetAvailableSlots(c.Context(), uint(branchID),uint(tenantID), startDate, endDate)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "failed to get available slots")
	}

	dtoSlots := make([]barberBookingDto.CalendarSlot, 0, len(slots))
for _, s := range slots {
    dtoSlots = append(dtoSlots, barberBookingDto.CalendarSlot{
        Start:  s.Start,
        End:    s.End,
        Status: s.Status,
    })
}
return c.JSON(dtoSlots)
}
