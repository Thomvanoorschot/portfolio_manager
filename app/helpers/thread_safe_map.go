package helpers

import "sync"

type ThreadSafeMap[K comparable, V interface{}] struct {
	sync.RWMutex
	Entries map[K]*V
}

func (sn *ThreadSafeMap[K, V]) Add(symbol K, holding *V) {
	sn.Lock()
	defer sn.Unlock()
	sn.Entries[symbol] = holding
}

func (sn *ThreadSafeMap[K, V]) Get(symbol K) *V {
	sn.RLock()
	defer sn.RUnlock()
	return sn.Entries[symbol]
}
