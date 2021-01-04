package system_test

import (
	"net/http/httptest"
	"testing"

	"auctioneer/app/api/system"
	server "auctioneer/app/auctioneer"
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	"github.com/stretchr/testify/assert"
)

func Test_Ping(t *testing.T) {
	cfg, _ := conf.NewConfig()
	logger, _ := logging.NewLogger("DEBUG")
	app := server.NewApp(logger, cfg)
	h := system.NewSystemHandler()

	app.Fib.Get("/ping", h.Ping)
	req := httptest.NewRequest("GET", "/ping", nil)

	resp, err := app.Fib.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
}
