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
	v1 "github.com/dog-sky/auctioneer/internal/api/v1"
	server "github.com/dog-sky/auctioneer/internal/app/auctioneer"
	"github.com/dog-sky/auctioneer/internal/client/blizz"
	blizzMock "github.com/dog-sky/auctioneer/internal/client/blizz/mocks"
	"github.com/dog-sky/auctioneer/internal/router"

	"github.com/stretchr/testify/assert"
)

func Test_SearchItemMedia(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		itemID    int
		init      func(t *testing.T) *server.Auctioneer
		expStatus int
		exp       interface{}
	}{
		{
			name:      "OK request",
			itemID:    200,
			expStatus: 200,
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetItemMediaMock.Expect("200").Return(&blizz.ItemMedia{
					ID: 200,
					Assets: []blizz.ItemAssets{
						{
							Key:        "hello",
							Value:      "world",
							FileDataID: 100,
						},
					},
				}, nil)

				app.BaseHandler = api.NewBasehandler(mockClient)
				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			exp: v1.ResponseV1ItemMedia{
				Success: true,
				ItemMedia: &blizz.ItemMedia{
					Assets: []blizz.ItemAssets{
						{
							Key:        "hello",
							Value:      "world",
							FileDataID: 100,
						},
					},
					ID: 200,
				},
			},
		},
		{
			name:      "404 request",
			itemID:    404,
			expStatus: 404,
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetItemMediaMock.Expect("404").Return(nil, nil)

				app.BaseHandler = api.NewBasehandler(mockClient)
				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			exp: v1.ResponseV1{
				Success: false,
				Message: "Item with ID 404 not found",
			},
		},
		{
			name:      "400 request",
			itemID:    400,
			expStatus: 400,
			init: func(t *testing.T) *server.Auctioneer {
				cfg := new(conf.Config)
				cfg.LogLvl = "INFO"
				ctx := context.Background()
				app, err := server.NewApp(ctx, cfg)
				assert.NoError(t, err)

				mockClient := blizzMock.NewClientMock(t)
				mockClient.GetItemMediaMock.Expect("400").Return(
					nil,
					fmt.Errorf("error making get item mdeia request, status: %d", 400),
				)

				app.BaseHandler = api.NewBasehandler(mockClient)
				router.SetupRoutes(app.Fib, app.BaseHandler)

				return app
			},
			exp: v1.ResponseV1{
				Success: false,
				Message: "error making get item mdeia request, status: 400",
			},
		},
	}

	for _, tc := range testCases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			app := tc.init(t)

			reqURI := fmt.Sprintf("/api/v1/item_media/%d", tc.itemID)
			req := httptest.NewRequest("GET", reqURI, nil)

			resp, err := app.Fib.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tc.expStatus, resp.StatusCode)

			bodyBytes, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)

			if tc.expStatus == 200 {
				var respV1 v1.ResponseV1ItemMedia
				err = json.Unmarshal(bodyBytes, &respV1)
				assert.NoError(t, err)

				assert.EqualValues(t, tc.exp, respV1)

				return
			}

			var respV1 v1.ResponseV1
			err = json.Unmarshal(bodyBytes, &respV1)
			assert.NoError(t, err)

			assert.EqualValues(t, tc.exp, respV1)
		})
	}
}
