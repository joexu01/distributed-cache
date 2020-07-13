package cache

import (
	"log"
)

func New(typ string, ttl int) Cache {
	var c Cache
	if typ == "in_memory" {
		c = newInMemoryCache(ttl)
	}
	if typ == "rocksdb" {
		c = newRocksDbCache(ttl)
	}
	if c == nil {
		panic("unknown cache type " + typ)
	}
	log.Println(typ, "ready to serve")
	return c
}
