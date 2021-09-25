package blizz

import (
	"fmt"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
)

type BlizzRealmsSearchResultResultsDataRealmsName struct {
	RuRU string `json:"ru_RU"`
	EnGB string `json:"en_GB"`
}

type BlizzRealmsSearchResultResultsDataRealms struct {
	Name BlizzRealmsSearchResultResultsDataRealmsName `json:"name"`
}

type BlizzRealmsSearchResultResultsData struct {
	Realms []BlizzRealmsSearchResultResultsDataRealms `json:"realms"`
	ID     int                                        `json:"id"`
}

type BlizzRealmsSearchResultResults struct {
	Data BlizzRealmsSearchResultResultsData `json:"data"`
}

// Итоговая структура. Из этого нужно брать имя сервера и значением будет айди коннектед риалма
type BlizzRealmsSearchResult struct {
	Results []BlizzRealmsSearchResultResults `json:"results"`
}

func (c *client) GetBlizzRealms() error {
	for _, region := range c.cfg.RegionList {
		if err := c.getBlizzRealms(region); err != nil {
			return errors.Wrapf(err, "GetBlizzRealms")
		}
	}

	return nil
}

func (c *client) GetRealmID(RealmName string) int {
	return c.cache.GetRealmID(RealmName)
}

func (c *client) setRealms(realms *BlizzRealmsSearchResult) {
	for _, connectedRealm := range realms.Results {
		for _, realm := range connectedRealm.Data.Realms {
			c.cache.SetRealmID(realm.Name.RuRU, connectedRealm.Data.ID)
			c.cache.SetRealmID(realm.Name.EnGB, connectedRealm.Data.ID)
		}
	}
}

func (c *client) getBlizzRealms(region string) error {
	requestURL := c.urls[region] + "/data/wow/search/connected-realm"

	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"namespace":    fmt.Sprintf("dynamic-%s", region),
			"access_token": c.token.AccessToken,
			"locale":       "ru_RU",
		},
	}

	response, err := c.makeGetRequest(requestURL, ro)
	if err != nil {
		return errors.Wrapf(err, "getBlizzRealms makeGetRequest")
	}

	realmData := new(BlizzRealmsSearchResult)
	if err := response.JSON(realmData); err != nil {
		return errors.Wrapf(err, "getBlizzRealms JSON")
	}

	c.setRealms(realmData)

	return nil
}
