package main

import (
	"github.com/joexu01/distributed-cache/cache"
	"github.com/joexu01/distributed-cache/http"
	"github.com/joexu01/distributed-cache/tcp"
)

func main() {
	c := cache.New("in_memory")
	go tcp.New(c).Listen()
	http.New(c).Listen()
}
