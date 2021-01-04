package router

import (
	"auctioneer/app/api/system"
	"auctioneer/app/api/v1"
	"github.com/gofiber/fiber/v2"
)

func SetupSystemRoutes(r fiber.Router, h system.Handler) {
	r.Get("/ping", h.Ping)
}

func SetupV1Routes(r fiber.Router, h v1.Handler) {
	r.Get("/auc_search", h.SearchItemData)
}
