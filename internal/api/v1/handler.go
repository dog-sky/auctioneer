package v1

import (
	"github.com/dog-sky/auctioneer/internal/client/blizz"

	"github.com/gofiber/fiber/v2"
)

type Handler interface {
	SearchItemData(*fiber.Ctx) error
	SearchItemMedia(c *fiber.Ctx) error
	MakeBlizzAuth() error
	GetBlizzRealms() error
}

type V1Handler struct {
	blizzClient blizz.Client
	// log      *logrus.Logger
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
