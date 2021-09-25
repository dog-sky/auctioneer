package blizz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_getAuctionData(t *testing.T) {
	blizzClient := makeTestBlizzClient()
	_ = blizzClient.MakeBlizzAuth()

	res, err := blizzClient.GetAuctionData(501, "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	// Второй раз для получения данных из кэша и проверка на ошибку.
	res, err = blizzClient.GetAuctionData(501, "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)
}

func TestClient_getAuctionDataError(t *testing.T) {
	errClient := makeTestBlizzClient()
	_ = errClient.MakeBlizzAuth()

	tests := []struct {
		name   string
		server int
	}{
		{
			name:   "Server status Err",
			server: 502,
		}, {
			name:   "Time Decode Err",
			server: 503,
		}, {
			name:   "JSON Decode Err",
			server: 504,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, err := errClient.GetAuctionData(tt.server, "eu")
			assert.Error(t, err)
			assert.Nil(t, res)
		})
	}
}
