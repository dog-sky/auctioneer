package blizz

type BlizzardToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type realms struct {
	Key struct {
		Href string `json:"href"`
	} `json:"key"`
	Name string `json:"name"`
	ID   int    `json:"id"`
	Slug string `json:"slug"`
}

type BlizzRealmsSearchResult struct {
	Links struct {
		Self struct {
			Href string `json:"href"`
		} `json:"self"`
	} `json:"_links"`
	Realms []realms `json:"realms"`
}
