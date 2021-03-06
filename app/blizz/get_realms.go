package blizz

import (
	"fmt"

	"github.com/levigross/grequests"
)

type realm struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type BlizzRealmsSearchResult struct {
	Realms []realm `json:"realms"`
}

func (c *client) GetBlizzRealms() error {
	for _, region := range c.cfg.RegionList {
		if err := c.getBlizzRealms(region); err != nil {
			return err
		}
	}

	return nil
}

func (c *client) GetRealmID(RealmName string) int {
	return c.cache.GetRealmID(RealmName)
}

func (c *client) setRealms(realms *BlizzRealmsSearchResult) {
	for _, realm := range realms.Realms {
		c.cache.SetRealmID(realm.Name, realm.ID)
	}
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
