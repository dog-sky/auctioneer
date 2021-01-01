package v1

type searchQueryParams struct {
	RealmName string `query:"realm_name"`
	ItemName  string `query:"item_name"`
	Region    string `query:"region"`
}

type responseResultv1 struct {
	ItemID          int    `json:"item_id"`
	ItemName        string `json:"item_name"`
	ItemMedia       string `json:"item_media"`
	ItemAuctionData string `json:"item_auc_data"` // будет структура
}

type ResponseV1 struct {
	Success bool               `json:"success"`
	Message string             `json:"message,omitempty"`
	Result  []responseResultv1 `json:"result,omitempty"`
}
