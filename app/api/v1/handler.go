package v1

import (
	"auctioneer/app/conf"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Handler interface {
	MakeBlizzAuth() error
}

type v1Handler struct {
	token       *BlizzardToken
	BlizzApiCfg *conf.BlizzApiCfg
	// log      *logging.Logger
}

func NewBasehandlerv1(blizzCfg *conf.BlizzApiCfg) Handler {
	return &v1Handler{
		BlizzApiCfg: blizzCfg,
	}

}

func (h *v1Handler) MakeBlizzAuth() error {
	body := strings.NewReader("grant_type=client_credentials")

	request, err := http.NewRequest(http.MethodPost, h.BlizzApiCfg.AUTHUrl.String(), body)
	if err != nil {
		return fmt.Errorf("Error creating request: %v", err)
	}
	request.SetBasicAuth(h.BlizzApiCfg.ClientID, h.BlizzApiCfg.ClientSecret)
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	httpClient := new(http.Client)
	response, err := httpClient.Do(request)
	if err != nil {
		return fmt.Errorf("Error making blizzard auth request: %v", err)
	}
	defer response.Body.Close()

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("Error reading blizzard auth response: %v", err)
	}

	tokenData := new(BlizzardToken)
	err = json.Unmarshal(responseData, &tokenData)
	if err != nil {
		return fmt.Errorf("Error unmarshaling blizzard auth response: %v", err)
	}

	h.token = tokenData

	return nil
}
