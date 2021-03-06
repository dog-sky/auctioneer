package blizz

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"auctioneer/app/cache"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"

	"github.com/levigross/grequests"
)

const layoutUS = "Mon, 2 Jan 2006 15:04:05 MST"

type Client interface {
	GetBlizzRealms() error
	MakeBlizzAuth() error
	GetRealmID(string) int
	BlizzAuthRoutine()
	GetItemMedia(itemID string) (*ItemMedia, error)
	SearchItem(itemName string, region string) (*Item, error)
	GetAuctionData(realmID int, region string) ([]*AuctionsDetail, error)
}

type client struct {
	cache   cache.Cache
	token   *BlizzardToken
	cfg     *conf.BlizzApiCfg
	session *grequests.Session
	urls    map[string]string
	log     *logging.Logger
}

func NewClient(logger *logging.Logger, blizzCfg *conf.BlizzApiCfg) Client {
	urlsMap := make(map[string]string)
	urlsMap["eu"] = blizzCfg.EuAPIUrl
	urlsMap["us"] = blizzCfg.UsAPIUrl

	session := grequests.NewSession(nil)
	session.HTTPClient.Transport = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	session.HTTPClient.Timeout = time.Second * 10

	return &client{
		cfg:     blizzCfg,
		cache:   cache.NewCache(),
		session: session,
		urls:    urlsMap,
		log:     logger,
	}
}

func (c *client) makeGetRequest(requestURL string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	response, err := c.session.Get(requestURL, ro)
	if err != nil {
		return nil, fmt.Errorf(
			"error making get request: %v", err,
		)
	}
	if !response.Ok {
		return nil, fmt.Errorf(
			"error making get request, status: %v", response.StatusCode,
		)
	}

	return response, nil
}
