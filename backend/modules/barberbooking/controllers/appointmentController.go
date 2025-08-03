package barberBookingController

import (
	"context"
	"fmt"
	"strconv"

	"strings"
	"time"

	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"
	coreModels "myapp/modules/core/models"

	"github.com/gofiber/fiber/v2"
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
		BranchID   uint                             `json:"branch_id"`
		ServiceID  uint                             `json:"service_id"`
		BarberID   uint                             `json:"barber_id,omitempty"`
		CustomerID uint                             `json:"customer_id"`
		StartTime  string                           `json:"start_time"`
		Notes      string                           `json:"notes,omitempty"`
		Customer   *barberBookingPort.CustomerInput `json:"customer,omitempty"`
	}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid request body",
		})
	}

	// 3. Validate required fields
	if payload.CustomerID == 0 {
		if payload.Customer == nil || payload.Customer.Name == "" || payload.Customer.Phone == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Guest appointment requires customer name and phone",
			})
		}
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
	// appt := &barberBookingModels.Appointment{
	// 	TenantID:   tenantID,
	// 	BranchID:   payload.BranchID,
	// 	ServiceID:  payload.ServiceID,
	// 	BarberID:   payload.BarberID,
	// 	CustomerID: payload.CustomerID,
	// 	StartTime:  startTime,
	// 	Notes:      payload.Notes,
	// }

	appt := &barberBookingModels.Appointment{
		TenantID:   tenantID,
		BranchID:   payload.BranchID,
		ServiceID:  payload.ServiceID,
		BarberID:   payload.BarberID,
		CustomerID: payload.CustomerID,
		StartTime:  startTime,
		Notes:      payload.Notes,
	}

	// ถ้า guest → แนบข้อมูล guest ไปให้ service ใช้สร้าง customer
	if payload.CustomerID == 0 && payload.Customer != nil {
		appt.Customer = &barberBookingModels.Customer{
			Name:  payload.Customer.Name,
			Phone: payload.Customer.Phone,
		}
	}

	// 6. Call service
	createdDTO, err := ctrl.Service.CreateAppointment(c.Context(), appt)
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

// @Summary      ดึงการนัดหมายทั้งหมดของช่างในสาขา
// @Description  คืนรายการการจองคิวของช่างทั้งหมดในสาขาที่กำหนด โดยสามารถกรองช่วงเวลา, สถานะ, ลูกค้า, และใช้ pagination ได้
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        branch_id       path      uint     true   "รหัสสาขา (Branch ID)"
// @Param        tenant_id       path     uint     true   "รหัสผู้เช่า (Tenant ID)"
// @Param        start           query     string   false  "วันที่เริ่มต้น (รูปแบบ: yyyy-MM-dd)"
// @Param        end             query     string   false  "วันที่สิ้นสุด (รูปแบบ: yyyy-MM-dd)"
// @Param        filter          query     string   false  "ประเภทการกรองช่วงเวลา: week (สัปดาห์นี้), month (เดือนนี้)"
// @Param        exclude_status  query     string   false  "สถานะที่ไม่ต้องการให้แสดง เช่น CANCELLED,NO_SHOW (คั่นด้วย ,)"
// @Param        status          query     string   false  "สถานะการนัดหมาย เช่น WAITING, COMPLETED"
// @Param        barber_id       query     uint     false  "รหัสช่าง"
// @Param        service_id      query     uint     false  "รหัสบริการ"
// @Param        created_date    query     string   false  "วันที่สร้างการนัดหมาย (yyyy-MM-dd)"
// @Param        cus_name        query     string   false  "ชื่อลูกค้า"
// @Param        phone           query     string   false  "เบอร์โทรลูกค้า"
// @Param        page            query     int      false  "หน้าปัจจุบัน (เริ่มที่ 1)"
// @Param        limit           query     int      false  "จำนวนรายการต่อหน้า"
// @Success      200             {object}  map[string]interface{}  "คืน status success และรายการการนัดหมาย"
// @Failure      400             {object}  map[string]string        "พารามิเตอร์ไม่ถูกต้อง เช่น วันที่ผิดรูปแบบ"
// @Failure      404             {object}  map[string]string        "ไม่พบข้อมูล"
// @Failure      500             {object}  map[string]string        "ข้อผิดพลาดภายในเซิร์ฟเวอร์"
// @Router       /tenants/{tenant_id}/branches/{branch_id}/appointments [get]
func (ctrl *AppointmentController) GetAppointments(c *fiber.Ctx) error {
	// ── branch_id
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil || branchID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch_id",
		})
	}

	// ── tenant_id (บังคับ)
	tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
	if err != nil || tenantID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Missing tenant_id",
		})
	}

	// ── ตัวกรองเพิ่มเติม
	search := c.Query("search")
	page, _ := strconv.Atoi(c.Query("page", "1"))
	limit, _ := strconv.Atoi(c.Query("limit", "10"))

	// ── filter: status (support หลายค่า)
	var statuses []barberBookingModels.AppointmentStatus
	rawStatuses := c.Query("status") // ex. "BOOKED,COMPLETED"
	if rawStatuses != "" {
		for _, s := range strings.Split(rawStatuses, ",") {
			status := barberBookingModels.AppointmentStatus(strings.ToUpper(strings.TrimSpace(s)))
			statuses = append(statuses, status)
		}
	}

	// ── filter: barber_id
	var barberID *uint
	if bid, err := helperFunc.ParseUintQuery(c, "barber_id"); err == nil && bid > 0 {
		barberID = &bid
	}

	// ── filter: service_id
	var serviceID *uint
	if sid, err := helperFunc.ParseUintQuery(c, "service_id"); err == nil && sid > 0 {
		serviceID = &sid
	}

	// ── created_date range
	var createdStart, createdEnd *time.Time
	layout := "2006-01-02"

	if val := c.Query("created_start"); val != "" {
		t, err := time.Parse(layout, val)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid created_start format (yyyy-MM-dd)",
			})
		}
		createdStart = &t
	}
	if val := c.Query("created_end"); val != "" {
		t, err := time.Parse(layout, val)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid created_end format (yyyy-MM-dd)",
			})
		}
		t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond)
		createdEnd = &t
	}

	// ── สร้าง filter object
	filter := barberBookingPort.GetAppointmentsFilter{
		TenantID:  &tenantID,
		BranchID:  &branchID,
		Search:    search,
		Statuses:  statuses,
		BarberID:  barberID,
		ServiceID: serviceID,
		StartDate: createdStart,
		EndDate:   createdEnd,
		Page:      page,
		Limit:     limit,
	}

	// ── เรียก service ใหม่
	data, total, err := ctrl.Service.GetAppointments(c.Context(), filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	if data == nil {
		data = []barberBookingPort.AppointmentBrief{}
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   data,
		"meta": fiber.Map{
			"pagination": fiber.Map{
				"total": total,
				"page":  page,
				"limit": limit,
			},
		},
	})
}

