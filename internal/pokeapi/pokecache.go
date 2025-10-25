package pokeapi

import (
	"sync"
	"time"
)

type Cache struct {
	Dat map[string]cacheEntry
	Mu  *sync.Mutex
}

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		Dat: make(map[string]cacheEntry),
		Mu:  &sync.Mutex{},
	}

	go c.reapLoop(interval)

	return c
}

func (C *Cache) Add(key string, val []byte) {
	C.Mu.Lock()
	defer C.Mu.Unlock()
	C.Dat[key] = cacheEntry{
		createdAt: time.Now().UTC(),
		val:       val,
	}
}

func (C *Cache) Get(key string) ([]byte, bool) {
	C.Mu.Lock()
	defer C.Mu.Unlock()
	val, ok := C.Dat[key]
	return val.val, ok
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	for range ticker.C {
		c.reap(time.Now().UTC(), interval)
	}
}

func (c *Cache) reap(now time.Time, last time.Duration) {
	c.Mu.Lock()
	defer c.Mu.Unlock()
	for k, v := range c.Dat {
		if v.createdAt.Before(now.Add(-last)) {
			delete(c.Dat, k)
		}
	}
}
