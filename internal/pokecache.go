package pokecache

import (
	"time"
	"sync"
)

type Cache struct {
	cached 		map[string]cacheEntry
	mutex 		sync.RWMutex
}

type cacheEntry struct {
	createdAt 	time.Time
	val			[]byte
}

func NewCache(Interval time.Duration) *Cache {
	cache := new(Cache)
	cache.cached = make(map[string]cacheEntry)
	go cache.reapLoop(5)
	return cache
}

func (cache *Cache) Add(key string, val []byte) {
	cache.mutex.Lock()
	cache.cached[key] = cacheEntry{createdAt: time.Now(), val: val}
	cache.mutex.Unlock()
}

func (cache *Cache) Get(key string) ([]byte, bool) {
	cache.mutex.RLock()
	cacheobj, exists := cache.cached[key]
	if exists{
		val := cacheobj.val
		cache.mutex.RUnlock()
		return val, exists
	}
	cache.mutex.RUnlock()
	return nil,false
}

func (cache *Cache) reapLoop(Interval time.Duration) {
	for {
		time.Sleep(Interval)
		for key, cacheobj := range cache.cached{
			if time.Now().Sub(cacheobj.createdAt)>Interval {
				delete (cache.cached, key)
			}
		}
	}
}