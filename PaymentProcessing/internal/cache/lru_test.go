package cache

import (
	"testing"
)

func TestLRUCachePutAndGet(t *testing.T) {
	cache := NewLRUCache[string](2)
	value1 := "Alpha"
	value2 := "Beta"

	cache.Put(1, value1)
	cache.Put(2, value2)

	if cache.Len() != 2 {
		t.Fatalf("expected len 2, got %d", cache.Len())
	}

	if got, ok := cache.Get(1); !ok || got != value1 {
		t.Fatalf("expected value 1 from cache, got %#v, ok=%v", got, ok)
	}

	if got, ok := cache.Get(2); !ok || got != value2 {
		t.Fatalf("expected value 2 from cache, got %#v, ok=%v", got, ok)
	}
}

func TestLRUCacheEvictsLeastRecentlyUsedItem(t *testing.T) {
	cache := NewLRUCache[string](2)
	value1 := "Alpha"
	value2 := "Beta"
	value3 := "Gamma"

	cache.Put(1, value1)
	cache.Put(2, value2)
	cache.Get(1)
	cache.Put(3, value3)

	if _, ok := cache.Get(2); ok {
		t.Fatal("expected least recently used item 2 to be evicted")
	}

	if got, ok := cache.Get(1); !ok || got != value1 {
		t.Fatalf("expected merchant 1 to remain in cache, got %#v, ok=%v", got, ok)
	}

	if got, ok := cache.Get(3); !ok || got != value3 {
		t.Fatalf("expected merchant 3 to be cached, got %#v, ok=%v", got, ok)
	}
}

func TestLRUCacheDelete(t *testing.T) {
	cache := NewLRUCache[string](2)
	value := "Delta"

	cache.Put(10, value)
	cache.Delete(10)

	if cache.Len() != 0 {
		t.Fatalf("expected len 0 after delete, got %d", cache.Len())
	}

	if _, ok := cache.Get(10); ok {
		t.Fatal("expected deleted item to be absent from cache")
	}
}
