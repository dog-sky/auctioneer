package system_test

import (
	"context"
	"net/http/httptest"
	"testing"

	conf "github.com/dog-sky/auctioneer/configs"
	"github.com/dog-sky/auctioneer/internal/api/system"

	server "github.com/dog-sky/auctioneer/internal/app/auctioneer"

	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	t.Parallel()

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
