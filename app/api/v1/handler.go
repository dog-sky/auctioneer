package v1

import (
	"auctioneer/app/blizz"
	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SearchItemData(*fiber.Ctx) error
	MakeBlizzAuth() error
	GetBlizzRealms() error
}

type V1Handler struct {
	blizzClient blizz.Client
	// log      *logging.Logger
}

func NewBasehandlerv1(blizzClient blizz.Client) Handler {
	return &V1Handler{
		blizzClient: blizzClient,
	}
}

func (h *V1Handler) MakeBlizzAuth() error {
	return h.blizzClient.MakeBlizzAuth()
}

func (h *V1Handler) GetBlizzRealms() error {
	return h.blizzClient.GetBlizzRealms()
}