// GetAppointmentsByBarber ดึงการนัดหมายทั้งหมดของช่างในช่วงเวลาหรือสถานะที่กำหนด
// @Summary      ดึงการนัดหมายทั้งหมดของช่าง
// @Description  คืนรายการการจองคิวของช่างในระบบ โดยสามารถกรองตามช่วงเวลา, สถานะ หรือ preset mode (today, week, past)
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        barber_id   path      uint     true   "รหัสช่าง (Barber ID)"
// @Param        start       query     string   false  "เวลาที่เริ่มต้นช่วง (รูปแบบ: yyyy-MM-dd)"
// @Param        end         query     string   false  "เวลาที่สิ้นสุดช่วง (รูปแบบ: yyyy-MM-dd)"
// @Param        status      query     string   false  "กรองสถานะ เช่น CONFIRMED,COMPLETED (คั่นด้วย comma)"
// @Param        mode        query     string   false  "โหมด preset ช่วงเวลา เช่น today, week, past"
// @Success      200         {object}  map[string]interface{}  "คืน status success และรายการการนัดหมาย"
// @Failure      400         {object}  map[string]string
// @Failure      404         {object}  map[string]string
// @Failure      500         {object}  map[string]string
// @Router       /barbers/{barber_id}/appointments [get]
func (ctrl *AppointmentController) GetAppointmentsByBarber(c *fiber.Ctx) error {
	barberID, err := helperFunc.ParseUintParam(c, "barber_id")
	if err != nil || barberID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid barber_id",
		})
	}

	layout := "2006-01-02"
	var startTime *time.Time
	var endTime *time.Time
	if startStr := c.Query("start"); startStr != "" {
		t, err := time.Parse(layout, startStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid start_time format (yyyy-MM-dd)",
			})
		}
		startTime = &t
	}
	if endStr := c.Query("end"); endStr != "" {
		t, err := time.Parse(layout, endStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid end_time format (yyyy-MM-dd)",
			})
		}
		t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond)
		endTime = &t
	}

	// timeMode = today, week, past
	timeMode := c.Query("mode", "")

	// status=CONFIRMED,COMPLETED
	statusParam := c.Query("status", "")
	var statusList []barberBookingModels.AppointmentStatus
	if statusParam != "" {
		parts := strings.Split(statusParam, ",")
		for _, s := range parts {
			statusList = append(statusList, barberBookingModels.AppointmentStatus(strings.ToUpper(strings.TrimSpace(s))))
		}
	}

	filter := barberBookingPort.AppointmentFilter{
		Start:    startTime,
		End:      endTime,
		Status:   statusList,
		TimeMode: timeMode,
	}

	appts, err := ctrl.Service.GetAppointmentsByBarber(c.Context(), barberID, filter)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   appts,
	})
}

