package main

import (
	"post_server/controllers"
	"post_server/interfaces"
)

func main() {
	projectController := controllers.NewProjectController()

	controllers := []interfaces.ControllerInterface{
		projectController,
	}

	app := NewApp(controllers, ":8080")

	app.Run()
}
