package cache

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/karlseguin/ccache/v2"
)

type Cache interface {
	GetRealmID(string) int
	SetRealmID(string, int)
	GetAuctionData(realmID int, region string) interface{}
	SetAuctionData(realmID int, region string, auctionData interface{}, updatedAt *time.Time)
}

type cache struct {
	mux       sync.RWMutex
	realmList map[string]int
	cache     *ccache.Cache
}

func NewCache() Cache {
	return &cache{
		realmList: make(map[string]int),
		cache:     ccache.New(ccache.Configure().MaxSize(1000).ItemsToPrune(100)),
	}
}

func (c *cache) GetRealmID(RealmName string) int {
	c.mux.RLock()


	return c.realmList[strings.ToLower(RealmName)]
}

func (c *cache) SetRealmID(RealmName string, RealmID int) {
	if RealmName == "" {
		return
	}

	c.mux.Lock()
	defer c.mux.RUnlock()

	c.realmList[strings.ToLower(RealmName)] = RealmID
}

func (c *cache) GetAuctionData(realmID int, region string) interface{} {
	if realmID == 0 || region == "" {
		return nil
	}

	key := fmt.Sprintf("%d_%s", realmID, region)
	if item := c.cache.Get(key); item != nil {
		return item.Value()
	}

	return nil
}

func (c *cache) SetAuctionData(realmID int, region string, auctionData interface{}, updatedAt *time.Time) {
	key := fmt.Sprintf("%d_%s", realmID, region)
	ttl := time.Until(updatedAt.Add(time.Hour))

	c.cache.Set(key, auctionData, ttl)
}
