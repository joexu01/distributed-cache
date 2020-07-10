package tcp

import (
	"bufio"
	"errors"
	"io"
)

func (s *Server) readKey(r *bufio.Reader) (string, error) {
	keyLen, err := readLen(r)
	if err != nil {
		return "", err
	}

	key := make([]byte, keyLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return "", err
	}
	keyStr := string(key)
	addr, ok := s.ShouldProcess(keyStr)
	if !ok {
		return "", errors.New("redirect: " + addr)
	}
	return keyStr, nil
}

func (s *Server) readKeyAndValue(r *bufio.Reader) (string, []byte, error) {
	keyLen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}
	valueLen, err := readLen(r)
	if err != nil {
		return "", nil, err
	}

	key := make([]byte, keyLen)
	_, err = io.ReadFull(r, key)
	if err != nil {
		return "", nil, err
	}
	keyStr := string(key)
	addr, ok := s.ShouldProcess(keyStr)
	if !ok {
		return "", nil, errors.New("redirect: " + addr)
	}
	value := make([]byte, valueLen)
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", nil, err
	}
	return keyStr, value, nil
}

