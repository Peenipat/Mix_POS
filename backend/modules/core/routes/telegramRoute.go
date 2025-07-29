package coreRoutes

import (
    "github.com/gofiber/fiber/v2"

    coreControllers "myapp/modules/core/controllers"
)

func RegisterTelegramRoutes(router fiber.Router, ctrl *coreControllers.TelegramController) {
    
    telegramGroup := router.Group("/telegram")
	telegramGroup.Post("/webhook", ctrl.HandleWebhook)
	telegramGroup.Post("/send", ctrl.HandleSendMessage)

}