package main

import (
	"flag"
	"github.com/joexu01/distributed-cache/cache"
	"github.com/joexu01/distributed-cache/http"
	"github.com/joexu01/distributed-cache/tcp"
	"log"
)

var typ = flag.String("type", "in_memory", "cache type")

func main() {
	flag.Parse()
	log.Println("type is", *typ)
	c := cache.New(*typ)
	go tcp.New(c).Listen()
	http.New(c).Listen()
}
