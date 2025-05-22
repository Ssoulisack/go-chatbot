package routes

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Setup(app *fiber.App, db *gorm.DB, client *http.Client) {
	api := app.Group("/api/v1", func(ctx *fiber.Ctx) error {
		return ctx.Next()
	})

	NewUserRouter(api, db)

	NewMessageRoutes(api, db, client)

}
