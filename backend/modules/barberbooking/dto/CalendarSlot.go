// modules/barberbooking/dto/calendar.go

package barberBookingDto

import "time"

type CalendarSlot struct {
	Start  time.Time `json:"start" example:"2025-07-01T09:00:00Z"`
	End    time.Time `json:"end" example:"2025-07-01T09:30:00Z"`
	Status string    `json:"status" example:"open"` // "open" or "closed"
}
