package server

import (
	"context"
	"testing"

	"auctioneer/app/conf"
	"github.com/stretchr/testify/assert"
)

func TestSetup(t *testing.T) {
	cfg := new(conf.Config)
	cfg.LogLvl = "INFO"
	ctx := context.Background()

	_, err := Setup(ctx, cfg)
	assert.NoError(t, err)
}
