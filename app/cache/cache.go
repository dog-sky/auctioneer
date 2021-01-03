package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"
)

type Cache interface {
	GetRealmID(string) int
	SetRealmID(string, int)
	GetAuctionData(realmID int, region string) interface{}
	SetAuctionData(realmID int, region string, auctionData interface{}, updatedAt *time.Time)
}

type aucDataWithTTL struct {
	ttl  *time.Time
	data interface{}
}

type cache struct {
	mux         sync.RWMutex
	realmList   map[string]int
	auctionData map[string]*aucDataWithTTL
}

func NewCache() Cache {
	return &cache{
		realmList:   make(map[string]int),
		auctionData: make(map[string]*aucDataWithTTL),
	}
}

func (c *cache) GetRealmID(RealmName string) int {
	if RealmName == "" {
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
	if RealmName == "" {
		return
	}

	c.mux.Lock()
	c.realmList[strings.ToLower(RealmName)] = RealmID
	c.mux.Unlock()
}

func (c *cache) GetAuctionData(realmID int, region string) interface{} {
	if realmID == 0 {
		return nil
	}
	if region == "" {
		return nil
	}
	key := fmt.Sprintf("%d_%s", realmID, region)

	c.mux.RLock()
	defer c.mux.RUnlock()

	data, ok := c.auctionData[key]
	if !ok {
		return nil
	}

	now := time.Now().In(data.ttl.Location())
	if data.ttl.After(now) {
		return data.data
	}

	return nil
}

func (c *cache) SetAuctionData(realmID int, region string, auctionData interface{}, updatedAt *time.Time) {
	key := fmt.Sprintf("%d_%s", realmID, region)
	ttl := updatedAt.Add(time.Hour) // Кеш действителен час
	data := aucDataWithTTL{
		ttl:  &ttl,
		data: auctionData,
	}

	c.mux.Lock()
	c.auctionData[key] = &data
	c.mux.Unlock()
}
