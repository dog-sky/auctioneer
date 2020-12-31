package v1

type ResponseV1 struct {
	Success bool   `json:"success,omitempty"`
	Result  string `json:"result,omitempty"`
}

type BlizzardToken struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
}

type Realms struct {
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
	Realms []Realms `json:"realms"`
}
