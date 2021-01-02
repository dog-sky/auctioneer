package blizz

type BlizzardToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type DetailedName struct {
	RuRU string `json:"ru_RU"`
	EnGB string `json:"en_GB"`
	EnUS string `json:"en_US"`
}

type realms struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type BlizzRealmsSearchResult struct {
	Realms []realms `json:"realms"`
}

type ItemData struct {
	Media struct {
		ID int `json:"id"`
	} `json:"media"`
	Name    DetailedName `json:"name"`
	ID      int          `json:"id"`
	Quality struct {
		Type string `json:"type"`
	} `json:"quality"`
}

type ItemResult struct {
	Results []struct {
		Data ItemData `json:"data"`
	} `json:"results"`
}

type AuctionsDetail struct {
	ID   int `json:"id"`
	Item struct {
		ID        int `json:"id"`
		Context   int `json:"context"`
		Modifiers []struct {
			Type  int `json:"type"`
			Value int `json:"value"`
		} `json:"modifiers"`
		PetBreedID   int `json:"pet_breed_id"`
		PetLevel     int `json:"pet_level"`
		PetQualityID int `json:"pet_quality_id"`
		PetSpeciesID int `json:"pet_species_id"`
	} `json:"item"`
	Buyout   int    `json:"buyout"`
	Quantity int    `json:"quantity"`
	TimeLeft string `json:"time_left"`
}

type AuctionData struct {
	Auctions []*AuctionsDetail `json:"auctions"`
}
