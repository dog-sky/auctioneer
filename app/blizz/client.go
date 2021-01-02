package blizz

import (
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	"crypto/tls"
	"encoding/json"
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const layoutUS = "Mon, 2 Jan 2006 15:04:05 MST"

type Client interface {
	GetBlizzRealms() error
	getBlizzRealms(string) error
	MakeBlizzAuth() error
	setRealms(*BlizzRealmsSearchResult)
	GetRealmID(string) int
	SearchItem(itemName string, region string) (*ItemResult, error)
	GetAuctionData(realmID int, region string) ([]*AuctionsDetail, error)
	getAuctionData(realmID int, region string) []*AuctionsDetail
	setAuctionData(realmID int, region string, auctionData *AuctionData, updatedAt *time.Time)
}

type client struct {
	Cache      cache.Cache
	token      *BlizzardToken
	cfg        *conf.BlizzApiCfg
	httpClient *http.Client
}

func NewClient(blizzCfg *conf.BlizzApiCfg, cache cache.Cache) Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &client{
		cfg:        blizzCfg,
		httpClient: &http.Client{Transport: tr},
		Cache:      cache,
	}
}

func (c *client) SearchItem(itemName string, region string) (*ItemResult, error) {
	requestURL, err := url.Parse(
		fmt.Sprintf(c.cfg.APIUrl+"/data/wow/search/item", region),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"Error creating item search request url: %v",
			err,
		)
	}
	q := requestURL.Query()
	q.Set("namespace", fmt.Sprintf("static-%s", region))
	q.Set("access_token", c.token.AccessToken)
	if itemName != "" {
		if isRussian(itemName) {
			// Проверяем либо кирилицу
			q.Set("name.ru_RU", itemName)
		} else {
			// либо устанавливает английский язык для поиска предмета
			q.Set("name.en_US", itemName)
		}
	}
	requestURL.RawQuery = q.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(
			"Error creating item search request: %v",
			err,
		)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"Error making search item request: %v",
			err,
		)
	}
	if response.StatusCode != fiber.StatusOK {
		return nil, fmt.Errorf(
			"Error making search item request, status: %v",
			response.Status,
		)
	}

	defer response.Body.Close()

	itemData := new(ItemResult)
	if err := json.NewDecoder(response.Body).Decode(itemData); err != nil {
		return nil, fmt.Errorf(
			"Error unmarshaling realm list response: %v",
			err,
		)
	}

	return itemData, nil
}

func (c *client) GetBlizzRealms() error {
	if err := c.getBlizzRealms("eu"); err != nil {
		return err
	}
	if err := c.getBlizzRealms("us"); err != nil {
		return err
	}

	return nil
}

func (c *client) getBlizzRealms(region string) error {
	requestURL, err := url.Parse(
		fmt.Sprintf(c.cfg.APIUrl+"/data/wow/realm/index", region),
	)
	if err != nil {
		return fmt.Errorf(
			"Error creating realm request url: %v",
			err,
		)
	}
	q := requestURL.Query()
	q.Set("namespace", fmt.Sprintf("dynamic-%s", region))
	q.Set("locale", "ru_RU")
	q.Set("access_token", c.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf(
			"Error creating realm request: %v, region %s",
			err, region,
		)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(
			"Error making get realm request: %v, region %s",
			err, region,
		)
	}
	if response.StatusCode != fiber.StatusOK {
		return fmt.Errorf(
			"Error making get realm request, status: %v, region %s",
			response.Status, region,
		)
	}
	defer response.Body.Close()

	realmData := new(BlizzRealmsSearchResult)
	if err := json.NewDecoder(response.Body).Decode(realmData); err != nil {
		return fmt.Errorf(
			"Error unmarshaling realm list response: %v, region %s",
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

	requestURL, err := url.Parse(
		fmt.Sprintf(c.cfg.APIUrl+"/data/wow/connected-realm/%d/auctions", region, realmID),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"Error creating action request url: %v", err,
		)
	}

	q := requestURL.Query()
	q.Set("namespace", fmt.Sprintf("dynamic-%s", region))
	q.Set("access_token", c.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return nil, fmt.Errorf(
			"Error creating action request: %v", err,
		)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf(
			"Error making get auction request: %v", err,
		)
	}
	if response.StatusCode != fiber.StatusOK {
		return nil, fmt.Errorf(
			"Error making get auction request, status: %v", response.Status,
		)
	}
	defer response.Body.Close()

	auctionData := new(AuctionData)
	if err := json.NewDecoder(response.Body).Decode(auctionData); err != nil {
		return nil, fmt.Errorf(
			"Error unmarshaling action data response: %v", err,
		)
	}

	updatedAt := response.Header.Get("last-modified") // GMT! (-3)
	updatedAtParsed, err := time.Parse(layoutUS, updatedAt)
	if err != nil {
		return nil, fmt.Errorf(
			"Error parsing last-modified header in auction response: %v", err,
		)
	}
	c.setAuctionData(realmID, region, auctionData, &updatedAtParsed)

	return auctionData.Auctions, nil
}

func (c *client) MakeBlizzAuth() error {
	body := strings.NewReader("grant_type=client_credentials")

	request, err := http.NewRequest(
		http.MethodPost,
		c.cfg.AUTHUrl.String(),
		body,
	)
	if err != nil {
		return fmt.Errorf(
			"Error creating request: %v",
			err,
		)
	}

	request.SetBasicAuth(c.cfg.ClientID, c.cfg.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Error making blizzard auth request: %v", err)
	}
	defer response.Body.Close()

	tokenData := new(BlizzardToken)
	if err := json.NewDecoder(response.Body).Decode(tokenData); err != nil {
		return fmt.Errorf(
			"Error unmarshaling blizzard auth response: %v",
			err,
		)
	}

	c.token = tokenData

	return nil
}

func (c *client) setRealms(realms *BlizzRealmsSearchResult) {
	for _, realm := range realms.Realms {
		c.Cache.SetRealmID(realm.Name, realm.ID)
	}
}

func (c *client) GetRealmID(RealmName string) int {
	return c.Cache.GetRealmID(RealmName)
}

func (c *client) getAuctionData(realmID int, region string) []*AuctionsDetail {
	data := c.Cache.GetAcutionData(realmID, region)

	switch t := data.(type) {
	case *AuctionData:
		return t.Auctions
	default:
		return nil
	}
}

func (c *client) setAuctionData(realmID int, region string, auctionData *AuctionData, updatedAt *time.Time) {
	c.Cache.SetAuctionData(realmID, region, auctionData, updatedAt)
}
