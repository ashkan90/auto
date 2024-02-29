package src

import (
	"github.com/ashkan90/auto/utils"
	"log"
)

type Cache struct {
	//cache *cache.Cache
	cache *utils.SyncMap
}

func NewCache() *Cache {
	return &Cache{
		cache: utils.NewSyncMap(),
	}
}

func (c *Cache) Get(key string) (any, bool) {
	log.Println("[Cache/SyncMap] get a value", key)
	return c.cache.Get(key)
}

func (c *Cache) Set(key string, val any) {
	log.Println("[Cache/SyncMap] set a value", key)
	c.cache.Add(key, val)
}

func (c *Cache) Delete(key string) {
	log.Println("[Cache/SyncMap] delete a value", key)
	c.cache.Delete(key)
}

func (c *Cache) Clone() *Cache {
	return NewCache()
}

func (c *Cache) Reset() {
	c.cache = utils.NewSyncMap()
}
