package tcp

import (
	"bufio"
	"io"
	"log"
	"net"
)

type result struct {
	value []byte
	err   error
}

func (s *Server) get(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	key, err := s.readKey(r)
	if err != nil {
		c <- &result{nil, err}
		return
	}

	go func() {
		value, err := s.Get(key)
		c <- &result{value, err}
	}()
}

func (s *Server) set(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c

	key, value, err := s.readKeyAndValue(r)
	if err != nil {
		c <- &result{nil, err}
		return
	}

	go func() {
		c <- &result{nil, s.Set(key, value)}
	}()
}

func (s *Server) del(ch chan chan *result, r *bufio.Reader) {
	c := make(chan *result)
	ch <- c
	key, err := s.readKey(r)
	if err != nil {
		c <- &result{nil, err}
		return
	}

	go func() {
		c <- &result{nil, s.Del(key)}
	}()
}

func reply(conn net.Conn, resultChan chan chan *result) {
	defer conn.Close()
	for {
		c, open := <-resultChan
		if !open {
			return
		}
		r := <-c
		err := sendResponse(r.value, r.err, conn)
		if err != nil {
			log.Println("connection closed due to err:", err)
			return
		}
	}
}

func (s *Server) process(conn net.Conn) {
	r := bufio.NewReader(conn)
	resultChan := make(chan chan *result, 5000)
	defer close(resultChan)
	go reply(conn, resultChan)
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
			s.set(resultChan, r)
		} else if op == 'G' {
			s.get(resultChan, r)
		} else if op == 'D' {
			s.del(resultChan, r)
		} else {
			log.Println("connection closed due to invalid operation:", op)
			return
		}
	}
}
