package blizz

import (
	"fmt"
	"time"

	"github.com/levigross/grequests"
)

type ItemAssets struct {
	Key        string `json:"key"`
	Value      string `json:"value"`
	FileDataID int    `json:"file_data_id"`
}

type ItemMedia struct {
	Assets []ItemAssets `json:"assets"`
	ID     int          `json:"id"`
}

type AucItemModifiers struct {
	Type  int `json:"type"`
	Value int `json:"value"`
}

type AcuItem struct {
	ID           int                `json:"id"`
	Context      int                `json:"context"`
	Modifiers    []AucItemModifiers `json:"modifiers"`
	PetBreedID   int                `json:"pet_breed_id"`
	PetLevel     int                `json:"pet_level"`
	PetQualityID int                `json:"pet_quality_id"`
	PetSpeciesID int                `json:"pet_species_id"`
}

type AuctionsDetail struct {
	ID       int                       `json:"id"`
	Item     AcuItem                   `json:"item"`
	Buyout   int                       `json:"buyout"`
	Quantity int                       `json:"quantity"`
	TimeLeft string                    `json:"time_left"`
	ItemName ItemResultResultsDataName `json:"item_name"`
	Quality  string                    `json:"quality"`
	Price    int                       `json:"unit_price"`
}

type AuctionData struct {
	Auctions []*AuctionsDetail `json:"auctions"`
}

func (c *client) getCachedAuctionData(realmID int, region string) []*AuctionsDetail {
	data := c.cache.GetAuctionData(realmID, region)

	if t, ok := data.(*AuctionData); ok {
		return t.Auctions
	}

	return nil
}

func (c *client) GetAuctionData(realmID int, region string) ([]*AuctionsDetail, error) {
	// Аукцион по реалму обновляется раз в час. В заголовке приходит дата обновления
	// last-modified: Thu, 31 Dec 2020 15:08:43 GMT
	// нужно сохранить данные локально для реалма и отдавать их из кеша в течение часа

	if data := c.getCachedAuctionData(realmID, region); data != nil {
		return data, nil
	}

	requestURL := c.urls[region] + fmt.Sprintf("/data/wow/connected-realm/%d/auctions", realmID)
	ro := &grequests.RequestOptions{
		Params: map[string]string{
			"namespace":    fmt.Sprintf("dynamic-%s", region),
			"access_token": c.token.AccessToken,
		},
	}
	response, err := c.makeGetRequest(requestURL, ro)
	if err != nil {
		return nil, fmt.Errorf("err making AUCTION DATA request: %v", err)
	}

	auctionData := new(AuctionData)
	if err := response.JSON(auctionData); err != nil {
		return nil, fmt.Errorf(
			"error unmarshaling item media response: %v", err,
		)
	}

	updatedAt := response.Header.Get("last-modified") // GMT! (-3)
	updatedAtParsed, err := time.Parse(layoutUS, updatedAt)
	if err != nil {
		return nil, fmt.Errorf(
			"error parsing last-modified header in auction response: %v", err,
		)
	}
	c.cache.SetAuctionData(realmID, region, auctionData, &updatedAtParsed)

	return auctionData.Auctions, nil
}
