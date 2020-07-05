package tcp

import (
	"github.com/joexu01/distributed-cache/cache"
	"log"
	"net"
)

type Server struct {
	cache.Cache
}

func (s *Server) Listen() {
	listen, err := net.Listen("tcp", ":12346")
	if err != nil {
		panic(err)
	}

	for {
		conn, err := listen.Accept()
		if err != nil {
			log.Fatal("cannot establish connection")
		}
		go s.process(conn)
	}
}

func New(c cache.Cache) *Server {
	return &Server{c}
}
