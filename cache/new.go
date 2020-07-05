package cache

import (
	"log"
)

func New(cacheType string) Cache {
	var c Cache
	if cacheType == "in_memory" {
		c = newInMemoryCache()
	}
	if c == nil {
		panic("unknown cache type: " + cacheType)
	}
	log.Println(cacheType, "ready to serve")
	return c
}
