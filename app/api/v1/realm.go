package v1

import (
	"encoding/json"
	"fmt"
	fiber "github.com/gofiber/fiber/v2"
	"net/http"
	"net/url"
)

func (h *v1Handler) GetBlizzRealms() error {
	requestURL, err := url.Parse(
		h.BlizzApiCfg.APIUrl.String() + "/data/wow/realm/index",
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
	q.Set("access_token", h.token.AccessToken)
	requestURL.RawQuery = q.Encode()

	request, err := http.NewRequest(http.MethodGet, requestURL.String(), nil)
	if err != nil {
		return fmt.Errorf(
			"Error creating realm request: %v",
			err,
		)
	}

	response, err := h.httpClient.Do(request)
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
	if err := json.NewDecoder(response.Body).Decode(&realmData); err != nil {
		return fmt.Errorf(
			"Error unmarshaling realm list response: %v",
			err,
		)
	}

	h.setRealms(realmData)

	return nil
}

func (h *v1Handler) setRealms(realms *BlizzRealmsSearchResult) {
	for _, realm := range realms.Realms {
		h.cache.SetRealmID(realm.Name, realm.ID)
	}
}
