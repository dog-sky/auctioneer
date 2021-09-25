package system

import (
	"github.com/gofiber/fiber/v2"
)

func (h *SystemHandler) Ping(c *fiber.Ctx) error {
	return c.SendStatus(fiber.StatusOK)
}
