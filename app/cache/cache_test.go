package cache_test

import (
	"auctioneer/app/cache"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCache(t *testing.T) {
	c := cache.NewCache()

	testCases := []struct {
		name      string
		realmName string
		realmID   int
		getKey    string
		exp       int
	}{
		{
			name:      "SET GET ok Гордунни",
			realmName: "Гордунни",
			realmID:   12,
			getKey:    "Гордунни",
			exp:       12,
		},
		{
			name:      "SET GET no such key Гордунни",
			realmName: "Гордунни",
			realmID:   12,
			getKey:    "Страж Смерти",
			exp:       0,
		},
		{
			name:      "SET GET ok Гордунни разный регистра",
			realmName: "ГордуННи",
			realmID:   12,
			getKey:    "Гордунни",
			exp:       12,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			c.SetRealmID(tc.realmName, tc.realmID)
			val := c.GetRealmID(tc.getKey)
			assert.Equal(t, tc.exp, val)
		})
	}
}
