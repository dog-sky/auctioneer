package system_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"auctioneer/app/api/system"
	server "auctioneer/app/auctioneer"
	"auctioneer/app/conf"

	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	cfg := new(conf.Config)
	cfg.LogLvl = "INFO"
	ctx := context.Background()
	app, err := server.NewApp(ctx, cfg)
	assert.NoError(t, err)
	h := system.NewSystemHandler()

	app.Fib.Get("/ping", h.Ping)
	req := httptest.NewRequest("GET", "/ping", nil)

	resp, err := app.Fib.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
