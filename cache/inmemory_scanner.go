package cache

type inMemoryScanner struct {
	pair
	pairChan  chan *pair
	closeChan chan struct{}
}

func (s *inMemoryScanner) Scan() bool {
	p, ok := <-s.pairChan
	if ok {
		s.key, s.value = p.key, p.value
	}
	return ok
}

func (s *inMemoryScanner) Key() string {
	return s.key
}

func (s *inMemoryScanner) Value() []byte {
	return s.value
}

func (s *inMemoryScanner) Close() {
	close(s.closeChan)
}

func (c *inMemoryCache) NewScanner() Scanner {
	pairCh := make(chan *pair)
	closeCh := make(chan struct{})
	go func() {
		defer close(pairCh)
		c.mutex.RLock()
		for k, v := range c.cache {
			c.mutex.RUnlock()
			select {
			case <-closeCh:
				return
			case pairCh <- &pair{k, v.v}:
			}
			c.mutex.RLock()
		}
		c.mutex.RUnlock()
	}()
	return &inMemoryScanner{pair{}, pairCh, closeCh}
}
