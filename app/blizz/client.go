package blizz

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"auctioneer/app/cache"
	"auctioneer/app/conf"

	"github.com/levigross/grequests"
	"github.com/sirupsen/logrus"
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
	ctx     context.Context
	log     *logrus.Logger
}

func NewClient(ctx context.Context, logger *logrus.Logger, blizzCfg *conf.BlizzApiCfg) Client {
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
		ctx:     ctx,
		log:     logger,
	}
}

func (c *client) makeGetRequest(requestURL string, ro *grequests.RequestOptions) (*grequests.Response, error) {
	response, err := c.session.Get(requestURL, ro)
	if err != nil {
		return nil, errors.Wrapf(err, "makeGetRequest")
	}
	if !response.Ok {

		return nil, fmt.Errorf(
			"error making get request, status: %v", response.StatusCode,
		)
	}

	return response, nil
}
