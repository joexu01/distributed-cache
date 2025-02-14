package cache

type Cache interface {
	Set(string, []byte) error
	Get(string) ([]byte, error)
	Del(string) error
	GetStat() Stat
	NewScanner() Scanner
}

// Scanner
type Scanner interface {
	Scan() bool
	Key() string
	Value() []byte
	Close()
}
