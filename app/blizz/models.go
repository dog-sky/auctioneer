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
