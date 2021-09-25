package blizz

import (
	"fmt"

	"github.com/levigross/grequests"
	"github.com/pkg/errors"
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
		return nil, errors.Wrapf(err, "GetItemMedia makeGetRequest")
	}

	itemMedia := new(ItemMedia)
	if err := response.JSON(itemMedia); err != nil {
		return nil, errors.Wrapf(err, "GetItemMedia JSON")
	}

	return itemMedia, nil
}
