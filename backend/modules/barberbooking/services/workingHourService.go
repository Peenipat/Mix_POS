package barberBookingService

import (
	"context"
	"errors"

	"fmt"
	"gorm.io/gorm"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	"strings"
	"time"
)

type WorkingHourService struct {
	DB *gorm.DB
}

func NewWorkingHourService(db *gorm.DB) barberBookingPort.IWorkingHourService {
	return &WorkingHourService{DB: db}
}

func (s *WorkingHourService) GetWorkingHours(ctx context.Context, branchID uint, tenantID uint) ([]barberBookingModels.WorkingHour, error) {
	var hours []barberBookingModels.WorkingHour
	err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND tenant_id = ? AND deleted_at IS NULL", branchID, tenantID).
		Order("weekday asc").
		Find(&hours).Error
	if err != nil {
		return nil, err
	}
	return hours, nil
}

func (s *WorkingHourService) UpdateWorkingHours(ctx context.Context, branchID uint, tenantID uint, input []barberBookingDto.WorkingHourInput) error {
	tx := s.DB.WithContext(ctx).Begin()

	for _, wh := range input {
		if wh.Weekday < 0 || wh.Weekday > 6 {
			tx.Rollback()
			return fmt.Errorf("invalid weekday: %d", wh.Weekday)
		}

		fmt.Println( "weekday=", wh.Weekday) 
		var existing barberBookingModels.WorkingHour
		err := tx.
			Where("branch_id = ? AND tenant_id = ? AND weekday = ?", branchID, tenantID, wh.Weekday).
			First(&existing).Error

		if errors.Is(err, gorm.ErrRecordNotFound) {
			newWH := barberBookingModels.WorkingHour{
				BranchID:  branchID,
				TenantID:  tenantID,
				Weekday:   wh.Weekday,
				StartTime: wh.StartTime,
				EndTime:   wh.EndTime,
				IsClosed:  wh.IsClosed,
			}
			if err := tx.Create(&newWH).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else if err == nil {
			// อัปเดต
			existing.StartTime = wh.StartTime
			existing.EndTime = wh.EndTime
			existing.IsClosed = wh.IsClosed
			if err := tx.Save(&existing).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}

func (s *WorkingHourService) CreateWorkingHours(ctx context.Context, branchID uint, input barberBookingDto.WorkingHourInput) error {
	if input.Weekday < 0 || input.Weekday > 6 {
		return fmt.Errorf("invalid weekday: %d", input.Weekday)
	}
	if input.StartTime.After(input.EndTime) || input.StartTime.Equal(input.EndTime) {
		return fmt.Errorf("start time must be before end time")
	}

	wh := barberBookingModels.WorkingHour{
		BranchID:  branchID,
		Weekday:   input.Weekday,
		StartTime: input.StartTime,
		EndTime:   input.EndTime,
	}

	if err := s.DB.WithContext(ctx).Create(&wh).Error; err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			return fmt.Errorf("working hour for weekday %d already exists", input.Weekday)
		}
		return err
	}

	return nil
}

func (s *WorkingHourService) GetAvailableSlots(
	ctx context.Context,
	branchID uint,
	tenantID uint,
	filter string, // "week", "month", หรือ ""
	fromTime *string, // nullable
	toTime *string,
) (map[string][]string, error) {
	loc := time.Now().Location()
	now := time.Now().In(loc)

	var dates []time.Time

	switch filter {
	case "week":
		isoWeekday := getWeekday(now)
		monday := now.AddDate(0, 0, -(isoWeekday - 1)).Truncate(24 * time.Hour)
		for i := 0; i < 7; i++ {
			d := monday.AddDate(0, 0, i)
			if d.Before(now.Truncate(24 * time.Hour)) {
				continue
			}
			dates = append(dates, d)
		}
	case "month":
		firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, loc)
		for d := firstDay; d.Month() == now.Month(); d = d.AddDate(0, 0, 1) {
			if d.Before(now.Truncate(24 * time.Hour)) {
				continue
			}
			dates = append(dates, d)
		}
	default:
		dates = []time.Time{now.Truncate(24 * time.Hour)}
	}

	defaultHours, err := s.GetWorkingHours(ctx, branchID, tenantID)
	if err != nil {
		return nil, err
	}

	startStr := dates[0].Format("2006-01-02")
	endStr := dates[len(dates)-1].Format("2006-01-02")
	var overrides []barberBookingModels.WorkingDayOverride
	if err := s.DB.WithContext(ctx).
		Where("branch_id = ? AND work_date BETWEEN ? AND ? AND deleted_at IS NULL", branchID, startStr, endStr).
		Find(&overrides).Error; err != nil {
		return nil, err
	}
	overrideMap := make(map[string]barberBookingModels.WorkingDayOverride)
	for _, o := range overrides {
		overrideMap[o.WorkDate.Format("2006-01-02")] = o
	}

	result := make(map[string][]string)
	for _, date := range dates {
		dayStr := date.Format("2006-01-02")
		var start, end string

		if o, ok := overrideMap[dayStr]; ok {
			if o.IsClosed {
				result[dayStr] = []string{} 
				continue
			}
			start = o.StartTime.Format("15:04")
			end = o.EndTime.Format("15:04")
		} else {
			weekday := getWeekday(date)
			found := false
			for _, d := range defaultHours {
				if d.Weekday == weekday {
					if d.IsClosed {
						result[dayStr] = []string{} 
						found = true
						break
					}
					start = d.StartTime.Format("15:04")
					end = d.EndTime.Format("15:04")
					found = true
					break
				}
			}
			if !found {
				result[dayStr] = []string{} 
				continue
			}
		}
		

		realStart := start
		realEnd := end
		if fromTime != nil {
			realStart = maxTime(start, *fromTime)
		}
		if toTime != nil {
			realEnd = minTime(end, *toTime)
		}
		if realStart >= realEnd {
			continue
		}

		result[dayStr] = generateSlots(realStart, realEnd, 30)
	}

	return result, nil
}

func generateSlots(start string, end string, interval int) []string {
	slots := []string{}
	layout := "15:04"

	startTime, _ := time.Parse(layout, start)
	endTime, _ := time.Parse(layout, end)

	for t := startTime; !t.After(endTime); t = t.Add(time.Minute * time.Duration(interval)) {
		s := t.Format("15:04")
		slots = append(slots, s)
	}
	return slots
}



func getWeekday(t time.Time) int {
	return int(t.Weekday())
}
func maxTime(a, b string) string {
	if a > b {
		return a
	}
	return b
}
func minTime(a, b string) string {
	if a < b {
		return a
	}
	return b
}
