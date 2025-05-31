package Core_controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	Core_authDto "myapp/modules/core/dto/auth"
	"myapp/modules/core/models"
	"myapp/modules/core/services"
)

// logService เก็บ instance ของ SystemLogService
var logService coreServices.SystemLogService

// InitSystemLogHandler ต้องเรียกก่อนผูก routes เพื่อ inject service เพื่อ ให้แยก controller และ service ออกจากกัน
func InitSystemLogHandler(svc coreServices.SystemLogService) {
	logService = svc
}
// CreateLog godoc
// @Summary      สร้าง Log Entry
// @Description  รับข้อมูล Log (action, resource, status, HTTPMethod, endpoint, optional user_id, branch_id, details) แล้วบันทึกลงระบบ
// @Tags         Log
// @Accept       json
// @Produce      json
// @Param        body  body      Core_authDto.CreateLogRequest  true  "ข้อมูลสำหรับสร้าง Log Entry"
// @Success      201   {object}  Core_authDto.CreateLogResponse   "คืนค่า LogID และ CreatedAt ของ Log ที่สร้าง"
// @Failure      400   {object}  map[string]string               "Invalid input หรือ invalid details format"
// @Failure      500   {object}  map[string]string               "เกิดข้อผิดพลาดระหว่างการสร้าง Log"
// @Router       /none [post]
// @Security     ApiKeyAuth
func CreateLog(ctx *fiber.Ctx) error {
	//ผูกตัวแปรเข้ากับ request โดย check ด้วยว่าข้อมูลที่เข้ามามี type ตรงตามที่กำหนดไว้หรือเปล่า
	var req Core_authDto.CreateLogRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Map request to model Sytemlogs
	entry := &coreModels.SystemLog{
		CreatedAt:  time.Now(),
		Action:     req.Action,
		Resource:   "", //กำหนดไว้รอรับจาก ตัวอื่น
		Status:     req.Status,
		HTTPMethod: req.HTTPMethod,
		Endpoint:   req.Endpoint,
	}

	// รับค่า Resource
	if req.Resource != nil {
		entry.Resource = *req.Resource
	}
	// แปลงค่า user id จาก string เป็น unit
	if req.UserID != nil {
		if id, err := strconv.ParseUint(*req.UserID, 10, 32); err == nil {
			u := uint(id)
			entry.UserID = &u
		}
	}
	if req.BranchID != nil {
		if bid, err := strconv.ParseUint(*req.BranchID, 10, 32); err == nil {
			b := uint(bid)
			entry.BranchID = &b
		}
	}
	// แปลง Details เป็น JSON bytes
	if req.Details != nil {
		raw, err := json.Marshal(req.Details)
		if err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid details format"})
		}
		entry.Details = raw
	}

	// Save log entry
	if err := logService.Create(ctx.Context(), entry); err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "cannot create log"})
	}

	//เตรียมของเพื่อส่งกลับ
	res := Core_authDto.CreateLogResponse{
		LogID:     entry.LogID,
		CreatedAt: entry.CreatedAt,
	}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}

// GetSystemLogs godoc
// @Summary      ดึงประวัติระบบ (System Logs)
// @Description  ดึงรายการ System Logs พร้อม pagination และ filter ตาม action, endpoint, status, ช่วงเวลา `from`–`to` (RFC3339)
// @Tags         Log
// @Produce      json
// @Param        page      query     int     false  "หน้า (default = 1)"             default(1)
// @Param        limit     query     int     false  "จำนวนต่อหน้า (default = 20)"   default(20)
// @Param        action    query     string  false  "กรองตาม action"
// @Param        endpoint  query     string  false  "กรองตาม endpoint"
// @Param        status    query     string  false  "กรองตาม status"
// @Param        from      query     string  false  "วันที่เริ่มต้น (RFC3339)"
// @Param        to        query     string  false  "วันที่สิ้นสุด (RFC3339)"
// @Success      200       {object}  map[string]interface{}  "คืนค่า total และ logs[]"
// @Failure      500       {object}  map[string]string       "เกิดข้อผิดพลาดระหว่างดึง System Logs"
// @Router       /logs [get]
// @Security     ApiKeyAuth
func GetSystemLogs(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	filter := coreServices.LogFilter{Page: page, Limit: limit}

	if v := ctx.Query("action"); v != "" {
		filter.Action = &v
	}
	if v := ctx.Query("endpoint"); v != "" {
		filter.Endpoint = &v
	}
	if v := ctx.Query("status"); v != "" {
		filter.Status = &v
	}
	//กำหนดช่วงเวลา
	if v := ctx.Query("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.From = &t
		}
	}
	//กำหนดช่วงเวลา
	if v := ctx.Query("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.To = &t
		}
	}
	//เรียก service มาดึง log ทั้งหมด และนับค่าที่ได้
	logs, total, err := logService.Query(ctx.Context(), filter)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	//ส่งของคืน
	return ctx.JSON(fiber.Map{"total": total, "logs": logs})
}

// GetSystemLogByID godoc
// @Summary      ดึง System Log ตาม ID
// @Description  ดึงข้อมูล System Log รายการเดียวตามรหัส `log_id`
// @Tags         Log
// @Produce      json
// @Param        log_id  path      int  true  "รหัส Log"
// @Success      200     {object}  coreModels.SystemLog  "คืนค่า SystemLog object"
// @Failure      400     {object}  map[string]string     "invalid log_id"
// @Failure      404     {object}  map[string]string     "log not found"
// @Failure      500     {object}  map[string]string     "เกิดข้อผิดพลาดระหว่างดึง Log"
// @Router       /none [get]
// @Security     ApiKeyAuth
func GetSystemLogByID(ctx *fiber.Ctx) error {
	param := ctx.Params("log_id")
	id64, err := strconv.ParseUint(param, 10, 32)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid log_id"})
	}
	//ดึง log ตาม id
	entry, err := logService.GetByID(ctx.Context(), uint(id64))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "log not found"})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return ctx.JSON(entry)
}
