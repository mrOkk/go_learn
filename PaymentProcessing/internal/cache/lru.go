package cache

import (
	"container/list"
	"sync"
)

type entry[T any] struct {
	key   int64
	value T
}

type LRUCache[T any] struct {
	mu       sync.RWMutex
	capacity int
	items    map[int64]*list.Element
	order    *list.List
}

func NewLRUCache[T any](capacity int) *LRUCache[T] {
	return &LRUCache[T]{
		capacity: capacity,
		items:    make(map[int64]*list.Element),
		order: list.New(),
	}
}

func (c *LRUCache[T]) Get(key int64) (T, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		var empty T
		return empty, false
	}

	c.order.MoveToFront(elem)
	return elem.Value.(entry[T]).value, true
}

func (c *LRUCache[T]) Put(key int64, value T) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.items[key]; ok {
		elem.Value = entry[T]{key: key, value: value}
		c.order.MoveToFront(elem)
		return
	}

	elem := c.order.PushFront(entry[T]{key: key, value: value})
	c.items[key] = elem

	if c.order.Len() > c.capacity {
		oldest := c.order.Back()
		c.order.Remove(oldest)
		delete(c.items, oldest.Value.(entry[T]).key)
	}
}

func (c *LRUCache[T]) Delete(key int64)  {
	c.mu.Lock()
	defer c.mu.Unlock()

	elem, ok := c.items[key]
	if !ok {
		return
	}

	c.order.Remove(elem)
	delete(c.items, key)
}

func (c *LRUCache[T]) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.order.Len()
}