package system

import (
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	Ping(*fiber.Ctx) error
}

type SystemHandler struct{}

func NewSystemHandler() Handler {
	return new(SystemHandler)
}
