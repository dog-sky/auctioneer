package blizz

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"auctioneer/app/cache"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"

	"github.com/gofiber/fiber/v2"
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
	cache      cache.Cache
	token      *BlizzardToken
	cfg        *conf.BlizzApiCfg
	httpClient *http.Client
	urls       map[string]string
	log        *logging.Logger
}

func NewClient(logger *logging.Logger, blizzCfg *conf.BlizzApiCfg) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	urlsMap := make(map[string]string)
	urlsMap["eu"] = blizzCfg.EuAPIUrl
	urlsMap["us"] = blizzCfg.UsAPIUrl
	return &client{
		cfg: blizzCfg,
		httpClient: &http.Client{
			Transport: tr,
			Timeout:   time.Second * 10,
		},
		cache: cache.NewCache(),
		urls:  urlsMap,
		log:   logger,
	}
}

func (c *client) makeGetRequest(requestURL string) (*http.Response, error) {
	request, _ := http.NewRequest(http.MethodGet, requestURL, nil)
	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"error making get request: %v", err,
		)
	}
	if response.StatusCode != fiber.StatusOK {
		return nil, fmt.Errorf(
			"error making get request, status: %v", response.Status,
		)
	}

	return response, nil
}

func (c *client) SearchItem(itemName string, region string) (*ItemResult, error) {
	requestURL, _ := url.Parse(c.urls[region] + "/data/wow/search/item")
	q := requestURL.Query()
	q.Set("namespace", fmt.Sprintf("static-%s", region))
	q.Set("access_token", c.token.AccessToken)
	q.Set("_page", "1")
	q.Set("_pageSize", "25")
	if isRussian(itemName) {
		// Проверяем либо кирилицу
		q.Set("name.ru_RU", itemName)
	} else {
		// либо устанавливает английский язык для поиска предмета
		q.Set("name.en_US", itemName)
	}
	requestURL.RawQuery = q.Encode()

	response, err := c.makeGetRequest(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("err making SEARCH ITEM request: %v", err)
	}
	defer response.Body.Close()

	itemData := new(ItemResult)
	if err := json.NewDecoder(response.Body).Decode(itemData); err != nil {
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
	requestURL, _ := url.Parse(c.urls[region] + "/data/wow/realm/index")
	q := requestURL.Query()
	q.Set("namespace", fmt.Sprintf("dynamic-%s", region))
	q.Set("locale", "ru_RU")
	q.Set("access_token", c.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	response, err := c.makeGetRequest(requestURL.String())
	if err != nil {
		return fmt.Errorf("err making GET REALM request: %v", err)
	}
	defer response.Body.Close()

	realmData := new(BlizzRealmsSearchResult)
	if err := json.NewDecoder(response.Body).Decode(realmData); err != nil {
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

	data := c.getAuctionData(realmID, region)
	if data != nil {
		return data, nil
	}

	requestURL, _ := url.Parse(c.urls[region] + fmt.Sprintf("/data/wow/connected-realm/%d/auctions", realmID))
	q := requestURL.Query()
	q.Set("namespace", fmt.Sprintf("dynamic-%s", region))
	q.Set("access_token", c.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	response, err := c.makeGetRequest(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("err making AUCTION DATA request: %v", err)
	}
	defer response.Body.Close()

	auctionData := new(AuctionData)
	if err := json.NewDecoder(response.Body).Decode(auctionData); err != nil {
		return nil, fmt.Errorf(
			"error unmarshaling action data response: %v", err,
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
	requestURL, _ := url.Parse(c.urls["eu"] + fmt.Sprintf("/data/wow/media/item/%s", itemID))
	q := requestURL.Query()
	q.Set("namespace", "static-eu")
	q.Set("access_token", c.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	response, err := c.makeGetRequest(requestURL.String())
	if err != nil {
		return nil, fmt.Errorf("err making ITEM MEDIA request: %v", err)
	}
	defer response.Body.Close()

	itemMedia := new(ItemMedia)
	if err := json.NewDecoder(response.Body).Decode(itemMedia); err != nil {
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

	request, _ := http.NewRequest(http.MethodPost, c.cfg.AUTHUrl, body)
	request.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("error making blizzard auth request: %v", err)
	}
	defer response.Body.Close()

	tokenData := new(BlizzardToken)
	if err := json.NewDecoder(response.Body).Decode(tokenData); err != nil {
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

	switch t := data.(type) {
	case *AuctionData:
		return t.Auctions
	default:
		return nil
	}
}
