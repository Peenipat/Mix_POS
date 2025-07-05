package barberBookingService

import (
	"context"
	"time"

	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingPort "myapp/modules/barberbooking/port"

	"gorm.io/gorm"

)

type calendarService struct {
	DB                 *gorm.DB
	workingHourService barberBookingPort.IWorkingHourService
	overrideService    barberBookingPort.IWorkingDayOverrideService
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


func generateSlots(date time.Time, start time.Time, end time.Time) []barberBookingDto.CalendarSlot {
	slots := []barberBookingDto.CalendarSlot{}
	startTime := time.Date(date.Year(), date.Month(), date.Day(), start.Hour(), start.Minute(), 0, 0, time.Local)
	endTime := time.Date(date.Year(), date.Month(), date.Day(), end.Hour(), end.Minute(), 0, 0, time.Local)

	for t := startTime; t.Before(endTime); t = t.Add(30 * time.Minute) {
		slots = append(slots, barberBookingDto.CalendarSlot{
			Start:  t,
			End:    t.Add(30 * time.Minute),
			Status: "open",
		})
	}
	return slots
}


func (s *calendarService) GetAvailableSlots(ctx context.Context, branchID uint, tenantID uint, startDate, endDate time.Time) ([]barberBookingDto.CalendarSlot, error) {
	workingHours, err := s.workingHourService.GetWorkingHours(ctx, branchID, tenantID)
	if err != nil {
		return nil, err
	}

	overrides, err := s.overrideService.GetOverridesByDateRange(ctx, branchID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	overrideMap := map[string]barberBookingModels.WorkingDayOverride{}
	for _, o := range overrides {
		overrideMap[o.WorkDate.Format("2006-01-02")] = o
	}

	var slots []barberBookingDto.CalendarSlot
	for d := startDate; !d.After(endDate); d = d.AddDate(0, 0, 1) {
		dateStr := d.Format("2006-01-02")
		weekday := int(d.Weekday()) // 0 = Sunday

		if override, ok := overrideMap[dateStr]; ok {
			if override.IsClosed {
				continue // skip เพราะปิด
			}
			slots = append(slots, generateSlots(d, override.StartTime.ToTime(d), override.EndTime.ToTime(d))...)

		} else {
			for _, wh := range workingHours {
				if wh.Weekday == weekday && !wh.IsClosed {
					slots = append(slots, generateSlots(d, wh.StartTime, wh.EndTime)...)
				}
			}
		}
	}

	return slots, nil
}
