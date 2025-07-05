package barberBookingController

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	helperFunc "myapp/modules/barberbooking"
	barberBookingDto "myapp/modules/barberbooking/dto"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"
)

// AppointmentController handles endpoints related to appointments
type AppointmentController struct {
	Service barberBookingPort.IAppointment
}

// NewAppointmentController creates a new controller
func NewAppointmentController(service barberBookingPort.IAppointment) *AppointmentController {
	return &AppointmentController{Service: service}
}

// CheckBarberAvailability godoc
// @Summary      ตรวจสอบความพร้อมของช่างตัดผม
// @Description  ตรวจสอบว่า Barber ที่ระบุสามารถรับงานได้ในช่วงเวลาที่กำหนด (RFC3339) หรือไม่
// @Tags         Appointment
// @Produce      json
// @Param        tenant_id   path      int     true  "รหัส Tenant"
// @Param        barber_id   path      int     true  "รหัส Barber"
// @Param        start       query     string  true  "เวลาเริ่มต้น (RFC3339) เช่น 2025-05-29T09:00:00Z"
// @Param        end         query     string  true  "เวลาสิ้นสุด (RFC3339) เช่น 2025-05-29T10:00:00Z"
// @Success      200         {object}  map[string]interface{}  "คืนค่า status และ available (true/false)"
// @Failure      400         {object}  map[string]string       "พารามิเตอร์ไม่ถูกต้อง หรือขาด start/end"
// @Failure      500         {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างตรวจสอบความพร้อม"
// @Router       /tenants/:tenant_id/barbers/:barber_id/availability [get]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) CheckBarberAvailability(c *fiber.Ctx) error {
	// 1. Parse tenant_id from path
	tID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
	}

	// 2. Parse barber_id from path
	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid barber_id"})
	}

	// 3. Parse start & end from query
	startStr := c.Query("start", "")
	if startStr == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing start time"})
	}
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid start time format. Expect RFC3339"})
	}

	endStr := c.Query("end", "")
	if endStr == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing end time"})
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid end time format. Expect RFC3339"})
	}

	// 4. Call service
	available, err := ctrl.Service.CheckBarberAvailability(
		c.Context(), // หรือ context.Background()
		tID,
		barberID,
		startTime,
		endTime,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "Failed to check availability", "error": err.Error()})
	}

	// 5. Return result
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "available": available})
}

