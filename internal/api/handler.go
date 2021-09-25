package api

import (
	"github.com/dog-sky/auctioneer/internal/api/system"

	v1 "github.com/dog-sky/auctioneer/internal/api/v1"
	"github.com/dog-sky/auctioneer/internal/client/blizz"
)

type Handler interface {
	V1MakeBlizzAuth() error
	V1GetBlizzRealms() error
	V1Handler() v1.Handler
	SystemHandler() system.Handler
}

type BaseHandler struct {
	V1     v1.Handler
	system system.Handler
}

func NewBasehandler(blizzClient blizz.Client) Handler {
	return &BaseHandler{
		V1:     v1.NewBasehandlerv1(blizzClient),
		system: system.NewSystemHandler(),
	}
}

func (h *BaseHandler) V1MakeBlizzAuth() error {
	return h.V1.MakeBlizzAuth()
}

func (h *BaseHandler) V1GetBlizzRealms() error {
	return h.V1.GetBlizzRealms()
}

func (h *BaseHandler) V1Handler() v1.Handler {
	return h.V1
}

func (h *BaseHandler) SystemHandler() system.Handler {
	return h.system
}
