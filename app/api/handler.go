package api

import (
	"auctioneer/app/api/v1"
	"auctioneer/app/blizz"
)

type Handler interface {
	V1MakeBlizzAuth() error
	V1GetBlizzRealms() error
	V1Handler() v1.Handler
}

type BaseHandler struct {
	V1 v1.Handler
}

func NewBasehandler(blizzClient blizz.Client) Handler {
	return &BaseHandler{
		V1: v1.NewBasehandlerv1(blizzClient),
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
