package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -ldl -O3
import "C"
import (
	"errors"
	"regexp"
	"strconv"
	"unsafe"
)

func (c *rocksDbCache) Set(key string, value []byte) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	v := C.CBytes(value)
	defer C.free(v)
	C.rocksdb_put(
		c.db, c.wo, k, C.size_t(len(key)), (*C.char)(v), C.size_t(len(value)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil
}

func (c *rocksDbCache) Get(key string) ([]byte, error) {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	var length C.size_t
	value := C.rocksdb_get(
		c.db, c.ro, k, C.size_t(len(key)), &length, &c.e)
	if c.e != nil {
		return nil, errors.New(C.GoString(c.e))
	}
	defer C.free(unsafe.Pointer(value))
	return C.GoBytes(unsafe.Pointer(value), C.int(length)), nil
}

func (c *rocksDbCache) Del(key string) error {
	k := C.CString(key)
	defer C.free(unsafe.Pointer(k))
	C.rocksdb_delete(c.db, c.wo, k, C.size_t(len(key)), &c.e)
	if c.e != nil {
		return errors.New(C.GoString(c.e))
	}
	return nil
}

func (c *rocksDbCache) GetStat() Stat {
	k := C.CString("rocksdb.aggregated-table-properties")
	defer C.free(unsafe.Pointer(k))
	v := C.rocksdb_property_value(c.db, k)
	defer C.free(unsafe.Pointer(v))
	p := C.GoString(v)

	r := regexp.MustCompile(`([^;]+)=([^;]+);`)
	s := Stat{}
	for _, subMatches := range r.FindAllStringSubmatch(p, -1) {
		if subMatches[1] == " # entries" {
			s.Count, _ = strconv.ParseInt(subMatches[2], 10, 64)
		} else if subMatches[1] == " raw key size" {
			s.KeySize, _ = strconv.ParseInt(subMatches[2], 10, 64)
		} else if subMatches[1] == " raw value size" {
			s.ValueSize, _ = strconv.ParseInt(subMatches[2], 10, 64)
		}
	}
	return s
}
