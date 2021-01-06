package router

import (
	"auctioneer/app/api"
	"auctioneer/app/api/system"
	"auctioneer/app/api/v1"

	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(a *fiber.App, h api.Handler) {
	v1 := a.Group("/api/v1")
	system := a.Group("/")

	setupV1Routes(v1, h.V1Handler())
	setupSystemRoutes(system, h.SystemHandler())
}

func setupSystemRoutes(r fiber.Router, h system.Handler) {
	r.Get("/ping", h.Ping)
}

func setupV1Routes(r fiber.Router, h v1.Handler) {
	r.Get("/auc_search", h.SearchItemData)
	r.Get("/item_media/:item_id", h.SearchItemMedia)
}
