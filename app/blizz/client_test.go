package blizz

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"

	"auctioneer/app/conf"
	logging "auctioneer/app/logger"

	"github.com/twinj/uuid"
)

func serverMock() *httptest.Server {
	handler := http.NewServeMux()
	handler.HandleFunc("/oauth/token", authMock)
	handler.HandleFunc("/data/wow/realm/index", realmListMock)
	handler.HandleFunc("/data/wow/search/item", searchItemMock)
	handler.HandleFunc("/data/wow/connected-realm/", auctionDataMock)
	handler.HandleFunc("/data/wow/media/item/", itemMediaMock)

	srv := httptest.NewServer(handler)

	return srv
}

func makeTestBlizzClient() Client {
	srv := serverMock()
	cfg := conf.Config{
		BlizzApiCfg: conf.BlizzApiCfg{
			EuAPIUrl:     srv.URL,
			UsAPIUrl:     srv.URL,
			AUTHUrl:      srv.URL + "/oauth/token",
			ClientSecret: "secret",
			RegionList:   []string{"eu", "us"},
		},
	}

	log, _ := logging.NewLogger("ERROR")
	return NewClient(log, &cfg.BlizzApiCfg)
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
		Results: []BlizzRealmsSearchResultResults{
			{
				Data: BlizzRealmsSearchResultResultsData{
					ID: 501,
					Realms: []BlizzRealmsSearchResultResultsDataRealms{
						{
							Name: BlizzRealmsSearchResultResultsDataRealmsName{
								RuRU: "Arathor",
								EnGB: "Arathor",
							},
						},
					},
				},
			},
			{
				Data: BlizzRealmsSearchResultResultsData{
					ID: 500,
					Realms: []BlizzRealmsSearchResultResultsDataRealms{
						{
							Name: BlizzRealmsSearchResultResultsDataRealmsName{
								RuRU: "Aggramar",
								EnGB: "Aggramar",
							},
						},
					},
				},
			},
			{
				Data: BlizzRealmsSearchResultResultsData{
					ID: 503,
					Realms: []BlizzRealmsSearchResultResultsDataRealms{
						{
							Name: BlizzRealmsSearchResultResultsDataRealmsName{
								RuRU: "WhronJson",
								EnGB: "WhronJson",
							},
						},
					},
				},
			},
			{
				Data: BlizzRealmsSearchResultResultsData{
					ID: 504,
					Realms: []BlizzRealmsSearchResultResultsDataRealms{
						{
							Name: BlizzRealmsSearchResultResultsDataRealmsName{
								RuRU: "TimeDecodeErr",
								EnGB: "TimeDecodeErr",
							},
						},
					},
				},
			},
			{
				Data: BlizzRealmsSearchResultResultsData{
					ID: 502,
					Realms: []BlizzRealmsSearchResultResultsDataRealms{
						{
							Name: BlizzRealmsSearchResultResultsDataRealmsName{
								RuRU: "ServerStatus",
								EnGB: "ServerStatus",
							},
						},
					},
				},
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

	items := &Item{
		Results: []ItemResultResults{
			{
				Data: ItemResultResultsData{
					Name: ItemResultResultsDataName{
						RuRU: "Оправдание Гарроша",
						EnGB: "Garrosh's Pardon",
						EnUS: "Garrosh's Pardon",
					},
					ID: 1,
					Quality: ItemResultResultsDataQuality{
						Type: "EPIC",
					},
				},
			},
			{
				Data: ItemResultResultsData{
					Name: ItemResultResultsDataName{
						RuRU: "Боевой топор авангарда Гарроша",
						EnGB: "Garrosh's Vanguard Battleaxe",
						EnUS: "Garrosh's Vanguard Battleaxe",
					},
					ID: 2,
					Quality: ItemResultResultsDataQuality{
						Type: "UNCOMMON",
					},
				},
			},
		},
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(items)
}

func itemMediaMock(w http.ResponseWriter, r *http.Request) {
	q := r.RequestURI
	if strings.Contains(q, "/502") {
		w.WriteHeader(404)
		return
	}

	if strings.Contains(q, "/504") {
		_, _ = io.WriteString(w, "hello")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		return
	}

	itemMedia := ItemMedia{
		Assets: []ItemAssets{
			{
				Key:        "icon",
				Value:      "https://render-eu.worldofwarcraft.com/icons/56/inv_sword_39.jpg",
				FileDataID: 135349,
			},
		},
		ID: 19019,
	}

	w.Header().Set("Content-Type", "application/json;charset=UTF-8")
	_ = json.NewEncoder(w).Encode(itemMedia)
}

func auctionDataMock(w http.ResponseWriter, r *http.Request) {
	q := r.RequestURI
	if strings.Contains(q, "/502/") {
		w.WriteHeader(404)
		return
	}

	if strings.Contains(q, "/504/") {
		aucData := AuctionData{
			Auctions: []*AuctionsDetail{
				{
					ID: 1,
					Item: AcuItem{
						ID:      2,
						Context: 1,
						Modifiers: []AucItemModifiers{
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
					ItemName: ItemResultResultsDataName{
						RuRU: "Боевой топор авангарда Гарроша",
						EnGB: "Garrosh's Vanguard Battleaxe",
						EnUS: "Garrosh's Vanguard Battleaxe",
					},
					Quality: "UNCOMMON",
					Price:   120000,
				},
			},
		}

		w.Header().Set("last-modified", "11/11/2020")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		_ = json.NewEncoder(w).Encode(aucData)
		return
	}

	if strings.Contains(q, "/503/") {
		w.Header().Set("last-modified", "Sat, 2 Jan 2021 12:08:43 GMT")
		_, _ = io.WriteString(w, "hello")
		return
	}

	aucData := AuctionData{
		Auctions: []*AuctionsDetail{
			{
				ID: 1,
				Item: AcuItem{
					ID:      2,
					Context: 1,
					Modifiers: []AucItemModifiers{
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
				ItemName: ItemResultResultsDataName{
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
