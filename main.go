package main

import (
	"github.com/joexu01/distributed-cache/cache"
	"github.com/joexu01/distributed-cache/http"
)

func main() {
	c := cache.New("in_memory")
	http.New(c).Listen()
}
