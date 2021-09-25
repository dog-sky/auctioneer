package cache_test

import (
	"testing"
	"time"

	"github.com/dog-sky/auctioneer/internal/client/blizz"
	"github.com/dog-sky/auctioneer/internal/pkg/cache"

	"github.com/stretchr/testify/assert"
)

func Test_cache_getSetRealm(t *testing.T) {
	c := cache.NewCache()

	t.Parallel()

	type tc struct {
		name      string
		realmName string
		realmID   int
		getKey    string
		exp       int
	}

	testCases := []tc{
		{
			name:      "SET GET ok Гордунни",
			realmName: "Гордунни",
			realmID:   12,
			getKey:    "Гордунни",
			exp:       12,
		},
		{
			name:      "SET GET no such key Гордунни",
			realmName: "Гордунни",
			realmID:   12,
			getKey:    "Страж Смерти",
			exp:       0,
		},
		{
			name:      "SET GET ok Гордунни разный регистра",
			realmName: "ГордуННи",
			realmID:   12,
			getKey:    "Гордунни",
			exp:       12,
		},
		{
			name:      "GET не передано значение",
			realmName: "ГордуННи",
			realmID:   12,
			getKey:    "",
			exp:       0,
		},
		{
			name:      "SET не передано значение",
			realmName: "",
			realmID:   12,
			getKey:    "Годунни",
			exp:       0,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			c.SetRealmID(testCase.realmName, testCase.realmID)
			val := c.GetRealmID(testCase.getKey)
			assert.Equal(t, testCase.exp, val)
		})
	}
}

func Test_cache_SetAuctionData(t *testing.T) {
	t.Parallel()

	c := cache.NewCache()
	now := time.Now()
	pastHour := time.Now().Add(-1 * time.Hour)
	pastMinutes := time.Now().Add(-20 * time.Minute)

	type args struct {
		realmID     int
		region      string
		auctionData interface{}
		updatedAt   *time.Time
	}

	tests := []struct {
		name       string
		args       args
		getRealmID int
		getRegion  string
		exp        interface{}
	}{
		{
			name: "OK now",
			args: args{
				realmID: 504,
				region:  "eu",
				auctionData: &blizz.AuctionData{
					Auctions: []*blizz.AuctionsDetail{
						{
							ID: 1,
							Item: blizz.AcuItem{
								ID:      2,
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
								RuRU: "Боевой топор авангарда Гарроша",
								EnGB: "Garrosh's Vanguard Battleaxe",
								EnUS: "Garrosh's Vanguard Battleaxe",
							},
							Quality: "UNCOMMON",
						},
					},
				},
				updatedAt: &now,
			},
			getRealmID: 504,
			getRegion:  "eu",
			exp:        &blizz.AuctionData{},
		},
		{
			name: "No get region",
			args: args{
				realmID: 504,
				region:  "eu",
				auctionData: &blizz.AuctionData{
					Auctions: []*blizz.AuctionsDetail{
						{
							ID: 1,
							Item: blizz.AcuItem{
								ID:      2,
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
								RuRU: "Боевой топор авангарда Гарроша",
								EnGB: "Garrosh's Vanguard Battleaxe",
								EnUS: "Garrosh's Vanguard Battleaxe",
							},
							Quality: "UNCOMMON",
						},
					},
				},
				updatedAt: &now,
			},
			getRealmID: 504,
			getRegion:  "",
			exp:        nil,
		},
		{
			name: "No such data region",
			args: args{
				realmID: 504,
				region:  "eu",
				auctionData: &blizz.AuctionData{
					Auctions: []*blizz.AuctionsDetail{
						{
							ID: 1,
							Item: blizz.AcuItem{
								ID:      2,
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
								RuRU: "Боевой топор авангарда Гарроша",
								EnGB: "Garrosh's Vanguard Battleaxe",
								EnUS: "Garrosh's Vanguard Battleaxe",
							},
							Quality: "UNCOMMON",
						},
					},
				},
				updatedAt: &now,
			},
			getRealmID: 50423,
			getRegion:  "us",
			exp:        nil,
		},
		{
			name: "OK past",
			args: args{
				realmID: 504,
				region:  "eu",
				auctionData: &blizz.AuctionData{
					Auctions: []*blizz.AuctionsDetail{
						{
							ID: 2,
							Item: blizz.AcuItem{
								ID:      3,
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
								RuRU: "Боевой топор авангарда Гарроша",
								EnGB: "Garrosh's Vanguard Battleaxe",
								EnUS: "Garrosh's Vanguard Battleaxe",
							},
							Quality: "UNCOMMON",
						},
					},
				},
				updatedAt: &pastMinutes,
			},
			getRealmID: 504,
			getRegion:  "eu",
			exp:        &blizz.AuctionData{},
		},
		{
			name: "OK past hour",
			args: args{
				realmID: 504,
				region:  "eu",
				auctionData: &blizz.AuctionData{
					Auctions: []*blizz.AuctionsDetail{
						{
							ID: 2,
							Item: blizz.AcuItem{
								ID:      3,
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
								RuRU: "Боевой топор авангарда Гарроша",
								EnGB: "Garrosh's Vanguard Battleaxe",
								EnUS: "Garrosh's Vanguard Battleaxe",
							},
							Quality: "UNCOMMON",
						},
					},
				},
				updatedAt: &pastHour,
			},
			getRealmID: 2,
			getRegion:  "eu",
			exp:        nil,
		},
		{
			name: "OK realm id 0",
			args: args{
				realmID: 0,
				region:  "eu",
				auctionData: &blizz.AuctionData{
					Auctions: []*blizz.AuctionsDetail{
						{
							ID: 2,
							Item: blizz.AcuItem{
								ID:      3,
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
								RuRU: "Боевой топор авангарда Гарроша",
								EnGB: "Garrosh's Vanguard Battleaxe",
								EnUS: "Garrosh's Vanguard Battleaxe",
							},
							Quality: "UNCOMMON",
						},
					},
				},
				updatedAt: &now,
			},
			getRealmID: 0,
			getRegion:  "",
			exp:        nil,
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c.SetAuctionData(tt.args.realmID, tt.args.region, tt.args.auctionData, tt.args.updatedAt)

			data := c.GetAuctionData(tt.getRealmID, tt.getRegion)
			assert.IsType(t, tt.exp, data)
		})
	}
}
