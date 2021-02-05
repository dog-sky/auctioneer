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

	app, err := NewApp(ctx, cfg)
	assert.NoError(t, err)
	app.Setup()
}
