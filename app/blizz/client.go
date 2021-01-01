package blizz

import (
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	"encoding/json"
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"net/http"
	"net/url"
	"strings"
)

type Client interface {
	GetBlizzRealms() error
	MakeBlizzAuth() error
	setRealms(*BlizzRealmsSearchResult)
	GetRealmID(string) int
	SearchItem(itemName string, region string) (*ItemResult, error)
}

type client struct {
	Cache      cache.Cache
	token      *BlizzardToken
	cfg        *conf.BlizzApiCfg
	httpClient *http.Client
}

func NewClient(blizzCfg *conf.BlizzApiCfg, cache cache.Cache) Client {
	return &client{
		cfg:        blizzCfg,
		httpClient: new(http.Client),
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
	// q.Set("region", region)
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
	requestURL, err := url.Parse(
		fmt.Sprintf(c.cfg.APIUrl+"/data/wow/realm/index", "eu"),
	)
	if err != nil {
		return fmt.Errorf(
			"Error creating realm request url: %v",
			err,
		)
	}
	q := requestURL.Query()
	q.Set("namespace", "dynamic-eu")
	q.Set("locale", "ru_RU")
	q.Set("access_token", c.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf(
			"Error creating realm request: %v",
			err,
		)
	}

	response, err := c.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf(
			"Error making get realm request: %v",
			err,
		)
	}
	if response.StatusCode != fiber.StatusOK {
		return fmt.Errorf(
			"Error making get realm request, status: %v",
			response.Status,
		)
	}
	defer response.Body.Close()

	realmData := new(BlizzRealmsSearchResult)
	if err := json.NewDecoder(response.Body).Decode(realmData); err != nil {
		return fmt.Errorf(
			"Error unmarshaling realm list response: %v",
			err,
		)
	}

	c.setRealms(realmData)

	return nil
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
