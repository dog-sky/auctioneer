package blizz

type BlizzardToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type realms struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

type BlizzRealmsSearchResult struct {
	Realms []realms `json:"realms"`
}

type ItemName struct {
	ItIT string `json:"it_IT"`
	RuRU string `json:"ru_RU"`
	EnGB string `json:"en_GB"`
	ZhTW string `json:"zh_TW"`
	KoKR string `json:"ko_KR"`
	EnUS string `json:"en_US"`
	EsMX string `json:"es_MX"`
	PtBR string `json:"pt_BR"`
	EsES string `json:"es_ES"`
	ZhCN string `json:"zh_CN"`
	FrFR string `json:"fr_FR"`
	DeDE string `json:"de_DE"`
}

type ItemData struct {
	Media struct {
		ID int `json:"id"`
	} `json:"media"`
	Name ItemName `json:"name"`
	ID   int      `json:"id"`
}

type ItemResult struct {
	Results []struct {
		Data ItemData `json:"data"`
	} `json:"results"`
}
