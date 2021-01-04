package v1_test

import (
	"fmt"
	"io/ioutil"
	"strings"
	"testing"

	"auctioneer/app/api/system"
	"auctioneer/app/api/v1"
	server "auctioneer/app/auctioneer"
	"auctioneer/app/blizz"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http/httptest"
)

func newV1handler() v1.Handler {
	return v1.NewBasehandlerv1(&mockBlizzClient{})
}

func newSystemHandler() system.Handler {
	return system.NewSystemHandler()
}

type mockHandler struct {
	v1     v1.Handler
	system system.Handler
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
	if s == "errRealm" {
		return 2
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
	if strings.Contains(itemName, "riseError") {
		return nil, fmt.Errorf(
			"error making get auction request, status: %d", 404,
		)
	}
	return &blizz.ItemResult{
		Results: []blizz.ItemTesult{},
	}, nil
}

func (c *mockBlizzClient) GetAuctionData(realmID int, region string) ([]*blizz.AuctionsDetail, error) {
	if realmID == 1 {
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
	return nil, fmt.Errorf(
		"error making GetAuctionData request, status: %d", 404,
	)
}

func (h *mockHandler) V1MakeBlizzAuth() error { return nil }

func (h *mockHandler) V1GetBlizzRealms() error { return nil }

func (h *mockHandler) V1Handler() v1.Handler {
	return h.v1
}

func (h *mockHandler) SystemHandler() system.Handler {
	return h.system
}

func TestV1Handler_SearchItemData(t *testing.T) {
	cfg, _ := conf.NewConfig()
	logger, _ := logging.NewLogger("DEBUG")
	app := server.NewApp(logger, cfg)
	app.BaseHandler = &mockHandler{
		v1:     newV1handler(),
		system: newSystemHandler(),
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
				Message: "Invalid query params. Err: item_name is empty",
			},
		},
		{
			name:      "region not in request",
			reqURI:    "?item_name=алмаз&realm_name=Killrog",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Invalid query params. Err: region is empty",
			},
		},
		{
			name:      "realm_name not in request",
			reqURI:    "?item_name=алмаз&region=eu",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Invalid query params. Err: realm_name is empty",
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
		{
			name:      "Valid request but SearchItem raises error",
			reqURI:    "?item_name=riseError&region=eu&realm_name=Killrog",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "error making get auction request, status: 404",
			},
		},
		{
			name:      "Valid request but SearchItem raises error",
			reqURI:    "?item_name=гаррош&region=eu&realm_name=errRealm",
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "error making GetAuctionData request, status: 404",
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

func TestV1Handler_TestHandler(t *testing.T) {
	h := newV1handler()

	err := h.MakeBlizzAuth()
	assert.NoError(t, err)

	err = h.GetBlizzRealms()
	assert.NoError(t, err)
}
