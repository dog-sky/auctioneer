package v1_test

import (
	"fmt"
	"strings"

	"auctioneer/app/api/system"
	"auctioneer/app/api/v1"
	"auctioneer/app/blizz"
)

func newV1handler() v1.Handler {
	return v1.NewBasehandlerv1(&mockBlizzClient{})
}

func newSystemHandler() system.Handler {
	return system.NewSystemHandler()
}

type mockHandler struct {
	v1     v1.Handler
	system system.Handler
}

func (h *mockHandler) V1MakeBlizzAuth() error { return nil }

func (h *mockHandler) V1GetBlizzRealms() error { return nil }

func (h *mockHandler) V1Handler() v1.Handler {
	return h.v1
}

func (h *mockHandler) SystemHandler() system.Handler {
	return h.system
}

type mockBlizzClient struct{}

func (c *mockBlizzClient) BlizzAuthRoutine() {}

func (c *mockBlizzClient) GetBlizzRealms() error {
	return nil
}

func (c *mockBlizzClient) MakeBlizzAuth() error {
	return nil
}

func (c *mockBlizzClient) GetItemMedia(itemID string) (*blizz.ItemMedia, error) {

	if itemID == "404" {
		return nil, nil
	}

	if itemID == "400" {
		return nil, fmt.Errorf("error making get item mdeia request, status: %d", 400)
	}

	res := blizz.ItemMedia{
		ID: 200,
		Assets: []blizz.ItemAssets{
			blizz.ItemAssets{
				Key:        "hello",
				Value:      "world",
				FileDataID: 100,
			},
		},
	}

	return &res, nil
}

func (c *mockBlizzClient) GetRealmID(s string) int {
	if s == "Killrog" {
		return 1
	}
	if s == "errRealm" {
		return 2
	}
	return 0
}

func (c *mockBlizzClient) SearchItem(itemName string, region string) (*blizz.ItemResult, error) {
	if strings.Contains(itemName, "гаррош") {
		return &blizz.ItemResult{
			Results: []blizz.ItemTesult{
				{
					Data: blizz.ItemData{
						Name: blizz.DetailedName{
							RuRU: "Оправдание Гарроша",
							EnGB: "Garrosh's Pardon",
							EnUS: "Garrosh's Pardon",
						},
						ID: 1,
						Quality: blizz.ItemQuality{
							Type: "EPIC",
						},
					},
				},
			},
		}, nil
	}
	if strings.Contains(itemName, "опал") {
		return &blizz.ItemResult{
			Results: []blizz.ItemTesult{
				{
					Data: blizz.ItemData{
						Name: blizz.DetailedName{
							RuRU: "Большой опал",
							EnGB: "Large Opal",
							EnUS: "Large Opal",
						},
						ID: 2,
						Quality: blizz.ItemQuality{
							Type: "UNCOMMON",
						},
					},
				},
			},
		}, nil
	}
	if strings.Contains(itemName, "riseError") {
		return nil, fmt.Errorf(
			"error making get auction request, status: %d", 404,
		)
	}
	return &blizz.ItemResult{
		Results: []blizz.ItemTesult{},
	}, nil
}

func (c *mockBlizzClient) GetAuctionData(realmID int, region string) ([]*blizz.AuctionsDetail, error) {
	if realmID == 1 {
		return []*blizz.AuctionsDetail{
			&blizz.AuctionsDetail{
				ID: 1,
				Item: blizz.AcuItem{
					ID:      1,
					Context: 1,
					Modifiers: []blizz.AucItemModifiers{
						blizz.AucItemModifiers{
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
				ItemName: blizz.DetailedName{
					RuRU: "Оправдание Гарроша",
					EnGB: "Garrosh's Pardon",
					EnUS: "Garrosh's Pardon",
				},
				Quality: "EPIC",
			},
		}, nil
	}
	return nil, fmt.Errorf(
		"error making GetAuctionData request, status: %d", 404,
	)
}
