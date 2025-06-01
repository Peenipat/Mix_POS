package barberBookingService

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type appointmentService struct {
	DB         *gorm.DB
	LogService barberBookingPort.IAppointmentStatusLogService
}

func NewAppointmentService(
	db *gorm.DB,
	logSvc barberBookingPort.IAppointmentStatusLogService,
) *appointmentService {
	return &appointmentService{
		DB:         db,
		LogService: logSvc,
	}
}

// checkBarberAvailabilityTx...
func (s *appointmentService) checkBarberAvailabilityTx(
	tx *gorm.DB,
	tenantID, barberID uint,
	start, end time.Time,
) (bool, error) {
	// 1) Normalize to UTC
	start = start.UTC()
	end = end.UTC()

	// 2) Lock the barber row to prevent concurrent bookings
	var barber barberBookingModels.Barber
	if err := tx.
		Clauses(clause.Locking{Strength: "UPDATE"}).
		Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", barberID, tenantID).
		First(&barber).Error; err != nil {

		if errors.Is(err, gorm.ErrRecordNotFound) {
			// ช่างไม่อยู่ในระบบ ให้ถือว่าไม่ว่างเลย
			return false, nil
		}
		return false, err
	}

	// 3) ตรวจสอบ overlap กับ existing appointments
	var count int64
	if err := tx.
		Model(&barberBookingModels.Appointment{}).
		Where("tenant_id = ? AND barber_id = ? AND status IN ? AND deleted_at IS NULL",
			tenantID, barberID,
			[]string{
				string(barberBookingModels.StatusPending),
				string(barberBookingModels.StatusConfirmed),
			},
		).
		// Time comparisons in UTC
		Where("start_time < ? AND end_time > ?", end, start).
		Count(&count).Error; err != nil {
		return false, err
	}

	return count == 0, nil
}

func (s *appointmentService) CheckBarberAvailability(
	ctx context.Context,
	tenantID, barberID uint,
	start, end time.Time,
) (bool, error) {
	// เริ่ม transaction เพื่อให้ Lock ใช้ได้จริง
	var available bool
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var err error
		available, err = s.checkBarberAvailabilityTx(tx, tenantID, barberID, start, end)
		return err
	})
	return available, err
}

func (s *appointmentService) CreateAppointment(
	ctx context.Context,
	input *barberBookingModels.Appointment,
) (*barberBookingModels.Appointment, error) {
	// ตรวจ input เบื้องต้น
	if input == nil {
		return nil, errors.New("input appointment data is required")
	}
	if input.TenantID == 0 || input.BranchID == 0 || input.ServiceID == 0 || input.CustomerID == 0 || input.StartTime.IsZero() {
		return nil, errors.New("missing required fields")
	}

	var appt *barberBookingModels.Appointment

	// 1. สร้าง appointment ภายใน transaction
	err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 0. ตรวจว่า branch มีอยู่และสังกัด tenant เดียวกัน
		var branch coreModels.Branch
		if err := tx.
			Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", input.BranchID, input.TenantID).
			First(&branch).Error; err != nil {
			return fmt.Errorf("branch not found or access denied")
		}

		// 1. ดึง service + ตรวจ tenant
		var service barberBookingModels.Service
		if err := tx.
			Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", input.ServiceID, input.TenantID).
			First(&service).Error; err != nil {
			return fmt.Errorf("service not found or access denied")
		}
		if service.Duration <= 0 {
			return fmt.Errorf("duration must be > 0")
		}

		// 2. คำนวณ EndTime
		startTime := input.StartTime
		endTime := startTime.Add(time.Duration(service.Duration) * time.Minute)
		input.EndTime = endTime

		// 3. ตรวจ availability ถ้ามี barber_id
		if input.BarberID != 0 {
			var barber barberBookingModels.Barber
			if err := tx.
				Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", input.BarberID, input.TenantID).
				First(&barber).Error; err != nil {
				return fmt.Errorf("barber not found or mismatched branch")
			}
			if barber.BranchID != input.BranchID {
				return fmt.Errorf("barber not found or mismatched branch")
			}
			available, err := s.checkBarberAvailabilityTx(
				tx, input.TenantID, input.BarberID, startTime, endTime,
			)
			if err != nil {
				return fmt.Errorf("check barber availability failed: %w", err)
			}
			if !available {
				return fmt.Errorf("barber is not available during this time")
			}
		}

		// 4. ตั้งค่า Status/Timestamps แล้วสร้างแถว appointment
		if input.Status == "" {
			input.Status = barberBookingModels.StatusPending
		}
		now := time.Now().UTC()
		input.CreatedAt = now
		input.UpdatedAt = now

		if err := tx.Create(input).Error; err != nil {
			return fmt.Errorf("failed to create appointment: %w", err)
		}

		// เซ็ตผลลัพธ์เพื่อคืนค่าหลัง transaction
		appt = input
		return nil
	})
	if err != nil {
		return nil, err
	}

	
	// 2. นอก transaction: เขียน status log
	var userID *uint
	var custID *uint
	if appt.UserID != nil {
		userID = appt.UserID
	} else {
		custID = &appt.CustomerID
	}
	log.Printf(">>> writing initial creation log for appt=%d", appt.ID) 
	if logErr := s.LogService.LogStatusChange(
		ctx,
		appt.ID,
		"",                  // oldStatus (ยังไม่มีสถานะก่อนหน้า)
		string(appt.Status), // newStatus
		userID,
		custID,
		"initial creation",
	); logErr != nil {
		// แค่เตือนใน log ฝั่งเซิร์ฟเวอร์ ไม่ rollback appointment
		log.Printf("warning: failed to write status log: %v", logErr)
	}

	return appt, nil
}

