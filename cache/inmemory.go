package cache

import "sync"

type inMemoryCache struct {
	cache map[string][]byte
	mutex sync.RWMutex
	Stat
}

func newInMemoryCache() *inMemoryCache {
	return &inMemoryCache{make(map[string][]byte), sync.RWMutex{}, Stat{}}
}

//Set(key string, value byte[]) is used to set a
//key-value pair in the memory.
func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	tmp, exist := c.cache[k]
	if exist {
		c.del(k, tmp)
	}
	c.cache[k] = v
	c.add(k, v)
	return nil
}

//Get(key string) Given the key, the function returns the
//value of the K-V pair.
func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache[k], nil
}

//Del(key string) is used to delete the K-V pair in the
//memory.
func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.cache[k]
	if exist {
		delete(c.cache, k)
		c.del(k, v)
	}
	return nil
}

//GetStat() returns a struct indicating the status of cache.
func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}
