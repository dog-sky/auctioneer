package blizz

import (
	"auctioneer/app/conf"
	logging "auctioneer/app/logger"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_getRealms(t *testing.T) {
	blizzClient := makeTestBlizzClient()
	_ = blizzClient.MakeBlizzAuth()
	c := blizzClient.(*client)

	err := c.GetBlizzRealms()
	assert.NoError(t, err)

	// второй раз для получения из кэша
	err = c.GetBlizzRealms()
	assert.NoError(t, err)
}

func TestClient_getRealmsErr(t *testing.T) {
	srv := serverMock()
	blizzCfg := conf.BlizzApiCfg{
		EuAPIUrl:     srv.URL,
		UsAPIUrl:     srv.URL,
		AUTHUrl:      srv.URL + "/oauth/token",
		ClientSecret: "secret",
		RegionList:   []string{"gb"},
	}
	cfgErr := &conf.Config{
		BlizzApiCfg: blizzCfg,
	}

	log, _ := logging.NewLogger("ERROR")
	errClient := NewClient(log, &cfgErr.BlizzApiCfg)
	_ = errClient.MakeBlizzAuth()

	err := errClient.GetBlizzRealms()
	assert.Error(t, err)
}
