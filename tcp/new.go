package tcp

import (
	"github.com/joexu01/distributed-cache/cache"
	"github.com/joexu01/distributed-cache/cluster"
	"log"
	"net"
)

type Server struct {
	cache.Cache
	cluster.Node
}

func (s *Server) Listen() {
	listen, err := net.Listen("tcp", s.Addr()+":12346")
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

func New(c cache.Cache, n cluster.Node) *Server {
	return &Server{c, n}
}
