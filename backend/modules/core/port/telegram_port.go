package corePort

import (
	"github.com/gofiber/fiber/v2"
)

type ITelegramService interface {
	ProcessWebhook(c *fiber.Ctx) error
	SendTelegramMessage(chatID int64, message string) error
}

type TelegramSendRequest struct {
	ChatID  int64  `json:"chat_id"`
	Message string `json:"message"`
}
