package tcp

import (
	"bufio"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func readLen(r *bufio.Reader) (int, error) {
	tmp, err := r.ReadString(' ')
	if err != nil {
		return 0, err
	}
	length, err := strconv.Atoi(strings.TrimSpace(tmp))
	if err != nil {
		return 0, err
	}
	return length, nil
}

func sendResponse(value []byte, err error, conn net.Conn) error {
	if err != nil {
		errString := err.Error()
		tmp := fmt.Sprintf("-%d ", len(errString)) + errString
		_, e := conn.Write([]byte(tmp))
		return e
	}
	valueLen := fmt.Sprintf("%d ", len(value))
	_, e := conn.Write(append([]byte(valueLen), value...))
	return e
}
