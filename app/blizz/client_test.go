package blizz

import (
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/twinj/uuid"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var blizzClient Client

func init() {
	srv := serverMock()
	cache := cache.NewCache()
	cfg := conf.Config{
		BlizzApiCfg: conf.BlizzApiCfg{
			EuAPIUrl:     srv.URL,
			UsAPIUrl:     srv.URL,
			AUTHUrl:      srv.URL + "/oauth/token",
			ClientSecret: "secret",
			RegionList:   []string{"eu", "us"},
		},
	}

	blizzClient = NewClient(&cfg.BlizzApiCfg, cache)
}

func TestClient_auth(t *testing.T) {
	err := blizzClient.MakeBlizzAuth()
	assert.NoError(t, err)
}

func TestClient_getRealms(t *testing.T) {
	c := blizzClient.(*client)
	c.cfg.RegionList = append(c.cfg.RegionList, "gb")

	err := c.GetBlizzRealms()
	assert.Error(t, err)
}

func TestClient_getRealmsErr(t *testing.T) {
	srv := serverMock()
	cache := cache.NewCache()
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

	errClient := NewClient(&cfgErr.BlizzApiCfg, cache)
	_ = errClient.MakeBlizzAuth()

	err := errClient.GetBlizzRealms()
	assert.Error(t, err)
}

func TestClient_searchItem(t *testing.T) {
	res, err := blizzClient.SearchItem("Гаррош", "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = blizzClient.SearchItem("Garrosh", "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = blizzClient.SearchItem("Garrosh", "us")
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestClient_searchItemErrJson(t *testing.T) {
	srv := serverMock()
	cache := cache.NewCache()
	blizzCfg := conf.BlizzApiCfg{
		EuAPIUrl:     srv.URL,
		UsAPIUrl:     srv.URL,
		AUTHUrl:      srv.URL + "/oauth/token",
		ClientSecret: "secret",
		RegionList:   []string{"eu", "us"},
	}
	cfgErr := &conf.Config{
		BlizzApiCfg: blizzCfg,
	}

	errClient := NewClient(&cfgErr.BlizzApiCfg, cache)
	_ = errClient.MakeBlizzAuth()

	res, err := errClient.SearchItem("error_item_search", "eu")
	println("ERROR: ", err.Error())
	assert.Error(t, err)
	assert.Nil(t, res)
}

func TestClient_getAuctionData(t *testing.T) {
	res, err := blizzClient.GetAuctionData(501, "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestClient_getAuctionDataError(t *testing.T) {
	tests := []struct {
		name   string
		region string
	}{
		{
			name:   "Server status Err",
			region: "us",
		}, {
			name:   "Time Decode Err",
			region: "wrong_time",
		}, {
			name:   "JSON Decode Err",
			region: "wrong_json",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := blizzClient.GetAuctionData(501, "wrong_json")
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	}
}

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/oauth/token", authMock)
	handler.HandleFunc("/data/wow/realm/index", realmListMock)
	handler.HandleFunc("/data/wow/search/item", searchItemMock)
	handler.HandleFunc("/data/wow/connected-realm/501/auctions", auctionDataMock)

	srv := httptest.NewServer(handler)

	return srv
}

func authMock(w http.ResponseWriter, r *http.Request) {
	token := &BlizzardToken{
		AccessToken: uuid.NewV4().String(),
		TokenType:   "bearer",
		ExpiresIn:   86399,
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}

func realmListMock(w http.ResponseWriter, r *http.Request) {
	q := r.RequestURI
	if strings.Contains(q, "dynamic-gb") {
		w.WriteHeader(404)
		return
	}

	rlms := BlizzRealmsSearchResult{
		Realms: []realm{
			realm{
				ID:   501,
				Name: "Arathor",
			},
			realm{
				ID:   500,
				Name: "Aggramar",
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(rlms)
}

func searchItemMock(w http.ResponseWriter, r *http.Request) {
	q := r.RequestURI
	if strings.Contains(q, "static-us") {
		w.WriteHeader(404)
		return
	}

	if strings.Contains(q, "error_item_search") {
		_, _ = io.WriteString(w, "{hello, there}")
		return
	}

	items := &ItemResult{
		Results: []ItemTesult{
			{
				Data: ItemData{
					Media: ItemMedia{
						ID: 1,
					},
					Name: DetailedName{
						RuRU: "Оправдание Гарроша",
						EnGB: "Garrosh's Pardon",
						EnUS: "Garrosh's Pardon",
					},
					ID: 1,
					Quality: ItemQuality{
						Type: "EPIC",
					},
				},
			},
			{
				Data: ItemData{
					Media: ItemMedia{
						ID: 2,
					},
					Name: DetailedName{
						RuRU: "Боевой топор авангарда Гарроша",
						EnGB: "Garrosh's Vanguard Battleaxe",
						EnUS: "Garrosh's Vanguard Battleaxe",
					},
					ID: 2,
					Quality: ItemQuality{
						Type: "UNCOMMON",
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func auctionDataMock(w http.ResponseWriter, r *http.Request) {
	q := r.RequestURI
	if strings.Contains(q, "dynamic-us") {
		w.WriteHeader(404)
		return
	}

	if strings.Contains(q, "dynamic-wrong_time") {
		w.Header().Set("last-modified", "Not a time")
		return
	}

	if strings.Contains(q, "dynamic-wrong_json") {
		w.Header().Set("last-modified", "Sat, 2 Jan 2021 12:08:43 GMT")
		_, _ = io.WriteString(w, "hello")
		return
	}

	aucData := AuctionData{
		Auctions: []*AuctionsDetail{
			&AuctionsDetail{
				ID: 1,
				Item: AcuItem{
					ID:      2,
					Context: 1,
					Modifiers: []AucItemModifiers{
						AucItemModifiers{
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
				ItemName: DetailedName{
					RuRU: "Боевой топор авангарда Гарроша",
					EnGB: "Garrosh's Vanguard Battleaxe",
					EnUS: "Garrosh's Vanguard Battleaxe",
				},
				Quality: "UNCOMMON",
				Price:   120000,
			},
		},
	}

	w.Header().Set("last-modified", "Sat, 2 Jan 2021 12:08:43 GMT")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	_ = json.NewEncoder(w).Encode(aucData)
}
