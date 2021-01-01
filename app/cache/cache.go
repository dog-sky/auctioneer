package cache

import (
	"strings"
	"sync"
)

type Cache interface {
	GetRealmID(string) int
	SetRealmID(string, int)
}

type cache struct {
	mux       sync.RWMutex
	realmList map[string]int
}

func NewCache() Cache {
	return &cache{
		realmList: make(map[string]int),
	}
}

func (c *cache) GetRealmID(RealmName string) int {
	if len(RealmName) == 0 {
		return 0
	}

	c.mux.RLock()
	defer c.mux.RUnlock()

	if item, ok := c.realmList[strings.ToLower(RealmName)]; ok {
		return item
	}
	return 0
}

func (c *cache) SetRealmID(RealmName string, RealmID int) {
	if len(RealmName) == 0 {
		return
	}

	c.mux.Lock()
	c.realmList[strings.ToLower(RealmName)] = RealmID
	c.mux.Unlock()
}
