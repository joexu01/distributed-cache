package cache

import (
	"sync"
	"time"
)

type value struct {
	v       []byte
	created time.Time
}

type inMemoryCache struct {
	cache map[string]value
	mutex sync.RWMutex
	Stat
	ttl time.Duration
}

func newInMemoryCache(ttl int) *inMemoryCache {
	c := &inMemoryCache{
		make(map[string]value),
		sync.RWMutex{},
		Stat{},
		time.Duration(ttl) * time.Second}
	if ttl > 0 {
		go c.expire()
	}
	return c
}

//Set(key string, value byte[]) is used to set a
//key-value pair in the memory.
func (c *inMemoryCache) Set(k string, v []byte) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.cache[k] = value{v, time.Now()}
	c.add(k, v)
	return nil
}

//Get(key string) Given the key, the function returns the
//value of the K-V pair.
func (c *inMemoryCache) Get(k string) ([]byte, error) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.cache[k].v, nil
}

//Del(key string) is used to delete the K-V pair in the
//memory.
func (c *inMemoryCache) Del(k string) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	v, exist := c.cache[k]
	if exist {
		delete(c.cache, k)
		c.del(k, v.v)
	}
	return nil
}

//GetStat() returns a struct indicating the status of cache.
func (c *inMemoryCache) GetStat() Stat {
	return c.Stat
}

func (c *inMemoryCache) expire() {
	for {
		time.Sleep(c.ttl)
		c.mutex.RLock()
		for k, v := range c.cache {
			c.mutex.RUnlock()

			if v.created.Add(c.ttl).Before(time.Now()) {
				_ = c.Del(k)
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}
}
