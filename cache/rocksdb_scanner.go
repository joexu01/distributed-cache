package cache

// #include <stdlib.h>
// #include "rocksdb/c.h"
// #cgo CFLAGS: -I${SRCDIR}/../rocksdb/include
// #cgo LDFLAGS: -L${SRCDIR}/../rocksdb -lrocksdb -lz -lpthread -lsnappy -lstdc++ -lm -ldl -O3
import "C"
import "unsafe"

type rocksdbScanner struct {
	iterator    *C.rocksdb_iterator_t
	initialized bool
}

func (c *rocksDbCache) NewScanner() Scanner {
	return &rocksdbScanner{C.rocksdb_create_iterator(c.db, c.ro), false}
}

func (s *rocksdbScanner) Scan() bool {
	if !s.initialized {
		C.rocksdb_iter_seek_to_first(s.iterator)
		s.initialized = true
	} else {
		C.rocksdb_iter_next(s.iterator)
	}
	return C.rocksdb_iter_valid(s.iterator) != 0
}

func (s *rocksdbScanner) Key() string {
	var length C.size_t
	k := C.rocksdb_iter_key(s.iterator, &length)
	return C.GoString(k)
}

func (s *rocksdbScanner) Value() []byte {
	var length C.size_t
	v := C.rocksdb_iter_value(s.iterator, &length)
	return C.GoBytes(unsafe.Pointer(v), C.int(length))
}

func (s *rocksdbScanner) Close() {
	C.rocksdb_iter_destroy(s.iterator)
}
