package main

import (
	"fmt"
	"go-fiber/api/rest/routes"
	"go-fiber/bootstrap"
	"log"
)

func main() {
	app := bootstrap.App()
	globalEnv := app.Env
	fiber := app.Fiber
	client := app.Client
	db := app.DB
	routes.Setup(fiber, db, client)
	log.Fatal(fiber.Listen(fmt.Sprintf(":%d", globalEnv.App.Port)))
}
