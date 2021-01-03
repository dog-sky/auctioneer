package blizz

import (
	"auctioneer/app/cache"
	"auctioneer/app/conf"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/twinj/uuid"
	"net/http"
	"net/http/httptest"
	"testing"
	"strings"
)

var blizzClient Client

func init() {
	srv := serverMock()
	// defer srv.Close()

	cache := cache.NewCache()
	cfg := conf.Config{
		BlizzApiCfg: conf.BlizzApiCfg{
			EuAPIUrl:     srv.URL,
			UsAPIUrl:     srv.URL,
			AUTHUrl:      srv.URL + "/oauth/token",
			ClientSecret: "secret",
		},
	}

	blizzClient = NewClient(&cfg.BlizzApiCfg, cache)
}

func TestClient_auth(t *testing.T) {
	err := blizzClient.MakeBlizzAuth()
	assert.NoError(t, err)
}

func TestClient_getRealms(t *testing.T) {
	err := blizzClient.GetBlizzRealms()
	assert.NoError(t, err)
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

func TestClient_getAuctionData(t *testing.T) {
	res, err := blizzClient.GetAuctionData(501, "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = blizzClient.GetAuctionData(501, "us")
	assert.Error(t, err)
	assert.Nil(t, res)
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
			},
		},
	}

	w.Header().Set("last-modified", "Sat, 2 Jan 2021 12:08:43 GMT")
	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	_ = json.NewEncoder(w).Encode(aucData)
}
