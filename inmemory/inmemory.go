package cache

import "sync"

// inMemoryCache is an implement of Cache
type inMemoryCache struct {
	c     map[string][]byte
	mutex sync.RWMutex
}

func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.c[k] = v
	return nil
}

func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.c[k], nil
}

func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	_, exist := c.c[k]
	if exist {
		delete(c.c, k)
	}
	return nil
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{
		c:     make(map[string][]byte),
		mutex: sync.RWMutex{},
	}
}
