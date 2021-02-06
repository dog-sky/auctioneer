package middleware

import (
	v1 "auctioneer/app/api/v1"

	"github.com/gofiber/fiber/v2"
)

func ErrorHandler(c *fiber.Ctx, incomingError error) error {
	code := fiber.StatusInternalServerError
	resp := v1.ResponseV1{
		Success: false,
	}

	if e, ok := incomingError.(*fiber.Error); ok {
		code = e.Code
		resp.Message = e.Message
	} else {
		resp.Message = incomingError.Error()
	}

	return c.Status(code).JSON(resp)
}
