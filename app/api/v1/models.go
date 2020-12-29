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
