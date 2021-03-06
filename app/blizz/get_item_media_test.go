package blizz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_getItemMedia(t *testing.T) {
	blizzClient := makeTestBlizzClient()
	_ = blizzClient.MakeBlizzAuth()

	res, err := blizzClient.GetItemMedia("500")
	assert.NoError(t, err)
	assert.NotNil(t, res)

	res, err = blizzClient.GetItemMedia("504")
	assert.Error(t, err)
	assert.Nil(t, res)

	res, err = blizzClient.GetItemMedia("502")
	assert.Error(t, err)
	assert.Nil(t, res)
}
