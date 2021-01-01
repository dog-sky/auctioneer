package v1

import (
	"auctioneer/app/blizz"
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SearchItemData(*fiber.Ctx) error
	MakeBlizzAuth() error
	GetBlizzRealms() error
}

type V1Handler struct {
	BlizzClient blizz.Client
	// log      *logging.Logger
}

func NewBasehandlerv1(blizzCfg *conf.BlizzApiCfg, cache cache.Cache) Handler {
	return &V1Handler{
		BlizzClient: blizz.NewClient(blizzCfg, cache),
	}
}

func (h *V1Handler) MakeBlizzAuth() error {
	return h.BlizzClient.MakeBlizzAuth()
}

func (h *V1Handler) GetBlizzRealms() error {
	return h.BlizzClient.GetBlizzRealms()
}
