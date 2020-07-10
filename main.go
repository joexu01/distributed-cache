package main

import (
	"flag"
	"github.com/joexu01/distributed-cache/cache"
	"github.com/joexu01/distributed-cache/cluster"
	"github.com/joexu01/distributed-cache/http"
	"github.com/joexu01/distributed-cache/tcp"
	"log"
)

var (
	typ = flag.String("type", "in_memory", "cache type [in_memory|rocksdb]")
	nodeAddr = flag.String("node", "127.0.0.1", "node address, default to 127.0.0.1")
	clusterAddr = flag.String("cluster", "", "cluster address")
)

func main() {
	flag.Parse()
	log.Printf("Service type: %s.\n", *typ)
	log.Printf("Node address: %s.\n", *nodeAddr)
	log.Printf("Cluster address: %s.\n", *clusterAddr)
	c := cache.New(*typ)
	n, err := cluster.New(*nodeAddr, *clusterAddr)
	if err != nil {
		panic(err)
	}
	go tcp.New(c, n).Listen()
	http.New(c, n).Listen()
}
