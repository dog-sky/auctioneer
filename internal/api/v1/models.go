package v1

import (
	"github.com/dog-sky/auctioneer/internal/client/blizz"
)

type searchQueryParams struct {
	RealmName string `query:"realm_name,required"`
	ItemName  string `query:"item_name,required"`
	Region    string `query:"region,required"`
}

type ResponseV1 struct {
	Success bool                    `json:"success"`
	Message string                  `json:"message,omitempty"`
	Result  []*blizz.AuctionsDetail `json:"result"`
}

type ResponseV1ItemMedia struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
	*blizz.ItemMedia
}
