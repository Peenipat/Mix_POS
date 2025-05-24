package barberBookingController

import (
	// "fmt"
	"net/http"
	"strings"
	"time"

	helperFunc "myapp/modules/barberbooking"
	barberBookingModels "myapp/modules/barberbooking/models"
	barberBookingPort "myapp/modules/barberbooking/port"

	"github.com/gofiber/fiber/v2"
)

type AppointmentReviewController struct {
    ReviewService barberBookingPort.IAppointmentReview
}

func NewAppointmentReviewController(svc barberBookingPort.IAppointmentReview) *AppointmentReviewController {
    return &AppointmentReviewController{ReviewService: svc}
}

// GET /tenants/:tenant_id/reviews/:review_id
func (ctrl *AppointmentReviewController) GetReviewByID(c *fiber.Ctx) error {

    // 1. Parse tenant_id (เพื่อ consistency)
    if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant_id",
        })
    }

    // 2. Parse review_id และเช็คว่าไม่เป็น 0
    rid, err := helperFunc.ParseUintParam(c, "review_id")
    if err != nil || rid == 0 {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid review_id",
        })
    }

    // 3. Call service
    rev, err := ctrl.ReviewService.GetByID(c.Context(), rid)
    if err != nil {
        msg := err.Error()
        switch {
        case strings.Contains(msg, "invalid review id"):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
                "status":  "error",
                "message": msg,
            })
        case strings.Contains(msg, "not found"):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
                "status":  "error",
                "message": msg,
            })
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to fetch review",
                "error":   msg,
            })
        }
    }

    // 4. Return success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   rev,
    })
}

// POST /tenants/:tenant_id/appointments/:appointment_id/reviews
func (ctrl *AppointmentReviewController) CreateReview(c *fiber.Ctx) error {
    // 1. Parse tenant_id
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

    // 3. Bind JSON body
    var payload struct {
        CustomerID *uint  `json:"customer_id,omitempty"`
        Rating     int    `json:"rating"`
        Comment    string `json:"comment,omitempty"`
    }
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid JSON body",
        })
    }

    // 4. Build model
    review := &barberBookingModels.AppointmentReview{
        AppointmentID: apptID,
        CustomerID:    payload.CustomerID,
        Rating:        payload.Rating,
        Comment:       payload.Comment,
    }

    // 5. Call service
    created, err := ctrl.ReviewService.CreateReview(c.Context(), review)
    if err != nil {
        msg := err.Error()
        switch {
        case strings.Contains(msg, "invalid review input"):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":msg})
        case strings.Contains(msg, "appointment lookup failed"):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status":"error","message":strings.TrimPrefix(msg, "appointment lookup failed: ")})
        case strings.Contains(msg, "cannot review appointment"):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":msg})
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":"error","message":"Failed to create review","error":msg,
            })
        }
    }

    // 6. Return created
    return c.Status(fiber.StatusCreated).JSON(fiber.Map{
        "status":"success","data": created,
    })
}

// PUT /tenants/:tenant_id/appointments/:appointment_id/reviews/:review_id
func (ctrl *AppointmentReviewController) UpdateReview(c *fiber.Ctx) error {
    // 1. Validate tenant_id for path consistency
    if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant_id",
        })
    }

    // 2. Parse review_id
    rid, err := helperFunc.ParseUintParam(c, "review_id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid review_id",
        })
    }

    // 3. Bind JSON body
    var payload struct {
        Rating  int    `json:"rating"`
        Comment string `json:"comment,omitempty"`
    }
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid JSON body",
        })
    }

    // 4. Build review model to pass to service
    now := time.Now()
    input := &barberBookingModels.AppointmentReview{
        Rating:    payload.Rating,
        Comment:   payload.Comment,
        UpdatedAt: now, // service will overwrite with its own timestamp if needed
    }

    // 5. Call service
    updated, svcErr := ctrl.ReviewService.UpdateReview(c.Context(), rid, input)
    if svcErr != nil {
        msg := svcErr.Error()
        switch {
        case strings.Contains(msg, "invalid rating"):
            return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"status":"error","message":msg})
        case strings.Contains(msg, "not found"):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status":"error","message":msg})
        default:
            return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to update review",
                "error":   msg,
            })
        }
    }

    // 6. Return updated review
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   updated,
    })
}

type DeleteReviewRequest struct {
    ActorCustomerID uint `json:"actor_customer_id"`
}

// DELETE /tenants/:tenant_id/appointments/:appointment_id/reviews/:review_id
func (ctrl *AppointmentReviewController) DeleteReview(c *fiber.Ctx) error {
    // 1. Parse tenant_id (for URL consistency)
    if _, err := helperFunc.ParseUintParam(c, "tenant_id"); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid tenant_id",
        })
    }

    // 2. Parse appointment_id (not used in service, but keeps URL consistent)
    if _, err := helperFunc.ParseUintParam(c, "appointment_id"); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid appointment_id",
        })
    }

    // 3. Parse review_id
    rid, err := helperFunc.ParseUintParam(c, "review_id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid review_id",
        })
    }

    // 4. Parse body to get actorCustomerID
    var req DeleteReviewRequest
    if err := c.BodyParser(&req); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "status":  "error",
            "message": "Invalid JSON body",
        })
    }

    // 5. Call service
    err = ctrl.ReviewService.DeleteReview(c.UserContext(), rid, req.ActorCustomerID,"Customer")
    if err != nil {
        msg := err.Error()
        switch {
        case strings.Contains(msg, "not found"):
            return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"status":"error","message":msg})
        case strings.Contains(msg, "not authorized"):
            return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"status":"error","message":msg})
        default:
            return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
                "status":  "error",
                "message": "Failed to delete review",
                "error":   msg,
            })
        }
    }

    // 6. Success
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  "success",
        "message": "Review deleted",
    })
}

// GET /tenants/:tenant_id/barbers/:barber_id/average-rating
func (ctrl *AppointmentReviewController) GetAverageRatingByBarber(c *fiber.Ctx) error {
    // 1. Validate tenant_id (for URL consistency)
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

    // 3. Call service
    avg, svcErr := ctrl.ReviewService.GetAverageRatingByBarber(c.Context(), barberID)
    if svcErr != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "status":  "error",
            "message": "Failed to compute average rating",
            "error":   svcErr.Error(),
        })
    }

    // 4. Return result
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status": "success",
        "data":   avg, // float64
    })
}