func (s *appointmentService) GetAvailableBarbers(ctx context.Context, tenantID, branchID uint, start, end time.Time) ([]barberBookingModels.Barber, error) {
	var barbers []barberBookingModels.Barber

	activeStatuses := []barberBookingModels.AppointmentStatus{
		barberBookingModels.StatusPending,
		barberBookingModels.StatusConfirmed,
		barberBookingModels.StatusRescheduled,
	}

	tx := s.DB.WithContext(ctx)

	subQuery := tx.Model(&barberBookingModels.Appointment{}).
		Select("barber_id").
		Where(`
			tenant_id = ?
			AND barber_id IS NOT NULL
			AND status IN ?
			AND NOT (end_time <= ? OR start_time >= ?)
		`, tenantID, activeStatuses, start, end)

	err := tx.Model(&barberBookingModels.Barber{}).
		Where("tenant_id = ? AND branch_id = ? AND id NOT IN (?)", tenantID, branchID, subQuery).
		Find(&barbers).Error

	return barbers, err
}

func (s *appointmentService) UpdateAppointment(
    ctx context.Context,
    id uint,
    tenantID uint,
    input *barberBookingModels.Appointment,
) (*barberBookingModels.Appointment, error) {
    if input == nil {
        return nil, errors.New("input appointment data is required")
    }

    var updatedAppt *barberBookingModels.Appointment

    err := s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
        // 1. โหลด appointment ปัจจุบัน
        var ap barberBookingModels.Appointment
        if err := tx.
            Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", id, tenantID).
            First(&ap).Error; err != nil {
            return fmt.Errorf("appointment not found")
        }
        oldStatus := ap.Status

        // 2. อัปเดต ServiceID และคำนวณ EndTime ใหม่ถ้า service เปลี่ยน
        if input.ServiceID != 0 && input.ServiceID != ap.ServiceID {
            // ตรวจ service + ดึง duration
            var svc barberBookingModels.Service
            if err := tx.
                Where("id = ? AND tenant_id = ? AND deleted_at IS NULL", input.ServiceID, tenantID).
                First(&svc).Error; err != nil {
                return fmt.Errorf("service not found or access denied")
            }
            ap.ServiceID = input.ServiceID
            // recalc end time บน startTime ปัจจุบัน
            ap.EndTime = ap.StartTime.Add(time.Duration(svc.Duration) * time.Minute)
        }

        // 3. ถ้าเปลี่ยน startTime ให้ recalc EndTime ด้วยเช่นกัน
        if !input.StartTime.IsZero() && !input.StartTime.Equal(ap.StartTime) {
            ap.StartTime = input.StartTime
            // ดึง duration ของ service เดิม
            var svc barberBookingModels.Service
            if err := tx.
                Where("id = ?", ap.ServiceID).
                First(&svc).Error; err != nil {
                return fmt.Errorf("failed fetching service for recalc end time")
            }
            ap.EndTime = ap.StartTime.Add(time.Duration(svc.Duration) * time.Minute)
        }

        // 4. อัปเดต BarberID, CustomerID, Status, Notes ตาม input
        if input.BarberID != 0 {
            ap.BarberID = input.BarberID
        }
        if input.CustomerID != 0 && input.CustomerID != ap.CustomerID {
            ap.CustomerID = input.CustomerID
        }
        if input.Status != "" && input.Status != ap.Status {
            ap.Status = input.Status
        }
        if input.Notes != "" {
            ap.Notes = input.Notes
        }
        ap.UpdatedAt = time.Now().UTC()

        // 5. Save appointment
        if err := tx.Save(&ap).Error; err != nil {
            return fmt.Errorf("failed to update appointment: %w", err)
        }

        // 6. Log status change ถ้ามีการเปลี่ยนสถานะ
        if oldStatus != ap.Status {
            var userID *uint
            var custID *uint
            if input.UserID != nil {
                userID = input.UserID
            } else {
                custID = &ap.CustomerID
            }
            if err := s.LogService.LogStatusChange(
                ctx,
                ap.ID,
                string(oldStatus),
                string(ap.Status),
                userID,
                custID,
                "status updated via API",
            ); err != nil {
                return fmt.Errorf("failed to log status change: %w", err)
            }
        }

        // 7. ดึงข้อมูลใหม่พร้อม Preload relations
        var out barberBookingModels.Appointment
        if err := tx.
            Preload("Service").
            Preload("Customer").
            First(&out, ap.ID).Error; err != nil {
            return fmt.Errorf("failed to fetch updated appointment: %w", err)
        }
        updatedAppt = &out
        return nil
    })

    return updatedAppt, err
}


