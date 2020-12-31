package v1

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func (h *v1Handler) MakeBlizzAuth() error {
	body := strings.NewReader("grant_type=client_credentials")

	request, err := http.NewRequest(
		http.MethodPost,
		h.BlizzApiCfg.AUTHUrl.String(),
		body,
	)
	if err != nil {
		return fmt.Errorf(
			"Error creating request: %v",
			err,
		)
	}

	request.SetBasicAuth(h.BlizzApiCfg.ClientID, h.BlizzApiCfg.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response, err := h.httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Error making blizzard auth request: %v", err)
	}
	defer response.Body.Close()

	tokenData := new(BlizzardToken)

	if err := json.NewDecoder(response.Body).Decode(&tokenData); err != nil {
		return fmt.Errorf(
			"Error unmarshaling blizzard auth response: %v",
			err,
		)
	}

	h.token = tokenData

	return nil
}
