package v1

import (
	"auctioneer/app/cache"
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
	cache       *cache.Cache
	// log      *logging.Logger
}

func NewBasehandlerv1(blizzCfg *conf.BlizzApiCfg, cache *cache.Cache) Handler {
	return &v1Handler{
		BlizzApiCfg: blizzCfg,
		httpClient:  new(http.Client),
		cache:       cache,
	}

}
