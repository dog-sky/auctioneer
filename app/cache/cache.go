package cache

import (
	"strings"
	"sync"
)

type Cache struct {
	mux       sync.RWMutex
	realmList map[string]int
}

func NewCache() *Cache {
	return &Cache{
		realmList: make(map[string]int),
	}
}

func (c *Cache) GetRealmID(RealmName string) int {
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

func (c *Cache) SetRealmID(RealmName string, RealmID int) {
	if len(RealmName) == 0 {
		return
	}

	c.mux.Lock()
	c.realmList[strings.ToLower(RealmName)] = RealmID
	c.mux.Unlock()
}
