package cache

import (
	"App/internal/domain"
	"container/list"
	"sync"
)

type entry struct {
	key   int64
	value domain.Merchant
}

type LRUCache struct {
	mu       sync.RWMutex
	capacity int
	items    map[int64]*list.Element
	order    *list.List
}

func NewLRUCache(capacity int) *LRUCache {
	return &LRUCache{
		capacity: capacity,
		items:    make(map[int64]*list.Element),
		order: list.New(),
	}
}

func (c *LRUCache) Get(key int64) (domain.Merchant, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return domain.Merchant{}, false
	}

	c.order.MoveToFront(elem)
	return elem.Value.(entry).value, true
}

func (c *LRUCache) Put(key int64, value domain.Merchant) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value = entry{key: key, value: value}
		c.order.MoveToFront(elem)
		return
	}

	elem := c.order.PushFront(entry{key: key, value: value})
	c.items[key] = elem

	if c.order.Len() > c.capacity {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.items, oldest.Value.(entry).key)
	}
}

func (c *LRUCache) Delete(key int64)  {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return
	}

	c.order.Remove(elem)
	delete(c.items, key)
}

func (c *LRUCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.order.Len()
}