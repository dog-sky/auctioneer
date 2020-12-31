package cache_test

import (
	"auctioneer/app/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	c := cache.NewCache()

	testCases := []struct {
		name string
		realmName string
		reaimID int
		getKey string
		exp int
	}{
		{
			name: "SET GET ok Гордунни",
			realmName: "Гордунни",
			reaimID: 12,
			getKey: "Гордунни",
			exp: 12,
		},
		{
			name: "SET GET no such key Гордунни",
			realmName: "Гордунни",
			reaimID: 12,
			getKey: "Страж Смерти",
			exp: 0,
		},
	}

	for _, tc := range testCases{
		t.Run(tc.name, func(t *testing.T) {
			c.SetRealmID(tc.realmName, tc.reaimID)
			val := c.GetRealmID(tc.getKey)
			assert.Equal(t, tc.exp, val)
		})
	}
}
