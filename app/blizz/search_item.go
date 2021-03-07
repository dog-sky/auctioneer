package blizz

import (
	"fmt"

	"github.com/levigross/grequests"
)

type ItemResultResultsDataName struct {
	RuRU string `json:"ru_RU"`
	EnGB string `json:"en_GB"`
	EnUS string `json:"en_US"`
}

type ItemResultResultsDataQuality struct {
	Type string `json:"type"`
}

type ItemResultResultsData struct {
	Name    ItemResultResultsDataName    `json:"name"`
	ID      int                          `json:"id"`
	Quality ItemResultResultsDataQuality `json:"quality"`
}

type ItemResultResults struct {
	Data ItemResultResultsData `json:"data"`
}

type Item struct {
	Results []ItemResultResults `json:"results"`
}

func (c *client) SearchItem(itemName string, region string) (*Item, error) {
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

	itemData := new(Item)
	if err := response.JSON(itemData); err != nil {
		return nil, fmt.Errorf(
			"error unmarshaling realm list response: %v", err,
		)
	}

	return itemData, nil
}