func (s *appointmentService) GetAppointmentByID(ctx context.Context, id uint) (*barberBookingModels.Appointment, error) {
	var appt barberBookingModels.Appointment
	// ดึงเฉพาะเรคคอร์ดที่ยังไม่ถูกลบ (deleted_at IS NULL)
	err := s.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", id).
		First(&appt).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("appointment with ID %d not found", id)
		}
		return nil, fmt.Errorf("failed to fetch appointment: %w", err)
	}
	return &appt, nil
}

func (s *appointmentService) ListAppointments(ctx context.Context, filter barberBookingDto.AppointmentFilter) ([]barberBookingModels.Appointment, error) {
    var appointments []barberBookingModels.Appointment

    // เริ่มต้น tx พร้อม preload ความสัมพันธ์ทั้งหมด
    tx := s.DB.WithContext(ctx).Debug().
        Model(&barberBookingModels.Appointment{}).
        Preload("Service").
        Preload("Customer").
        Preload("Barber").

        Where("tenant_id = ?", filter.TenantID)

    // กรองตามเงื่อนไขต่าง ๆ
    if filter.BranchID != nil {
        tx = tx.Where("branch_id = ?", *filter.BranchID)
    }
    if filter.BarberID != nil {
        tx = tx.Where("barber_id = ?", *filter.BarberID)
    }
    if filter.CustomerID != nil {
        tx = tx.Where("customer_id = ?", *filter.CustomerID)
    }
    if filter.Status != nil {
        tx = tx.Where("status = ?", *filter.Status)
    }
    if filter.StartDate != nil {
        tx = tx.Where("start_time >= ?", *filter.StartDate)
    }
    if filter.EndDate != nil {
        tx = tx.Where("end_time <= ?", *filter.EndDate)
    }

    // กำหนดการจัดเรียง กรณีมี sort_by หรือ default เรียงตาม start_time
    if filter.SortBy != nil && *filter.SortBy != "" {
        tx = tx.Order(*filter.SortBy)
    } else {
        tx = tx.Order("start_time asc")
    }

    // กำหนด Limit/Offset เพื่อ Pagination (ถ้ามี)
    if filter.Limit != nil {
        tx = tx.Limit(*filter.Limit)
    }
    if filter.Offset != nil {
        tx = tx.Offset(*filter.Offset)
    }

    // สุดท้าย execute Find เพียงครั้งเดียว
    if err := tx.Find(&appointments).Error; err != nil {
        return nil, err
    }
    return appointments, nil
}

func (s *appointmentService) ListAppointmentsResponse(ctx context.Context, filter barberBookingDto.AppointmentFilter) ([]barberBookingPort.AppointmentResponse, error) {
    // 1. เรียก ListAppointments ดึงข้อมูล full
    apps, err := s.ListAppointments(ctx, filter)
    if err != nil {
        return nil, err
    }

    // 2. สร้าง slice ของ DTO
    var result []barberBookingPort.AppointmentResponse
    for _, a := range apps {
        // แปลงเวลาเป็น ISO string (format RFC3339)
        startISO := a.StartTime.Format(time.RFC3339)
        endISO := a.EndTime.Format(time.RFC3339)

        // สร้าง object ตัวนึง
        ar := barberBookingPort.AppointmentResponse{
            ID:       a.ID,
            BranchID: a.BranchID,
            TenantID: a.TenantID,
            StartTime: startISO,
            EndTime:   endISO,
            Status:   string(a.Status),
            Notes:    a.Notes,
        }

        // map ส่วน Service
        ar.Service.ID = a.Service.ID
        ar.Service.Name = a.Service.Name
        ar.Service.Duration = a.Service.Duration
        ar.Service.Price = a.Service.Price

        // map ส่วน Barber (coreModels.User)
        ar.Barber.ID = a.Barber.ID
        ar.Barber.Username = a.Barber.Username
        ar.Barber.Email = a.Barber.Email

        // map ส่วน Customer
        ar.Customer.ID = a.Customer.ID
        ar.Customer.Name = a.Customer.Name
        ar.Customer.Phone = a.Customer.Phone
        ar.Customer.Email = a.Customer.Email

        // เติมลง slice
        result = append(result, ar)
    }

    return result, nil
}

