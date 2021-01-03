package v1_test

import (
	"testing"
	// "time"
	"fmt"
	"io/ioutil"
	"strings"

	// "auctioneer/app/api"
	"auctioneer/app/api/v1"
	server "auctioneer/app/auctioneer"
	"auctioneer/app/blizz"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	// "github.com/gofiber/fiber/v2"
	"net/http/httptest"
)

func newV1handler() v1.Handler {
	return &v1.V1Handler{
		BlizzClient: &mockBlizzClient{},
	}
}

type mockHandler struct {
	v1 v1.Handler
}

type mockBlizzClient struct{}

func (c *mockBlizzClient) GetBlizzRealms() error {
	return nil
}

func (c *mockBlizzClient) MakeBlizzAuth() error {
	return nil
}

func (c *mockBlizzClient) GetRealmID(s string) int {
	if s == "Killrog" {
		return 1
	}
	return 0
}

func (c *mockBlizzClient) SearchItem(itemName string, region string) (*blizz.ItemResult, error) {
	if strings.Contains(itemName, "гаррош") {
		return &blizz.ItemResult{
			Results: []blizz.ItemTesult{
				{
					Data: blizz.ItemData{
						Media: blizz.ItemMedia{
							ID: 1,
						},
						Name: blizz.DetailedName{
							RuRU: "Оправдание Гарроша",
							EnGB: "Garrosh's Pardon",
							EnUS: "Garrosh's Pardon",
						},
						ID: 1,
						Quality: blizz.ItemQuality{
							Type: "EPIC",
						},
					},
				},
			},
		}, nil
	}
	if strings.Contains(itemName, "опал") {
		return &blizz.ItemResult{
			Results: []blizz.ItemTesult{
				{
					Data: blizz.ItemData{
						Media: blizz.ItemMedia{
							ID: 2,
						},
						Name: blizz.DetailedName{
							RuRU: "Большой опал",
							EnGB: "Large Opal",
							EnUS: "Large Opal",
						},
						ID: 2,
						Quality: blizz.ItemQuality{
							Type: "UNCOMMON",
						},
					},
				},
			},
		}, nil
	}
	return &blizz.ItemResult{
		Results: []blizz.ItemTesult{},
	}, nil
}

func (c *mockBlizzClient) GetAuctionData(realmID int, region string) ([]*blizz.AuctionsDetail, error) {
	return []*blizz.AuctionsDetail{
		&blizz.AuctionsDetail{
			ID: 1,
			Item: blizz.AcuItem{
				ID:      1,
				Context: 1,
				Modifiers: []blizz.AucItemModifiers{
					blizz.AucItemModifiers{
						Type:  1,
						Value: 1,
					},
				},
				PetBreedID:   1,
				PetLevel:     1,
				PetQualityID: 1,
				PetSpeciesID: 1,
			},
			Buyout:   10001,
			Quantity: 2,
			TimeLeft: "233",
			ItemName: blizz.DetailedName{
				RuRU: "Оправдание Гарроша",
				EnGB: "Garrosh's Pardon",
				EnUS: "Garrosh's Pardon",
			},
			Quality: "EPIC",
		},
	}, nil
}

func (h *mockHandler) V1MakeBlizzAuth() error { return nil }

func (h *mockHandler) V1GetBlizzRealms() error { return nil }

func (h *mockHandler) V1Handler() v1.Handler {
	return h.v1
}

func TestV1Handler_SearchItemData(t *testing.T) {
	cfg, _ := conf.NewConfig()
	logger, _ := logging.NewLogger("DEBUG")
	app := server.NewApp(logger, cfg)
	app.BaseHandler = &mockHandler{
		v1: newV1handler(),
	}
	app.SetupRoutes()

	testCases := []struct {
		name      string
		reqURI    string
		expStatus int
		exp       v1.ResponseV1
	}{
		{
			name:      "Valid request",
			reqURI:    "?item_name=гаррош&region=eu&realm_name=Killrog",
			expStatus: 200,
			exp: v1.ResponseV1{
				Success: true,
				Result: []*blizz.AuctionsDetail{
					&blizz.AuctionsDetail{
						ID: 1,
						Item: blizz.AcuItem{
							ID:      1,
							Context: 1,
							Modifiers: []blizz.AucItemModifiers{
								blizz.AucItemModifiers{
									Type:  1,
									Value: 1,
								},
							},
							PetBreedID:   1,
							PetLevel:     1,
							PetQualityID: 1,
							PetSpeciesID: 1,
						},
						Buyout:   10001,
						Quantity: 2,
						TimeLeft: "233",
						ItemName: blizz.DetailedName{
							RuRU: "Оправдание Гарроша",
							EnGB: "Garrosh's Pardon",
							EnUS: "Garrosh's Pardon",
						},
						Quality: "EPIC",
					},
				},
			},
		},
		{
			name:      "Valid request, but item not in auction right now",
			reqURI:    "?item_name=опал&region=eu&realm_name=Killrog",
			expStatus: 200,
			exp: v1.ResponseV1{
				Success: true,
				Result:  []*blizz.AuctionsDetail{},
			},
		},
		{
			name:      "Item not found",
			reqURI:    "?item_name=алмаз&region=eu&realm_name=Killrog",
			expStatus: 404,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Item алмаз not found",
			},
		},
		{
			name:      "item_name not in request",
			reqURI:    "?region=eu&realm_name=Killrog",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Item name must not be empty",
			},
		},
		{
			name:      "region not in request",
			reqURI:    "?item_name=алмаз&realm_name=Killrog",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Region must not be empty",
			},
		},
		{
			name:      "realm_name not in request",
			reqURI:    "?item_name=алмаз&region=eu",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Realm name must not be empty",
			},
		},
		{
			name:      "Realm not found",
			reqURI:    "?item_name=алмаз&region=eu&realm_name=Гордунни",
			expStatus: 404,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Realm Гордунни not found",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			reqURI := fmt.Sprintf("/api/v1/auc_search%s", tc.reqURI)
			req := httptest.NewRequest("GET", reqURI, nil)

			resp, err := app.Fib.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expStatus, resp.StatusCode)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			var respV1 v1.ResponseV1
			err = json.Unmarshal(bodyBytes, &respV1)
			assert.NoError(t, err)

			assert.EqualValues(t, tc.exp, respV1)
		})
	}
}