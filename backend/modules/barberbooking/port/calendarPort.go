package barberBookingPort

import (
	"time"
	"context"
	barberBookingDto "myapp/modules/barberbooking/dto"
)


type ICalendarService interface {
	GetAvailableSlots(
		ctx context.Context,
		branchID uint,
		tenantID uint,
		startDate time.Time,
		endDate time.Time,
	) ([]barberBookingDto.CalendarSlot, error)
}