package blizz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)


func TestClient_searchItemErrJson(t *testing.T) {
	errClient := makeTestBlizzClient()
	_ = errClient.MakeBlizzAuth()

	res, err := errClient.SearchItem("error_item_search", "eu")
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = errClient.SearchItem("error_item_search", "gr")
	assert.Error(t, err)
	assert.Nil(t, res)
}


func TestClient_searchItem(t *testing.T) {
	blizzClient := makeTestBlizzClient()
	_ = blizzClient.MakeBlizzAuth()

	res, err := blizzClient.SearchItem("Гаррош", "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = blizzClient.SearchItem("Garrosh", "eu")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = blizzClient.SearchItem("Garrosh", "us")
	assert.Error(t, err)
	assert.Nil(t, res)
}
