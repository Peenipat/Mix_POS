package barberBookingController

import (
	"fmt"
	"net/http"
	"strconv"
	"errors"

	"github.com/gofiber/fiber/v2"
	barberBookingPort "myapp/modules/barberbooking/port"
)

type WorkingDayOverrideController struct {
	Service barberBookingPort.IWorkingDayOverrideService
}

func NewWorkingDayOverrideController(service barberBookingPort.IWorkingDayOverrideService) *WorkingDayOverrideController {
	return &WorkingDayOverrideController{
		Service: service,
	}
}

// Create godoc
// @Summary      สร้างวันเปิด-ปิดเฉพาะกิจ
// @Description  ใช้สำหรับเพิ่มวันเปิดหรือปิดเฉพาะกิจของสาขา เช่น วันหยุดประจำปีหรือวันเปิดพิเศษ
// @Tags         WorkingDayOverride
// @Accept       json
// @Produce      json
// @Param        body body barberBookingPort.WorkingDayOverrideInput true "ข้อมูลวันที่และเวลาที่ต้องการ override"
// @Success      201  {object} map[string]interface{} "สร้างสำเร็จ"
// @Failure      400  {object} map[string]interface{} "กรณี input ไม่ถูกต้อง"
// @Failure      500  {object} map[string]interface{} "กรณีสร้างไม่สำเร็จหรือเกิด error ภายใน"
// @Router       /working-day-overrides [post]
// @Security     ApiKeyAuth
func (ctrl *WorkingDayOverrideController) Create(c *fiber.Ctx) error {
	var input barberBookingPort.WorkingDayOverrideInput

	// Bind JSON
	if err := c.BodyParser(&input); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input format",
		})
	}

	// Optional: validate input here ifใช้ validator

	// Call service
	override, err := ctrl.Service.Create(c.Context(), input)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "สร้างวันเปิด/ปิดสำเร็จ",
		"data":    override,
	})
}

// Update godoc
// @Summary      แก้ไขวันเปิด-ปิดเฉพาะกิจ
// @Description  ใช้สำหรับแก้ไขวัน override ของสาขา เช่น เปลี่ยนเวลาเปิด-ปิด หรือเปลี่ยนวัน
// @Tags         WorkingDayOverride
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "รหัส override ที่ต้องการแก้ไข"
// @Param        body body      barberBookingPort.WorkingDayOverrideInput true "ข้อมูลที่ต้องการอัปเดต"
// @Success      200  {object}  map[string]interface{} "อัปเดตสำเร็จ"
// @Failure      400  {object}  map[string]interface{} "ข้อมูลไม่ถูกต้อง"
// @Failure      404  {object}  map[string]interface{} "ไม่พบข้อมูล"
// @Failure      500  {object}  map[string]interface{} "เกิดข้อผิดพลาดภายใน"
// @Router       /working-day-overrides/{id} [put]
// @Security     ApiKeyAuth
func (ctrl *WorkingDayOverrideController) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID parameter",
		})
	}

	var input barberBookingPort.WorkingDayOverrideInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if err := ctrl.Service.Update(c.Context(), uint(id), input); err != nil {
		if err.Error() == fmt.Sprintf("override with ID %d not found", id) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "อัปเดตข้อมูลเรียบร้อยแล้ว",
	})
}


// @Summary        ดึงข้อมูล override ตาม ID
// @Description    ใช้ดึงข้อมูลวันเปิด-ปิดเฉพาะกิจตาม ID
// @Tags           WorkingDayOverride
// @Accept         json
// @Produce        json
// @Param          id path int true "WorkingDayOverride ID"
// @Success        200 {object} barberBookingModels.WorkingDayOverride
// @Failure        400 {object} map[string]string
// @Failure        404 {object} map[string]string
// @Router         /working-day-overrides/{id} [get]
// @Security       ApiKeyAuth
func (ctrl *WorkingDayOverrideController) GetByID(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID format",
		})
	}

	result, err := ctrl.Service.GetByID(c.Context(), uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return c.Status(fiber.StatusOK).JSON(result)
}

// @Summary      ลบวันเปิด-ปิดร้านเฉพาะวัน
// @Description  ใช้ลบ override เฉพาะวันจาก branch ที่ระบุ
// @Tags         WorkingDayOverride
// @Param        id   path      int  true  "WorkingDayOverride ID"
// @Success      204  {string}  string  "No Content"
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /working-day-overrides/{id} [delete]
// @Security     ApiKeyAuth
func (c *WorkingDayOverrideController) DeleteWorkingDayOverride(ctx *fiber.Ctx) error {
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil || id <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "invalid id parameter",
		})
	}

	if err := c.Service.Delete(ctx.Context(), uint(id)); err != nil {
		if errors.Is(err, fiber.ErrNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"error": err.Error(),
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

