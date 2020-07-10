package http

import (
	"github.com/joexu01/distributed-cache/cache"
	"github.com/joexu01/distributed-cache/cluster"
	"log"
	"net/http"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	http.Handle("/cluster", s.clusterHandler())
	log.Fatal(http.ListenAndServe(s.Addr()+":12345", nil))
}

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
