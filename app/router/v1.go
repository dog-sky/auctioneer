package router

import (
	"auctioneer/app/api/v1"
	"github.com/gofiber/fiber/v2"
)

func SetupV1Routes(r fiber.Router, h v1.Handler) {
	r.Get("/auc_search", h.SearchItemData)
}