func (s *appointmentService) CancelAppointment(
	ctx context.Context,
	appointmentID uint,
	actorUserID *uint,
	actorCustomerID *uint,
) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var ap barberBookingModels.Appointment
		if err := tx.
			Where("id = ? AND deleted_at IS NULL", appointmentID).
			First(&ap).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("appointment with ID %d not found", appointmentID)
			}
			return err
		}

		oldStatus := ap.Status
		if oldStatus == barberBookingModels.StatusComplete ||
			oldStatus == barberBookingModels.StatusCancelled {
			return errors.New("appointment cannot be cancelled in its current status")
		}

		// อัปเดตสถานะ
		ap.Status = barberBookingModels.StatusCancelled
		ap.UpdatedAt = time.Now().UTC()

		// บันทึกใครยกเลิก
		if actorUserID != nil {
			ap.UserID = actorUserID
		} else if actorCustomerID != nil {
			ap.CustomerID = *actorCustomerID // หรือเก็บลงฟิลด์อื่นถ้ามี
		}

		if err := tx.Save(&ap).Error; err != nil {
			return err
		}

		// เขียน log
		if err := s.LogService.LogStatusChange(
			ctx,
			appointmentID,
			string(oldStatus),
			string(ap.Status),
			actorUserID,
			actorCustomerID,
			"cancelled via API",
		); err != nil {
			return err
		}

		return nil
	})
}

func (s *appointmentService) RescheduleAppointment(
	ctx context.Context,
	appointmentID uint,
	newStartTime time.Time,
	actorUserID *uint,
	actorCustomerID *uint,
) error {
	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) Load appointment + Service
		var ap barberBookingModels.Appointment
		if err := tx.
			Preload("Service").
			Where("id = ? AND deleted_at IS NULL", appointmentID).
			First(&ap).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("appointment with ID %d not found", appointmentID)
			}
			return err
		}

		// 2) Only PENDING/CONFIRMED may be rescheduled
		if ap.Status == barberBookingModels.StatusComplete ||
			ap.Status == barberBookingModels.StatusCancelled {
			return fmt.Errorf("cannot reschedule a completed or cancelled appointment")
		}

		// remember old values
		oldStart := ap.StartTime
		oldEnd := ap.EndTime
		oldStatus := ap.Status

		// 3) Conflict check
		newEndTime := newStartTime.Add(time.Duration(ap.Service.Duration) * time.Minute)
		var conflict int64
		if err := tx.Model(&barberBookingModels.Appointment{}).
			Where(`barber_id = ? AND branch_id = ? AND id != ? 
                   AND status IN ? AND start_time < ? AND end_time > ?`,
				ap.BarberID, ap.BranchID, ap.ID,
				[]barberBookingModels.AppointmentStatus{
					barberBookingModels.StatusPending,
					barberBookingModels.StatusConfirmed,
				},
				newEndTime, newStartTime,
			).
			Count(&conflict).Error; err != nil {
			return err
		}
		if conflict > 0 {
			return fmt.Errorf("cannot reschedule: time slot conflicts with another appointment")
		}

		// 4) Apply new times & updater
		ap.StartTime = newStartTime
		ap.EndTime = newEndTime
		ap.UserID = actorUserID // may be nil
		ap.UpdatedAt = time.Now().UTC()
		ap.Status = barberBookingModels.StatusConfirmed

		if actorCustomerID != nil {
			// if customer themselves rescheduled, record that
			ap.UserID = nil
		}

		if err := tx.Save(&ap).Error; err != nil {
			return fmt.Errorf("failed to save rescheduled appointment: %w", err)
		}

		// 5) Log status-change (if it actually changed)
		if oldStatus != ap.Status {
			if err := s.LogService.LogStatusChange(
				ctx,
				ap.ID,
				string(oldStatus),
				string(ap.Status),
				actorUserID,
				actorCustomerID,
				"status updated via reschedule",
			); err != nil {
				return fmt.Errorf("failed to log status change: %w", err)
			}
		}

		// 6) Log the actual timeslot change
		note := fmt.Sprintf(
			"rescheduled from %s–%s to %s–%s",
			oldStart.Format(time.RFC3339), oldEnd.Format(time.RFC3339),
			ap.StartTime.Format(time.RFC3339), ap.EndTime.Format(time.RFC3339),
		)
		if err := s.LogService.LogStatusChange(
			ctx,
			ap.ID,
			"", // no status change
			"",
			actorUserID,
			actorCustomerID,
			note,
		); err != nil {
			return fmt.Errorf("failed to log reschedule detail: %w", err)
		}

		return nil
	})
}

