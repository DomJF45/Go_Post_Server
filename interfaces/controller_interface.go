package interfaces

import (
	"github.com/gofiber/fiber/v2"
)

type ControllerInterface interface {
	InitRouter(a *fiber.App)
}
