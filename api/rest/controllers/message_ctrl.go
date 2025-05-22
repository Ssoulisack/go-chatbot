package controllers

import (
	"fmt"
	"go-fiber/api/rest/middleware"
	"go-fiber/bootstrap"
	"go-fiber/core/logs"
	"go-fiber/data/services"
	"go-fiber/domain/models"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type MessageCtrl interface {
	HandleMessengerWebhook(c *fiber.Ctx) error
	VerifyWebhook(c *fiber.Ctx) error
}

type messageCtrl struct {
	geminiSvc  services.GeminiServices
	messageSvc services.MessageService
}

func (m *messageCtrl) HandleMessengerWebhook(c *fiber.Ctx) error {
	var req models.WebhookRequest
	if err := c.BodyParser(&req); err != nil {
		logs.Error(fmt.Sprintf("Error parsing request body: %v", err))
		return middleware.ErrorBadRequest("Invalid request body")
	}

	if len(req.Entry) == 0 || len(req.Entry[0].Messaging) == 0 {
		return middleware.ErrorBadRequest("Invalid payload")
	}
	message := req.Entry[0].Messaging[0]

	reply := m.geminiSvc.GenerateReply(message.Message.Text)

	logs.Info(fmt.Sprintf("Received message from %s: %s", message.Sender.ID, message.Message.Text))

	err := m.messageSvc.SendFacebookReply(message.Sender.ID, reply)
	if err != nil {
		logs.Error(fmt.Sprintf("Error sending reply: %v", err))
		return middleware.NewErrorErrMsgInternalServerError(c)
	}

	return c.SendStatus(http.StatusOK)
}

func (m *messageCtrl) VerifyWebhook(c *fiber.Ctx) error {
	mode := c.Query("hub.mode")
	token := c.Query("hub.verify_token")
	challenge := c.Query("hub.challenge")
	verifyToken := bootstrap.GlobalEnv.Keys.FbVerifyToken
	if mode == "subscribe" && token == verifyToken {
		logs.Info(fmt.Sprintf("Webhook verified with token: %s", token))
		return c.SendString(challenge)
	}
	logs.Info(fmt.Sprintf("Webhook verification failed with token: %s", token))
	return middleware.ErrorExpectationFailed("Invalid verification token")
}


func NewMessageCtrl(messageSvc services.MessageService, geminiSvc services.GeminiServices) MessageCtrl {
	return &messageCtrl{
		messageSvc: messageSvc,
		geminiSvc:  geminiSvc,
	}
}
