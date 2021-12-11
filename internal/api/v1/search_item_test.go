package v1_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http/httptest"
	"testing"

	conf "github.com/dog-sky/auctioneer/configs"
	"github.com/dog-sky/auctioneer/internal/api"
	server "github.com/dog-sky/auctioneer/internal/app/auctioneer"
	"github.com/dog-sky/auctioneer/internal/client/blizz"
	blizzMock "github.com/dog-sky/auctioneer/internal/client/blizz/mocks"
	"github.com/dog-sky/auctioneer/internal/router"

	v1 "github.com/dog-sky/auctioneer/internal/api/v1"

	"github.com/stretchr/testify/assert"
)

func Test_SearchItemData(t *testing.T) {
	t.Parallel()

	aucDetail := []*blizz.AuctionsDetail{{
		ID: 1,
		Item: blizz.AcuItem{
			ID:      1,
			Context: 1,
			Modifiers: []blizz.AucItemModifiers{
				{
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
		ItemName: blizz.ItemResultResultsDataName{
			RuRU: "Оправдание Гарроша",
			EnGB: "Garrosh's Pardon",
			EnUS: "Garrosh's Pardon",
		},
		Quality: "EPIC",
	}}

	testCases := []struct {
		name      string
		reqURI    string
		init      func(t *testing.T) *server.Auctioneer
		expStatus int
		exp       v1.ResponseV1
	}{
		{
			name:   "Valid request",
			reqURI: "?item_name=гаррош&region=eu&realm_name=Killrog",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Killrog").Return(1)
				mockClient.GetAuctionDataMock.Expect(1, "eu").Return(aucDetail, nil)
				mockClient.SearchItemMock.Expect("гаррош", "eu").Return(
					&blizz.Item{
						Results: []blizz.ItemResultResults{
							{
								Data: blizz.ItemResultResultsData{
									Name: blizz.ItemResultResultsDataName{
										RuRU: "Оправдание Гарроша",
										EnGB: "Garrosh's Pardon",
										EnUS: "Garrosh's Pardon",
									},
									ID: 1,
									Quality: blizz.ItemResultResultsDataQuality{
										Type: "EPIC",
									},
								},
							},
						},
					},
					nil,
				)

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 200,
			exp: v1.ResponseV1{
				Success: true,
				Result: []*blizz.AuctionsDetail{
					{
						ID: 1,
						Item: blizz.AcuItem{
							ID:      1,
							Context: 1,
							Modifiers: []blizz.AucItemModifiers{
								{
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
						ItemName: blizz.ItemResultResultsDataName{
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
			name:   "Valid request, but item not in auction right now",
			reqURI: "?item_name=опал&region=eu&realm_name=Killrog",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Killrog").Return(1)
				mockClient.GetAuctionDataMock.Expect(1, "eu").Return(aucDetail, nil)
				mockClient.SearchItemMock.Expect("опал", "eu").Return(&blizz.Item{
					Results: []blizz.ItemResultResults{
						{
							Data: blizz.ItemResultResultsData{
								Name: blizz.ItemResultResultsDataName{
									RuRU: "Большой опал",
									EnGB: "Large Opal",
									EnUS: "Large Opal",
								},
								ID: 2,
								Quality: blizz.ItemResultResultsDataQuality{
									Type: "UNCOMMON",
								},
							},
						},
					},
				}, nil)

				app.BaseHandler = api.NewBasehandler(mockClient)
				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 200,
			exp: v1.ResponseV1{
				Success: true,
				Result:  []*blizz.AuctionsDetail{},
			},
		},
		{
			name:   "Item not found",
			reqURI: "?item_name=алмаз&region=eu&realm_name=Killrog",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Killrog").Return(1)
				mockClient.GetAuctionDataMock.Expect(1, "eu").Return(aucDetail, nil)
				mockClient.SearchItemMock.Expect("алмаз", "eu").Return(&blizz.Item{
					Results: []blizz.ItemResultResults{},
				}, nil)

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 404,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Item алмаз not found",
			},
		},
		{
			name:   "item_name not in request",
			reqURI: "?region=eu&realm_name=Killrog",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Killrog").Return(1)
				mockClient.GetAuctionDataMock.Expect(1, "eu").Return(aucDetail, nil)

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Invalid query params. Err: item_name is empty",
			},
		},
		{
			name:   "region not in request",
			reqURI: "?item_name=алмаз&realm_name=Killrog",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Killrog").Return(1)
				mockClient.GetAuctionDataMock.Expect(1, "eu").Return(aucDetail, nil)

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Invalid query params. Err: region is empty",
			},
		},
		{
			name:   "realm_name not in request",
			reqURI: "?item_name=алмаз&region=eu",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Return(0)
				mockClient.GetAuctionDataMock.Expect(0, "eu").Return(nil, fmt.Errorf(
					"error making GetAuctionData request, status: %d", 404,
				))

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Invalid query params. Err: realm_name is empty",
			},
		},
		{
			name:   "Realm not found",
			reqURI: "?item_name=алмаз&region=eu&realm_name=Гордунни",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Гордунни").Return(0)
				mockClient.GetAuctionDataMock.Expect(0, "eu").Return(nil, fmt.Errorf(
					"error making GetAuctionData request, status: %d", 404,
				))

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 404,
			exp: v1.ResponseV1{
				Success: false,
				Message: "Realm Гордунни not found",
			},
		},
		{
			name:   "Valid request but SearchItem raises error",
			reqURI: "?item_name=riseError&region=eu&realm_name=Killrog",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("Killrog").Return(1)
				mockClient.GetAuctionDataMock.Expect(1, "eu").Return(aucDetail, nil)
				mockClient.SearchItemMock.Expect("riseError", "eu").Return(nil, fmt.Errorf(
					"error making get auction request, status: %d", 404,
				))

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 400,
			exp: v1.ResponseV1{
				Success: false,
				Message: "error making get auction request, status: 404",
			},
		},
		{
			name:   "Valid request but SearchItem raises error",
			reqURI: "?item_name=гаррош&region=eu&realm_name=errRealm",
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetRealmIDMock.Expect("errRealm").Return(2)
				mockClient.GetAuctionDataMock.Expect(2, "eu").Return(nil, fmt.Errorf(
					"error making GetAuctionData request, status: %d", 404,
				))
				mockClient.SearchItemMock.Expect("гаррош", "eu").Return(
					&blizz.Item{
						Results: []blizz.ItemResultResults{
							{
								Data: blizz.ItemResultResultsData{
									Name: blizz.ItemResultResultsDataName{
										RuRU: "Оправдание Гарроша",
										EnGB: "Garrosh's Pardon",
										EnUS: "Garrosh's Pardon",
									},
									ID: 1,
									Quality: blizz.ItemResultResultsDataQuality{
										Type: "EPIC",
									},
								},
							},
						},
					},
					nil,
				)

				app.BaseHandler = api.NewBasehandler(mockClient)

				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			expStatus: 500,
			exp: v1.ResponseV1{
				Success: false,
				Message: "error making GetAuctionData request, status: 404",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			reqURI := fmt.Sprintf("/api/v1/auc_search%s", tc.reqURI)
			req := httptest.NewRequest("GET", reqURI, nil)

			app := tc.init(t)

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

func Test_Handler(t *testing.T) {
	t.Parallel()

	mockBlizzClient := blizzMock.NewClientMock(t)
	mockBlizzClient.MakeBlizzAuthMock.Return(nil)
	mockBlizzClient.GetBlizzRealmsMock.Return(nil)

	h := v1.NewBasehandlerv1(mockBlizzClient)

	err := h.MakeBlizzAuth()
	assert.NoError(t, err)

	err = h.GetBlizzRealms()
	assert.NoError(t, err)
}
