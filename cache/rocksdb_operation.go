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
	"time"
	"unsafe"
)

const BatchSize = 100

func (c *rocksDbCache) Set(key string, value []byte) error {
	c.ch <- &pair{key, value}
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

func flushBach(
	db *C.rocksdb_t, b *C.rocksdb_writebatch_t, o *C.rocksdb_writeoptions_t) {
	var e *C.char
	C.rocksdb_write(db, o, b, &e)
	if e != nil {
		panic(C.GoString(e))
	}
	C.rocksdb_writebatch_clear(b)
}

func writeFunc(
	db *C.rocksdb_t, c chan *pair, o *C.rocksdb_writeoptions_t) {
	count := 0
	//计时器的触发时间是1s
	t := time.NewTimer(time.Second)
	b := C.rocksdb_writebatch_create()
	for {
		select {
		//如果channel中的数据先抵达
		case p := <-c:
			count++
			key := C.CString(p.key)
			value := C.CBytes(p.value)
			C.rocksdb_writebatch_put(
				b, key, C.size_t(len(p.key)), (*C.char)(value), C.size_t(len(p.value)))
			C.free(unsafe.Pointer(key))
			C.free(value)
			if count == BatchSize {
				flushBach(db, b, o)
				count = 0
			}
			if !t.Stop() {
				<-t.C
			}
			t.Reset(time.Second)
		//1s内没有写操作请求就先把内存中的数据写到磁盘中
		case <-t.C:
			if count != 0 {
				flushBach(db, b, o)
				count = 0
			}
			t.Reset(time.Second)
		}
	}
}
