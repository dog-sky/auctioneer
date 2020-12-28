package v1

import (
	logging "auctioneer/app/logger"
)

type Handler interface {
}

type v1Handler struct {
	log *logging.Logger
}

func NewBasehandlerv1(log *logging.Logger) Handler {
	return &v1Handler{
		log: log,
	}
}
