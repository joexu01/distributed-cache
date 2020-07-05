package tcp

import (
	"bufio"
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

	return string(key), nil
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
	value := make([]byte, valueLen)
	_, err = io.ReadFull(r, value)
	if err != nil {
		return "", nil, err
	}

	return string(key), value, nil
}

