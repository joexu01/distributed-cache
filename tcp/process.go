package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

func (s *Server) get(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}
	value, err := s.Get(key)
	return sendResponse(value, err, conn)
}

func (s *Server) set(conn net.Conn, r *bufio.Reader) error {
	key, value, err := s.readKeyAndValue(r)
	if err != nil {
		return err
	}

	return sendResponse(nil, s.Set(key, value), conn)
}

func (s *Server) del(conn net.Conn, r *bufio.Reader) error {
	key, err := s.readKey(r)
	if err != nil {
		return err
	}

	return sendResponse(nil, s.Del(key), conn)
}

func (s *Server) process(conn net.Conn) {
	defer conn.Close()

	r := bufio.NewReader(conn)
	for {
		//在无限循环中调用ReadByte()方法，在不发生错误的情况下
		//客户端可以复用这个连接
		op, err := r.ReadByte()
		if err != nil {
			if err != io.EOF {
				log.Println("connection closed due to error:", err)
			}
			return
		}
		if op == 'S' {
			err = s.set(conn, r)
		} else if op == 'G' {
			err = s.get(conn, r)
		} else if op == 'D' {
			err = s.del(conn, r)
		} else {
			log.Println("connection closed due to invalid operation:", op)
			return
		}
		if err != nil {
			log.Println("connection closed due to error:", err)
			return
		}
	}
}
