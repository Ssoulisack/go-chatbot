package bootstrap

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type Application struct {
	Env    *Env
	Fiber  *fiber.App
	DB     *gorm.DB
	Client *http.Client
}

var GlobalEnv Env

func App() *Application {
	app := &Application{}
	app.Env = NewEnv()
	GlobalEnv = *NewEnv()
	app.Fiber = NewFiber()
	app.Client = NewHttpClient()
	app.DB = NewDatabaseConnectionPostgres(app.Env)
	return app
}
