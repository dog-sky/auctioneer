package api

import (
	"auctioneer/app/api/v1"
	"auctioneer/app/conf"
)

type Handler interface {
}

type BaseHandler struct {
	V1 v1.Handler
}

func NewBasehandler(cfg *conf.Config) Handler {
	return &BaseHandler{
		V1: v1.NewBasehandlerv1(&cfg.BlizzApiCfg),
	}
}
