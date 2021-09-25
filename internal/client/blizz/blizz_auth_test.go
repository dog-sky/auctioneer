package blizz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_auth(t *testing.T) {
	blizzClient := makeTestBlizzClient()
	err := blizzClient.MakeBlizzAuth()
	assert.NoError(t, err)
}
