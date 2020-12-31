package v1

import (
	"auctioneer/app/conf"
	"net/http"
)

type Handler interface {
	MakeBlizzAuth() error
	GetBlizzRealms() error
}

type v1Handler struct {
	token       *BlizzardToken
	BlizzApiCfg *conf.BlizzApiCfg
	httpClient  *http.Client
	// log      *logging.Logger
}

func NewBasehandlerv1(blizzCfg *conf.BlizzApiCfg) Handler {
	return &v1Handler{
		BlizzApiCfg: blizzCfg,
		httpClient:  new(http.Client),
	}

}
