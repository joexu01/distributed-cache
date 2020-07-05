package http

import (
	"github.com/joexu01/distributed-cache/cache"
	"log"
	"net/http"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	http.Handle("/cache/", s.cacheHandler())
	http.Handle("/status", s.statusHandler())
	log.Fatal(http.ListenAndServe(":12345", nil))
}

func New(c cache.Cache) *Server {
	return &Server{c}
}
