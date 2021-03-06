package blizz

import (
	"fmt"

	"github.com/levigross/grequests"
)

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
