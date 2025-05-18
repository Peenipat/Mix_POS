package barberBookingController

import (
    "time"
	"context"
	"strings"
	"strconv"
    "fmt"

    "github.com/gofiber/fiber/v2"
    barberBookingPort "myapp/modules/barberbooking/port"
    barberBookingModels "myapp/modules/barberbooking/models"
    coreModels "myapp/modules/core/models"
	barberBookingDto "myapp/modules/barberbooking/dto"
    helperFunc "myapp/modules/barberbooking"

)

// AppointmentController handles endpoints related to appointments
type AppointmentController struct {
    Service barberBookingPort.IAppointment
}

// NewAppointmentController creates a new controller
func NewAppointmentController(service barberBookingPort.IAppointment) *AppointmentController {
    return &AppointmentController{Service: service}
}

// CheckBarberAvailability handles GET /tenants/:tenant_id/barbers/:barber_id/availability?start=&end=
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

// CreateAppointment handles POST /tenants/:tenant_id/appointments
func (ctrl *AppointmentController) CreateAppointment(c *fiber.Ctx) error {
    // 1. Parse tenant_id
    tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).
            JSON(fiber.Map{"status": "error", "message": "Invalid tenant_id"})
    }
    // 2. Parse body
    var payload struct {
        BranchID    uint   `json:"branch_id"`
        ServiceID   uint   `json:"service_id"`
        BarberID    *uint  `json:"barber_id,omitempty"`
        CustomerID  uint   `json:"customer_id"`
        StartTime   string `json:"start_time"`
        Notes       string `json:"notes,omitempty"`
    }
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).
            JSON(fiber.Map{"status": "error", "message": "Invalid request body"})
    }
    // 3. Validate required
    if payload.BranchID == 0 || payload.ServiceID == 0 || payload.CustomerID == 0 || payload.StartTime == "" {
        return c.Status(fiber.StatusBadRequest).
            JSON(fiber.Map{"status": "error", "message": "Missing required fields"})
    }
    // 4. Parse start_time
    startTime, err := time.Parse(time.RFC3339, payload.StartTime)
    if err != nil {
        return c.Status(fiber.StatusBadRequest).
            JSON(fiber.Map{"status": "error", "message": "Invalid start_time format. Expect RFC3339"})
    }
    // 5. Build model
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
    created, err := ctrl.Service.CreateAppointment(c.Context(), appt)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).
            JSON(fiber.Map{"status": "error", "message": err.Error()})
    }
    // 7. Return result
    return c.Status(fiber.StatusCreated).
        JSON(fiber.Map{"status": "success", "data": created})
}

// GetAvailableBarbers handles GET /tenants/:tenant_id/branches/:branch_id/available-barbers?start=&end=
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
func (ctrl *AppointmentController) ListAppointments(c *fiber.Ctx) error {
    // 1. Parse tenant_id
    tenantID, err := helperFunc.ParseUintParam(c, "tenant_id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid tenant_id"})
    }

    // 2. Build filter
    var f barberBookingDto.AppointmentFilter
    f.TenantID = tenantID

    if qs := c.Query("branch_id",""); qs != "" {
        if v, err := strconv.ParseUint(qs,10,64); err!=nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid branch_id"})
        } else {
            u := uint(v); f.BranchID = &u
        }
    }
    if qs := c.Query("barber_id",""); qs != "" {
        if v, err := strconv.ParseUint(qs,10,64); err!=nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid barber_id"})
        } else {
            u := uint(v); f.BarberID = &u
        }
    }
    if qs := c.Query("customer_id",""); qs != "" {
        if v, err := strconv.ParseUint(qs,10,64); err!=nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid customer_id"})
        } else {
            u := uint(v); f.CustomerID = &u
        }
    }
    if qs := c.Query("status",""); qs != "" {
        s := barberBookingModels.AppointmentStatus(qs)
        f.Status = &s
    }
    if qs := c.Query("start_date",""); qs != "" {
        t, err := time.Parse(time.RFC3339, qs)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid start_date format. Expect RFC3339"})
        }
        f.StartDate = &t
    }
    if qs := c.Query("end_date",""); qs != "" {
        t, err := time.Parse(time.RFC3339, qs)
        if err != nil {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid end_date format. Expect RFC3339"})
        }
        f.EndDate = &t
    }
    if qs := c.Query("limit",""); qs != "" {
        if v, err := strconv.Atoi(qs); err!=nil  || v < 0 {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid limit"})
        } else {
            f.Limit = &v
        }
    }
    if qs := c.Query("offset",""); qs != "" {
        if v, err := strconv.Atoi(qs); err!=nil || v < 0  {
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid offset"})
        } else {
            f.Offset = &v
        }
    }
    if qs := c.Query("sort_by",""); qs != "" {
        f.SortBy = &qs
    }

    // 3. Call service
    apps, err := ctrl.Service.ListAppointments(context.Background(), f)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":"error", "message":"Failed to list appointments", "error":err.Error(),
        })
    }

    // 4. Return
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":"success",
        "data": apps,
    })
}

type CancelRequest struct {
    ActorUserID     *uint `json:"actor_user_id,omitempty"`
    ActorCustomerID *uint `json:"actor_customer_id,omitempty"`
}
// POST /tenants/:tenant_id/appointments/:appointment_id/cancel
func (ctrl *AppointmentController) CancelAppointment(c *fiber.Ctx) error {
    // 1. Parse path params
    if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid tenant_id"})
    }
    apptID, err := helperFunc.ParseUintParam(c, "appointment_id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid appointment_id"})
    }

    // 2. Parse body
    var req CancelRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":"Invalid JSON body"})
    }

    // 3. Validate that exactly one actor is provided
    if (req.ActorUserID == nil && req.ActorCustomerID == nil) ||
       (req.ActorUserID != nil && req.ActorCustomerID != nil) {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":"error",
            "message":"Either actor_user_id or actor_customer_id must be provided, but not both",
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
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status":"error","message":msg})
        case strings.Contains(msg, "cannot be cancelled"):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":msg})
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"status":"error","message":"Failed to cancel appointment","error":msg})
        }
    }

    return c.Status(fiber.StatusOK).JSON(fiber.Map{"status":"success","message":"appointment cancelled"})
}


type RescheduleRequest struct {
    NewStartTime     string `json:"new_start_time"`
    ActorUserID      *uint  `json:"actor_user_id,omitempty"`
    ActorCustomerID  *uint  `json:"actor_customer_id,omitempty"`
}
// POST /tenants/:tenant_id/appointments/:appointment_id/reschedule
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
	coreModels.RoleNameSaaSSuperAdmin,
	coreModels.RoleNameTenant,
	coreModels.RoleNameTenantAdmin,
}

// DELETE /tenants/:tenant_id/appointments/:appointment_id
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