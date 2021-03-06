package blizz

import (
	"fmt"

	"github.com/levigross/grequests"
)

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
