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

type realm struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type BlizzRealmsSearchResult struct {
	Realms []realm `json:"realms"`
}

type ItemQuality struct {
	Type string `json:"type"`
}

type ItemMedia struct {
	ID int `json:"id"`
}

type ItemData struct {
	Media   ItemMedia    `json:"media"`
	Name    DetailedName `json:"name"`
	ID      int          `json:"id"`
	Quality ItemQuality  `json:"quality"`
}

type ItemTesult struct {
	Data ItemData `json:"data"`
}

type ItemResult struct {
	Results []ItemTesult `json:"results"`
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
	ID       int          `json:"id"`
	Item     AcuItem      `json:"item"`
	Buyout   int          `json:"buyout"`
	Quantity int          `json:"quantity"`
	TimeLeft string       `json:"time_left"`
	ItemName DetailedName `json:"item_name"`
	Quality  string       `json:"quality"`
	Price    int          `json:"unit_price"`
}

type AuctionData struct {
	Auctions []*AuctionsDetail `json:"auctions"`
}
