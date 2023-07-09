package main

import (
	"github.com/gofiber/fiber/v2"

	"post_server/interfaces"
)

type App struct {
	App  fiber.App
	Port string
}

func NewApp(controllers []interfaces.ControllerInterface, port string) *App {
	var app App
	app.App = *fiber.New()
	app.Port = port
	app.initControllers(controllers)
	return &app
}

func (a *App) initControllers(controllers []interfaces.ControllerInterface) {
	for _, controller := range controllers {
		controller.InitRouter(&a.App)
	}
}

func (a *App) Run() {
	a.App.Listen(a.Port)
}
