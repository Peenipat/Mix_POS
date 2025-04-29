package controllers

import (
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	authDto "myapp/dto/auth"
	"myapp/models"
	"myapp/services"
)

// logService เก็บ instance ของ SystemLogService
var logService services.SystemLogService

// InitSystemLogHandler ต้องเรียกก่อนผูก routes เพื่อ inject service เพื่อ ให้แยก controller และ service ออกจากกัน
func InitSystemLogHandler(svc services.SystemLogService) {
	logService = svc
}

// CreateLog handles POST /admin/system_logs
func CreateLog(ctx *fiber.Ctx) error {
	//ผูกตัวแปรเข้ากับ request โดย check ด้วยว่าข้อมูลที่เข้ามามี type ตรงตามที่กำหนดไว้หรือเปล่า
	var req authDto.CreateLogRequest
	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// Map request to model Sytemlogs
	entry := &models.SystemLog{
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
	res := authDto.CreateLogResponse{
		LogID:     entry.LogID,
		CreatedAt: entry.CreatedAt,
	}
	return ctx.Status(fiber.StatusCreated).JSON(res)
}

// GetSystemLogs handles GET /admin/system_logs
func GetSystemLogs(ctx *fiber.Ctx) error {
	page, _ := strconv.Atoi(ctx.Query("page", "1"))
	limit, _ := strconv.Atoi(ctx.Query("limit", "20"))
	filter := services.LogFilter{Page: page, Limit: limit}

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

// GetSystemLogByID handles GET /admin/system_logs/:log_id
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