func (s *appointmentService) GetAppointmentsByBarber(
	ctx context.Context,
	barberID uint,
	start *time.Time,
	end *time.Time,
) ([]barberBookingModels.Appointment, error) {
	// 0. validate IDs
	if barberID == 0 {
		return nil, fmt.Errorf("invalid barberID: %d", barberID)
	}
	// 1. optional: ถ้าส่งทั้ง start+end มา ให้ตรวจ start <= end
	if start != nil && end != nil && start.After(*end) {
		return nil, fmt.Errorf("start time %v is after end time %v", *start, *end)
	}

	// 2. ตรวจว่า barber ยังอยู่ในระบบ (และยังไม่ soft-deleted)
	var barber barberBookingModels.Barber
	if err := s.DB.WithContext(ctx).
		Where("id = ? AND deleted_at IS NULL", barberID).
		First(&barber).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("barber with ID %d not found", barberID)
		}
		return nil, fmt.Errorf("failed to lookup barber: %w", err)
	}

	// 3. Build query
	q := s.DB.WithContext(ctx).
		Model(&barberBookingModels.Appointment{}).
		Where("barber_id = ?", barberID)

	if start != nil {
		q = q.Where("start_time >= ?", *start)
	}
	if end != nil {
		q = q.Where("end_time <= ?", *end)
	}

	// 4. Optionally limit how many you return if no filters (to prevent full table scan)
	if start == nil && end == nil {
		q = q.Limit(1000) // หรือพาราม  config มากำหนด
	}

	// 5. Execute
	var appts []barberBookingModels.Appointment
	if err := q.Order("start_time ASC").Find(&appts).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch appointments: %w", err)
	}
	return appts, nil
}

func (s *appointmentService) DeleteAppointment(
	ctx context.Context,
	appointmentID uint,
) error {
	if appointmentID == 0 {
		return fmt.Errorf("invalid appointment ID: %d", appointmentID)
	}

	return s.DB.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1) ตรวจว่ามี appointment จริงๆ หรือเปล่า
		var ap barberBookingModels.Appointment
		if err := tx.
			Where("id = ? AND deleted_at IS NULL", appointmentID).
			First(&ap).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return fmt.Errorf("appointment with ID %d not found", appointmentID)
			}
			return err
		}

		// 3) Soft-delete appointment
		if err := tx.Delete(&ap).Error; err != nil {
			return err
		}

		// 4) ลบ log ที่เกี่ยวข้อง (ถ้าต้องการ clean-up)
		if err := tx.
			Where("appointment_id = ?", appointmentID).
			Delete(&barberBookingModels.AppointmentStatusLog{}).
			Error; err != nil {
			return err
		}

		// commit transaction
		return nil
	})
}

func (s *appointmentService) GetUpcomingAppointmentsByCustomer(ctx context.Context, customerID uint) (*barberBookingModels.Appointment, error) {
	var appointment barberBookingModels.Appointment
	err := s.DB.WithContext(ctx).
		Where("customer_id = ? AND start_time > ? AND status IN ?", customerID, time.Now(), []string{
			string(barberBookingModels.StatusPending),
			string(barberBookingModels.StatusConfirmed),
			string(barberBookingModels.StatusRescheduled),
		}).
		Order("start_time ASC").
		First(&appointment).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil // ไม่มีคิวถัดไป
	}
	return &appointment, err
}

// เก็บไว้ก่อนเอาไว้ใช้ตอนระบบเริ่มใหญ่
func (s *appointmentService) CalculateAppointmentEndTime(ctx context.Context, serviceID uint, startTime time.Time) (time.Time, error) {
	var service barberBookingModels.Service
	err := s.DB.WithContext(ctx).First(&service, serviceID).Error
	if err != nil {
		return time.Time{}, err
	}

	if service.Duration < 0 {
		return time.Time{}, fmt.Errorf("invalid service duration")
	}

	duration := time.Duration(service.Duration) * time.Minute
	endTime := startTime.Add(duration)
	return endTime, nil
}