// CreateAppointment godoc
// @Summary      สร้างนัดหมายใหม่ (Create Appointment)
// @Description  สร้าง Appointment ภายใต้ Tenant ที่ระบุ พร้อมกรอก branch, service, customer, optional barber, start_time (RFC3339) และ notes
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param tenant_id path uint true "รหัส Tenant"
// @Param body body barberBookingPort.CreateAppointmentRequest true "Payload สำหรับสร้างนัดหมาย"
// @Success      201         {object}  barberBookingModels.Appointment            "คืนค่า status success พร้อมข้อมูล Appointment ที่สร้าง"
// @Failure      400         {object}  map[string]string                          "Missing required fields หรือ Invalid format"
// @Failure      500         {object}  map[string]string                          "Internal Server Error"
// @Router       /tenants/{tenant_id}/appointments [post]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) CreateAppointment(c *fiber.Ctx) error {
	// 1. Parse tenant_id from URL
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant_id",
		})
	}

	// 2. Parse JSON body
	var payload struct {
		BranchID   uint   `json:"branch_id"`
		ServiceID  uint   `json:"service_id"`
		BarberID   uint   `json:"barber_id,omitempty"`
		CustomerID uint   `json:"customer_id"`
		StartTime  string `json:"start_time"`
		Notes      string `json:"notes,omitempty"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// 3. Validate required fields
	if payload.BranchID == 0 || payload.ServiceID == 0 || payload.CustomerID == 0 || payload.StartTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing required fields",
		})
	}

	// 4. Parse start_time
	startTime, err := time.Parse(time.RFC3339, payload.StartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid start_time format. Expect RFC3339",
		})
	}

	// 5. Build appointment model to send to service layer
	appt := &barberBookingModels.Appointment{
		TenantID:   tenantID,
		BranchID:   payload.BranchID,
		ServiceID:  payload.ServiceID,
		BarberID:   payload.BarberID,
		CustomerID: payload.CustomerID,
		StartTime:  startTime,
		Notes:      payload.Notes,
	}

	// 6. Call service
	createdDTO, err := ctrl.Service.CreateAppointment(c.Context(), appt) // 🔁 Return *AppointmentResponseDTO
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// 7. Return simplified DTO as JSON
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"status": "success",
		"data":   createdDTO,
	})
}

// GetAvailableBarbers handles GET /tenants/:tenant_id/branches/:branch_id/available-barbers?start=&end=
// GetAvailableBarbers godoc
// @Summary      ดึงรายชื่อช่างตัดผมว่าง
// @Description  คืนรายการ Barber ที่ว่างในช่วงเวลาที่กำหนด (RFC3339) ภายใต้ Tenant และ Branch ที่ระบุ
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint     true  "รหัส Tenant"
// @Param        branch_id  path      uint     true  "รหัส Branch"
// @Param        start      query     string   true  "เวลาเริ่มต้น (RFC3339) เช่น 2025-05-29T09:00:00Z"
// @Param        end        query     string   true  "เวลาสิ้นสุด (RFC3339) เช่น 2025-05-29T10:00:00Z"
// @Success      200        {object}  map[string][]barberBookingModels.Barber  "คืนค่า status success และ array ของ Barber ที่ว่าง"
// @Failure      400        {object}  map[string]string                        "Missing or invalid parameters"
// @Failure      500        {object}  map[string]string                        "Internal Server Error"
// @Router       /tenants/:tenant_id/branches/:branch_id/available-barbers [get]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) GetAvailableBarbers(c *fiber.Ctx) error {
	// 1. Parse tenant_id
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
	}
	// 2. Parse branch_id
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid branch_id"})
	}
	// 3. Parse start time
	startStr := c.Query("start", "")
	if startStr == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing start time"})
	}
	startTime, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid start time format. Expect RFC3339"})
	}
	// 4. Parse end time
	endStr := c.Query("end", "")
	if endStr == "" {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Missing end time"})
	}
	endTime, err := time.Parse(time.RFC3339, endStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).
			JSON(fiber.Map{"status": "error", "message": "Invalid end time format. Expect RFC3339"})
	}
	// 5. Call service
	barbers, err := ctrl.Service.GetAvailableBarbers(
		c.Context(), tenantID, branchID, startTime, endTime,
	)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).
			JSON(fiber.Map{"status": "error", "message": "Failed to get available barbers", "error": err.Error()})
	}
	// 6. Return JSON
	return c.Status(fiber.StatusOK).
		JSON(fiber.Map{"status": "success", "data": barbers})
}

// PUT /tenants/:tenant_id/appointments/:appointment_id
// UpdateAppointment godoc
// @Summary      แก้ไขข้อมูลนัดหมาย
// @Description  อัปเดต Appointment ตามรหัสที่ระบุ ภายใต้ Tenant ที่กำหนด โดยรับข้อมูล JSON ของ Appointment ใหม่
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id        path      uint                                         true  "รหัส Tenant"
// @Param        appointment_id   path      uint                                         true  "รหัส Appointment"
// @Param        body             body      barberBookingModels.Appointment              true  "ข้อมูล Appointment ที่ต้องการอัปเดต"
// @Success      200              {object}  barberBookingModels.Appointment              "คืนค่า status success และข้อมูล Appointment ที่อัปเดต"
// @Failure      400              {object}  map[string]string                             "Invalid tenant_id, appointment_id หรือ JSON body"
// @Failure      500              {object}  map[string]string                             "Internal Server Error"
// @Router       /tenants/:tenant_id/appointments/:appointment_id [put]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) UpdateAppointment(c *fiber.Ctx) error {
	// 1. Parse tenant_id
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
	}

	// 2. Parse appointment_id
	apptID, err := helperFunc.ParseUintParam(c, "appointment_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid appointment_id"})
	}

	// 3. Parse JSON body
	var input barberBookingModels.Appointment
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid JSON body"})
	}

	// 5. Call service
	updated, err := ctrl.Service.UpdateAppointment(context.Background(), apptID, tenantID, &input)
	if err != nil {
		// service returns generic fmt.Errorf with message
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	// 6. Return success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   updated,
	})
}

// GET /tenants/:tenant_id/appointments/:appointment_id
// GetAppointmentByID godoc
// @Summary      ดึงข้อมูลนัดหมายตาม ID
// @Description  คืนรายละเอียด Appointment ตามรหัสที่ระบุ ภายใต้ Tenant สำหรับความสอดคล้องของ URL
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id       path      uint   true  "รหัส Tenant"
// @Param        appointment_id  path      uint   true  "รหัส Appointment"
// @Success      200             {object}  barberBookingModels.Appointment  "คืนค่า status success และข้อมูล Appointment"
// @Failure      400             {object}  map[string]string               "Invalid tenant_id หรือ appointment_id"
// @Failure      404             {object}  map[string]string               "Appointment not found"
// @Failure      500             {object}  map[string]string               "Failed to fetch appointment"
// @Router       /tenants/:tenant_id/appointments/:appointment_id [get]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) GetAppointmentByID(c *fiber.Ctx) error {
	// 1. Parse tenant_id (for URL consistency; we don't pass it to service)
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
	appt, err := ctrl.Service.GetAppointmentByID(context.Background(), apptID)
	if err != nil {
		// Distinguish not found vs other errors
		if strings.Contains(err.Error(), "not found") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch appointment",
			"error":   err.Error(),
		})
	}

	// 4. Return success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   appt,
	})
}

// GET /tenants/:tenant_id/appointments
// ListAppointments godoc
// @Summary      ดึงรายการนัดหมายของ Tenant
// @Description  คืนรายการ Appointment ของ Tenant ที่ระบุ โดยรองรับการกรองตาม branch_id, barber_id, customer_id, status, ช่วงวัน (RFC3339), pagination (limit, offset) และการจัดเรียง (sort_by)
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id    path      uint     true   "รหัส Tenant"
// @Param        branch_id    query     uint     false  "กรองตาม Branch ID"
// @Param        barber_id    query     uint     false  "กรองตาม Barber ID"
// @Param        customer_id  query     uint     false  "กรองตาม Customer ID"
// @Param        status       query     string   false  "กรองตามสถานะ Appointment"
// @Param        start_date   query     string   false  "กรองวันที่เริ่มต้น (RFC3339)"
// @Param        end_date     query     string   false  "กรองวันที่สิ้นสุด (RFC3339)"
// @Param        limit        query     int      false  "จำนวนรายการสูงสุด (pagination)"
// @Param        offset       query     int      false  "เลื่อนข้ามรายการ (pagination)"
// @Param        sort_by      query     string   false  "จัดเรียงผลลัพธ์ เช่น start_time desc"
// @Success      200          {object}  map[string][]barberBookingModels.Appointment  "คืนค่า status success และ array ของ Appointment ใน key `data`"
// @Failure      400          {object}  map[string]string                            "Invalid query parameters"
// @Failure      500          {object}  map[string]string                            "Internal Server Error"
// @Router       /tenants/:tenant_id/appointments [get]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) ListAppointments(c *fiber.Ctx) error {
	// 1. Parse tenant_id
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
	}

	// 2. Build filter
	var f barberBookingDto.AppointmentFilter
	f.TenantID = tenantID

	if qs := c.Query("branch_id", ""); qs != "" {
		if v, err := strconv.ParseUint(qs, 10, 64); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid branch_id"})
		} else {
			u := uint(v)
			f.BranchID = &u
		}
	}
	if qs := c.Query("barber_id", ""); qs != "" {
		if v, err := strconv.ParseUint(qs, 10, 64); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid barber_id"})
		} else {
			u := uint(v)
			f.BarberID = &u
		}
	}
	if qs := c.Query("customer_id", ""); qs != "" {
		if v, err := strconv.ParseUint(qs, 10, 64); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid customer_id"})
		} else {
			u := uint(v)
			f.CustomerID = &u
		}
	}
	if qs := c.Query("status", ""); qs != "" {
		s := barberBookingModels.AppointmentStatus(qs)
		f.Status = &s
	}
	if qs := c.Query("start_date", ""); qs != "" {
		t, err := time.Parse(time.RFC3339, qs)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid start_date format. Expect RFC3339"})
		}
		f.StartDate = &t
	}
	if qs := c.Query("end_date", ""); qs != "" {
		t, err := time.Parse(time.RFC3339, qs)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid end_date format. Expect RFC3339"})
		}
		f.EndDate = &t
	}
	if qs := c.Query("limit", ""); qs != "" {
		if v, err := strconv.Atoi(qs); err != nil || v < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid limit"})
		} else {
			f.Limit = &v
		}
	}
	if qs := c.Query("offset", ""); qs != "" {
		if v, err := strconv.Atoi(qs); err != nil || v < 0 {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid offset"})
		} else {
			f.Offset = &v
		}
	}
	if qs := c.Query("sort_by", ""); qs != "" {
		f.SortBy = &qs
	}

	// 3. Call service for DTO
	apptResp, err := ctrl.Service.ListAppointmentsResponse(context.Background(), f)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to list appointments",
			"error":   err.Error(),
		})
	}

	// 4. Return slimmed-down DTO
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   apptResp,
	})
}

type CancelRequest struct {
	ActorUserID     *uint `json:"actor_user_id,omitempty"`
	ActorCustomerID *uint `json:"actor_customer_id,omitempty"`
}

// POST /tenants/:tenant_id/appointments/:appointment_id/cancel
// CancelAppointment godoc
// @Summary      ยกเลิกนัดหมาย
// @Description  ยกเลิก Appointment ตามรหัสที่ระบุ ภายใต้ Tenant ที่กำหนด โดยระบุ Actor เป็น User หรือ Customer ได้ครั้งละหนึ่งคน
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id         path      uint             true  "รหัส Tenant"
// @Param        appointment_id    path      uint             true  "รหัส Appointment"
// @Param        body              body      CancelRequest    true  "ข้อมูล Actor: ระบุ ActorUserID หรือ ActorCustomerID อย่างใดอย่างหนึ่ง"
// @Success      200               {object}  map[string]string  "คืนค่า status success และข้อความยืนยันการยกเลิก"
// @Failure      400               {object}  map[string]string  "Missing or invalid parameters หรือ cannot be cancelled"
// @Failure      404               {object}  map[string]string  "Appointment not found"
// @Failure      500               {object}  map[string]string  "Internal Server Error"
// @Router       /tenants/:tenant_id/appointments/:appointment_id/cancel [post]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) CancelAppointment(c *fiber.Ctx) error {
	// 1. Parse path params
	if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
	}
	apptID, err := helperFunc.ParseUintParam(c, "appointment_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid appointment_id"})
	}

	// 2. Parse body
	var req CancelRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": "Invalid JSON body"})
	}

	// 3. Validate that exactly one actor is provided
	if (req.ActorUserID == nil && req.ActorCustomerID == nil) ||
		(req.ActorUserID != nil && req.ActorCustomerID != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Either actor_user_id or actor_customer_id must be provided, but not both",
		})
	}

	// 4. Call service (assume signature has been updated to accept both)
	err = ctrl.Service.CancelAppointment(
		context.Background(),
		apptID,
		req.ActorUserID,
		req.ActorCustomerID,
	)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "not found"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status": "error", "message": msg})
		case strings.Contains(msg, "cannot be cancelled"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status": "error", "message": msg})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status": "error", "message": "Failed to cancel appointment", "error": msg})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"status": "success", "message": "appointment cancelled"})
}

type RescheduleRequest struct {
	NewStartTime    string `json:"new_start_time"`
	ActorUserID     *uint  `json:"actor_user_id,omitempty"`
	ActorCustomerID *uint  `json:"actor_customer_id,omitempty"`
}

// POST /tenants/:tenant_id/appointments/:appointment_id/reschedule
// RescheduleAppointment godoc
// @Summary      เปลี่ยนเวลานัดหมาย
// @Description  Reschedule Appointment ตามรหัสที่ระบุ ภายใต้ Tenant ที่กำหนด โดยระบุเวลาใหม่และ Actor (User หรือ Customer)
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id       path      uint               true  "รหัส Tenant"
// @Param        appointment_id  path      uint               true  "รหัส Appointment"
// @Param        body            body      RescheduleRequest  true  "ข้อมูลสำหรับเลื่อนนัดหมาย (new_start_time, actor_user_id หรือ actor_customer_id)"
// @Success      200             {object}  map[string]string  "คืนค่า status success และข้อความยืนยันการเลื่อนนัดหมาย"
// @Failure      400             {object}  map[string]string  "Missing or invalid parameters หรือ cannot reschedule"
// @Failure      404             {object}  map[string]string  "Appointment not found"
// @Failure      500             {object}  map[string]string  "Internal Server Error"
// @Router       /tenants/:tenant_id/appointments/:appointment_id/reschedule [post]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) RescheduleAppointment(c *fiber.Ctx) error {
	// 1. Parse path params
	if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant_id",
		})
	}
	apptID, err := helperFunc.ParseUintParam(c, "appointment_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid appointment_id",
		})
	}

	// 2. Parse request body
	var req RescheduleRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid JSON body",
		})
	}

	// 3. Validate new_start_time
	if req.NewStartTime == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing new_start_time",
		})
	}
	newStart, err := time.Parse(time.RFC3339, req.NewStartTime)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid new_start_time format, expect RFC3339",
		})
	}

	// 4. Validate actor: exactly one of user or customer
	if (req.ActorUserID == nil && req.ActorCustomerID == nil) ||
		(req.ActorUserID != nil && req.ActorCustomerID != nil) {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Either actor_user_id or actor_customer_id must be provided, but not both",
		})
	}

	// 5. Call service
	err = ctrl.Service.RescheduleAppointment(
		c.Context(),
		apptID,
		newStart,
		req.ActorUserID,
		req.ActorCustomerID,
	)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "not found"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": msg,
			})
		case strings.Contains(msg, "cannot reschedule"):
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": msg,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": "Failed to reschedule appointment",
				"error":   msg,
			})
		}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "appointment rescheduled",
	})
}

// GET /tenants/:tenant_id/barbers/:barber_id/appointments?start=...&end=...
// GetAppointmentsByBarber godoc
// @Summary      ดึงนัดหมายของช่างตัดผม
// @Description  คืนรายการ Appointment ของช่างตัดผมที่ระบุ ภายใน Tenant ที่กำหนด และช่วงเวลาเลือกได้ (RFC3339)
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id  path      uint      true   "รหัส Tenant"
// @Param        barber_id  path      uint      true   "รหัส Barber"
// @Param        start      query     string    false  "เวลาเริ่มต้นกรอง (RFC3339), เช่น 2025-05-30T09:00:00Z"
// @Param        end        query     string    false  "เวลาสิ้นสุดกรอง (RFC3339), เช่น 2025-05-30T17:00:00Z"
// @Success      200        {object}  map[string][]barberBookingModels.Appointment  "คืนค่า status success และ array ของ Appointment ใน key `data`"
// @Failure      400        {object}  map[string]string                            "Invalid tenant_id, barber_id หรือรูปแบบเวลาไม่ถูกต้อง"
// @Failure      404        {object}  map[string]string                            "Barber not found"
// @Failure      500        {object}  map[string]string                            "Internal Server Error"
// @Router       /none [get]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) GetAppointmentsByBarber(c *fiber.Ctx) error {
	// 1. Parse tenant_id (for consistency; not used in service but ensures valid URL)
	if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant_id",
		})
	}

	// 2. Parse barber_id
	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid barber_id",
		})
	}

	// 3. Parse optional start/end query params (RFC3339)
	var startPtr, endPtr *time.Time

	if s := c.Query("start", ""); s != "" {
		t, err := time.Parse(time.RFC3339, s)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid start format, expect RFC3339",
			})
		}
		startPtr = &t
	}
	if e := c.Query("end", ""); e != "" {
		t, err := time.Parse(time.RFC3339, e)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid end format, expect RFC3339",
			})
		}
		endPtr = &t
	}

	// 4. Call service
	appts, err := ctrl.Service.GetAppointmentsByBarber(c.Context(), barberID, startPtr, endPtr)
	if err != nil {
		// distinguish not-found vs other errors if service returns specific messages
		if strings.Contains(err.Error(), "barber with ID") {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": "Failed to fetch appointments",
			"error":   err.Error(),
		})
	}

	// 5. Return result
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   appts,
	})
}

var RolesCanManageAppointment = []coreModels.RoleName{
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
}

// DELETE /tenants/:tenant_id/appointments/:appointment_id
// DeleteAppointment godoc
// @Summary      ลบการนัดหมาย
// @Description  ลบ Appointment ตามรหัสที่ระบุ ภายใต้ Tenant ที่กำหนด (ต้องมีสิทธิ์ Tenant หรือ TenantAdmin)
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        tenant_id        path      uint     true  "รหัส Tenant"
// @Param        appointment_id   path      uint     true  "รหัส Appointment"
// @Success      200              {object}  map[string]string  "คืนค่า status success และข้อความยืนยันการลบ"
// @Failure      400              {object}  map[string]string  "Invalid tenant_id หรือ appointment_id"
// @Failure      403              {object}  map[string]string  "Permission denied"
// @Failure      404              {object}  map[string]string  "Appointment not found"
// @Failure      500              {object}  map[string]string  "Internal Server Error"
// @Router       /tenants/:tenant_id/appointments/:appointment_id [delete]
// @Security     ApiKeyAuth
func (ctrl *AppointmentController) DeleteAppointment(c *fiber.Ctx) error {
	// 1. Authorization: ตรวจสิทธิ์ก่อน
	roleStr, ok := c.Locals("role").(string)
	if !ok || !helperFunc.IsAuthorizedRole(roleStr, RolesCanManageAppointment) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"status":  "error",
			"message": "Permission denied",
		})
	}

	// 2. Parse tenant_id (validate format)
	if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid tenant_id",
		})
	}

	// 3. Parse appointment_id
	apptID, err := helperFunc.ParseUintParam(c, "appointment_id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid appointment_id",
		})
	}

	// 4. Call service
	err = ctrl.Service.DeleteAppointment(context.Background(), apptID)
	if err != nil {
		msg := err.Error()
		switch {
		case strings.Contains(msg, "not found"):
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"status":  "error",
				"message": msg,
			})
		default:
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"status":  "error",
				"message": fmt.Sprintf("Failed to delete appointment: %s", msg),
			})
		}
	}

	// 5. Success
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  "success",
		"message": "appointment deleted",
	})
}
