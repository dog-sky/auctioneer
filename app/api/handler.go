package api

import (
	"auctioneer/app/api/v1"
	logging "auctioneer/app/logger"
)

type Handler interface {
}

type BaseHandler struct {
	V1 v1.Handler
}

func NewBasehandler(log *logging.Logger) Handler {
	return &BaseHandler{
		V1: v1.NewBasehandlerv1(log),
	}
}
