package blizz

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"strings"
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
	SearchItem(itemName string, region string) (*ItemResult, error)
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

func (c *client) SearchItem(itemName string, region string) (*ItemResult, error) {
	requestURL := c.urls[region] + "/data/wow/search/item"

	var locale string
	if isRussian(itemName) {
		// Проверяем либо кирилицу
		locale = "name.ru_RU"
	} else {
		// либо устанавливает английский язык для поиска предмета
		locale = "name.en_US"
	}

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"namespace":    fmt.Sprintf("static-%s", region),
			"access_token": c.token.AccessToken,
			"locale":       "ru_RU",
			"_page":        "1",
			"_pageSize":    "25",
			locale:         itemName,
		},
	}

	response, err := c.makeGetRequest(requestURL, ro)
	if err != nil {
		return nil, fmt.Errorf("err making SEARCH ITEM request: %v", err)
	}

	itemData := new(ItemResult)
	if err := response.JSON(itemData); err != nil {
		return nil, fmt.Errorf(
			"error unmarshaling realm list response: %v", err,
		)
	}

	return itemData, nil
}

func (c *client) GetBlizzRealms() error {
	for _, region := range c.cfg.RegionList {
		if err := c.getBlizzRealms(region); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) getBlizzRealms(region string) error {
	requestURL := c.urls[region] + "/data/wow/realm/index"

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"namespace":    fmt.Sprintf("dynamic-%s", region),
			"access_token": c.token.AccessToken,
			"locale":       "ru_RU",
		},
	}

	response, err := c.makeGetRequest(requestURL, ro)
	if err != nil {
		return fmt.Errorf("err making GET REALM request: %v", err)
	}

	realmData := new(BlizzRealmsSearchResult)
	if err := response.JSON(realmData); err != nil {
		return fmt.Errorf(
			"error unmarshaling realm list response: %v, region %s",
			err, region,
		)
	}

	c.setRealms(realmData)

	return nil
}

func (c *client) GetAuctionData(realmID int, region string) ([]*AuctionsDetail, error) {
	// Аукцион по реалму обновляется раз в час. В заголовке приходит дата обновления
	// last-modified: Thu, 31 Dec 2020 15:08:43 GMT
	// нужно сохранить данные локально для реалма и отдавать их из кеша в течение часа

	if data := c.getAuctionData(realmID, region); data != nil {
		return data, nil
	}

	requestURL := c.urls[region] + fmt.Sprintf("/data/wow/connected-realm/%d/auctions", realmID)
	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"namespace":    fmt.Sprintf("dynamic-%s", region),
			"access_token": c.token.AccessToken,
		},
	}
	response, err := c.makeGetRequest(requestURL, ro)
	if err != nil {
		return nil, fmt.Errorf("err making AUCTION DATA request: %v", err)
	}

	auctionData := new(AuctionData)
	if err := response.JSON(auctionData); err != nil {
		return nil, fmt.Errorf(
			"error unmarshaling item media response: %v", err,
		)
	}

	updatedAt := response.Header.Get("last-modified") // GMT! (-3)
	updatedAtParsed, err := time.Parse(layoutUS, updatedAt)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing last-modified header in auction response: %v", err,
		)
	}
	c.cache.SetAuctionData(realmID, region, auctionData, &updatedAtParsed)

	return auctionData.Auctions, nil
}

func (c *client) GetItemMedia(itemID string) (*ItemMedia, error) {
	requestURL := c.urls["eu"] + fmt.Sprintf("/data/wow/media/item/%s", itemID)
	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"namespace":    "static-eu",
			"access_token": c.token.AccessToken,
		},
	}
	response, err := c.makeGetRequest(requestURL, ro)
	if err != nil {
		return nil, fmt.Errorf("err making ITEM MEDIA request: %v", err)
	}

	itemMedia := new(ItemMedia)
	if err := response.JSON(itemMedia); err != nil {
		return nil, fmt.Errorf(
			"error unmarshaling item media response: %v", err,
		)
	}

	return itemMedia, nil
}

func (c *client) BlizzAuthRoutine() {
	delay := 6 * time.Hour
	t := time.NewTicker(delay)
	defer t.Stop()

	for range t.C {
		c.log.Info("Making blizzard auth call")
		if err := c.MakeBlizzAuth(); err != nil {
			c.log.Errorf("error making blizzard auth request: %v", err)
		}
	}
}

func (c *client) MakeBlizzAuth() error {
	body := strings.NewReader("grant_type=client_credentials")
	ro := &grequests.RequestOptions{
		RequestBody: body,
		Auth:        []string{c.cfg.ClientID, c.cfg.ClientSecret},
		Headers: map[string]string{
			"Content-Type": "application/x-www-form-urlencoded",
		},
	}

	response, err := c.session.Post(c.cfg.AUTHUrl, ro)
	if err != nil {
		return fmt.Errorf("error making blizzard auth request: %v", err)
	}

	tokenData := new(BlizzardToken)
	if err := response.JSON(tokenData); err != nil {
		return fmt.Errorf("error unmarshaling blizzard auth response: %v", err)
	}

	c.token = tokenData

	return nil
}

func (c *client) setRealms(realms *BlizzRealmsSearchResult) {
	for _, realm := range realms.Realms {
		c.cache.SetRealmID(realm.Name, realm.ID)
	}
}

func (c *client) GetRealmID(RealmName string) int {
	return c.cache.GetRealmID(RealmName)
}

func (c *client) getAuctionData(realmID int, region string) []*AuctionsDetail {
	data := c.cache.GetAuctionData(realmID, region)

	if t, ok := data.(*AuctionData); ok {
		return t.Auctions
	}

	return nil
}
