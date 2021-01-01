package v1

import (
	"auctioneer/app/blizz"
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SearchItemData(*fiber.Ctx) error
}

type V1Handler struct {
	BlizzClient *blizz.BlizzClient
	// log      *logging.Logger
}

func NewBasehandlerv1(blizzCfg *conf.BlizzApiCfg, cache *cache.Cache) Handler {
	return &V1Handler{
		BlizzClient: blizz.NewBlizzClient(blizzCfg, cache),
	}
}
