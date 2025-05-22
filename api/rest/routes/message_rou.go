package routes

import (
	"go-fiber/adapter"
	"go-fiber/api/rest/controllers"
	"go-fiber/data/services"
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func NewMessageRoutes(router fiber.Router, db *gorm.DB, cli *http.Client) {
	client := adapter.NewCustomHTTPClient(cli)
	msvc := services.NewMessageService(client)
	gmnsvc := services.NewGeminiServices()
	mctrl := controllers.NewMessageCtrl(msvc, gmnsvc)

	router.Get("/messages/webhook", mctrl.VerifyWebhook)
	router.Post("/messages/webhook", mctrl.HandleMessengerWebhook)

}
