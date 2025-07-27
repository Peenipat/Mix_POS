package coreServices

import (
	coreModels "myapp/modules/core/models"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
)

type TelegramService struct{}

func NewTelegramService() *TelegramService {
	return &TelegramService{}
}

func (s *TelegramService) ProcessWebhook(c *fiber.Ctx) error {
	var req coreModels.TelegramWebhookRequest
	if err := c.BodyParser(&req); err != nil {
		log.Println("Invalid telegram payload:", err)
		return c.SendStatus(fiber.StatusBadRequest)
	}

	chatID := req.Message.Chat.ID
	username := req.Message.From.Username
	text := req.Message.Text

	log.Printf("📨 %s (chat_id: %d): %s", username, chatID, text)

	go s.SendTelegramMessage(chatID, "เชื่อมต่อกับระบบเรียบร้อยแล้ว")

	return c.SendStatus(fiber.StatusOK)
}

func (s *TelegramService) SendTelegramMessage(chatID int64, message string) error {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", os.Getenv("TELEGRAM_TOKEN"))

	payload := map[string]interface{}{
		"chat_id": chatID,
		"text":    message,
	}
	body, _ := json.Marshal(payload)

	resp, err := http.Post(url, "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("telegram error: %s", resp.Status)
	}
	return nil
}
