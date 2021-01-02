package v1

import (
	"auctioneer/app/blizz"
)

type searchQueryParams struct {
	RealmName string `query:"realm_name"`
	ItemName  string `query:"item_name"`
	Region    string `query:"region"`
}

type ResponseV1 struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message,omitempty"`
	Result  []*blizz.AuctionsDetail `json:"result,omitempty"`
}
