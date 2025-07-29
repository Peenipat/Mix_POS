package barberBookingService

import (
	"context"
	"time"

	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingPort "myapp/modules/barberbooking/port"

	"gorm.io/gorm"
)

type calendarService struct {
	DB                 *gorm.DB
	workingHourService barberBookingPort.IWorkingHourService
	overrideService    barberBookingPort.IWorkingDayOverrideService
}

// GetAvailableSlots implements barberBookingPort.ICalendarService.
func (c *calendarService) GetAvailableSlots(ctx context.Context, branchID uint, tenantID uint, startDate time.Time, endDate time.Time) ([]barberBookingDto.CalendarSlot, error) {
	panic("unimplemented")
}

func NewCalendarService(
	db *gorm.DB,
	workingHourSvc barberBookingPort.IWorkingHourService,
	overrideSvc barberBookingPort.IWorkingDayOverrideService,
) barberBookingPort.ICalendarService {
	return &calendarService{
		DB:                 db,
		workingHourService: workingHourSvc,
		overrideService:    overrideSvc,
	}
}
