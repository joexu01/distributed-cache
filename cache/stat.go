package cache

type Stat struct {
	Count     int64 `json:"count"`       //目前保存的K-V数量
	KeySize   int64 `json:"key_size" `   //Key的字节总数
	ValueSize int64 `json:"value_size" ` //Value的字节总数
}

func (s *Stat) add(k string, v []byte) {
	s.Count += 1
	s.KeySize += int64(len(k))
	s.ValueSize += int64(len(v))
}

func (s *Stat) del(k string, v []byte) {
	s.Count -= 1
	s.KeySize -= int64(len(k))
	s.ValueSize -= int64(len(v))
}