// @Summary      ดึงการนัดหมายทั้งหมดของช่างในสาขา
// @Description  คืนรายการการจองคิวของช่างทั้งหมดในสาขาที่กำหนด โดยสามารถกรองตามช่วงเวลาที่กำหนดได้ เช่น ช่วงเวลาเริ่มต้น/สิ้นสุด, หรือช่วงสัปดาห์นี้/เดือนนี้
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        branch_id   path      uint     true   "รหัสสาขา (Branch ID)"
// @Param        start       query     string   false  "เวลาที่เริ่มต้นช่วง (รูปแบบ: yyyy-MM-dd) เช่น 2025-07-15"
// @Param        end         query     string   false  "เวลาที่สิ้นสุดช่วง (รูปแบบ: yyyy-MM-dd) เช่น 2025-07-20"
// @Param        filter      query     string   false  "ประเภทการกรองเวลา: week (สัปดาห์นี้), month (เดือนนี้)"
// @Param        exclude_status query     string   false  "รายการสถานะที่ไม่ต้องการให้แสดง เช่น CANCELLED,NO_SHOW (คั่นด้วย ,)"
// @Success      200         {object}  map[string]interface{}  "คืน status success และรายการการนัดหมาย"
// @Failure      400         {object}  map[string]string        "กรณีพารามิเตอร์ไม่ถูกต้อง เช่น วันที่ผิดรูปแบบ"
// @Failure      404         {object}  map[string]string        "ไม่พบข้อมูล"
// @Failure      500         {object}  map[string]string        "เกิดข้อผิดพลาดภายในเซิร์ฟเวอร์"
// @Router       /branches/{branch_id}/appointments [get]
func (ctrl *AppointmentController) GetAppointmentsByBranch(c *fiber.Ctx) error {
	branchID, err := helperFunc.ParseUintParam(c, "branch_id")
	if err != nil || branchID == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "Invalid branch_id",
		})
	}

	filterType := c.Query("filter") 

	var startTime *time.Time
	var endTime *time.Time
	layout := "2006-01-02"

	if startStr := c.Query("start"); startStr != "" {
		t, err := time.Parse(layout, startStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid start_time format (yyyy-MM-dd)",
			})
		}
		startTime = &t
	}

	if endStr := c.Query("end"); endStr != "" {
		t, err := time.Parse(layout, endStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid end_time format (yyyy-MM-dd)",
			})
		}
		t = t.Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond)
		endTime = &t
	}

	rawStatuses := c.Query("exclude_status") // ตัวอย่าง: "CANCELLED,NO_SHOW"
	var excludeStatuses []barberBookingModels.AppointmentStatus
	if rawStatuses != "" {
		for _, s := range strings.Split(rawStatuses, ",") {
			status := barberBookingModels.AppointmentStatus(strings.ToUpper(strings.TrimSpace(s)))
			excludeStatuses = append(excludeStatuses, status)
		}
	}

	// เรียก service พร้อมส่ง filterType เพิ่ม
	appts, err := ctrl.Service.GetAppointmentsByBranch(c.Context(), branchID, startTime, endTime, filterType,excludeStatuses)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   appts,
	})
}


// GetAppointmentsByPhone ดึงรายการการจองของลูกค้าจากเบอร์โทร
// @Summary      ดึงการนัดหมายด้วยเบอร์โทร
// @Description  คืนรายการการจองทั้งหมดของลูกค้าที่ใช้เบอร์โทรนั้นในระบบ
// @Tags         Appointment
// @Accept       json
// @Produce      json
// @Param        phone     query     string   true   "เบอร์โทรลูกค้า (Customer Phone)"
// @Success      200       {object}  map[string]interface{}  "คืน status success และรายการการนัดหมาย"
// @Failure      400       {object}  map[string]string
// @Failure      500       {object}  map[string]string
// @Router       /appointments/by-phone [get]
func (ctrl *AppointmentController) GetAppointmentsByPhone(c *fiber.Ctx) error {
	phone := strings.TrimSpace(c.Query("phone"))
	if phone == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  "error",
			"message": "กรุณาระบุเบอร์โทรศัพท์ (phone)",
		})
	}

	appts, err := ctrl.Service.GetAppointmentsByPhone(c.Context(), phone)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  "error",
			"message": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status": "success",
		"data":   appts,
	})
}
