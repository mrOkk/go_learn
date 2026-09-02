package cache

import (
	"App/internal/domain"
	"testing"
)

// TODO: make tests domain-agnostic
func TestLRUCachePutAndGet(t *testing.T) {
	cache := NewLRUCache[domain.Merchant](2)
	merchant1 := domain.Merchant{Id: 1, Name: "Alpha", IsActive: true, Balance: 100}
	merchant2 := domain.Merchant{Id: 2, Name: "Beta", IsActive: true, Balance: 200}

	cache.Put(1, merchant1)
	cache.Put(2, merchant2)

	if cache.Len() != 2 {
		t.Fatalf("expected len 2, got %d", cache.Len())
	}

	if got, ok := cache.Get(1); !ok || got != merchant1 {
		t.Fatalf("expected merchant 1 from cache, got %#v, ok=%v", got, ok)
	}

	if got, ok := cache.Get(2); !ok || got != merchant2 {
		t.Fatalf("expected merchant 2 from cache, got %#v, ok=%v", got, ok)
	}
}

func TestLRUCacheEvictsLeastRecentlyUsedItem(t *testing.T) {
	cache := NewLRUCache[domain.Merchant](2)
	merchant1 := domain.Merchant{Id: 1, Name: "Alpha", IsActive: true, Balance: 100}
	merchant2 := domain.Merchant{Id: 2, Name: "Beta", IsActive: true, Balance: 200}
	merchant3 := domain.Merchant{Id: 3, Name: "Gamma", IsActive: true, Balance: 300}

	cache.Put(1, merchant1)
	cache.Put(2, merchant2)
	cache.Get(1)
	cache.Put(3, merchant3)

	if _, ok := cache.Get(2); ok {
		t.Fatal("expected least recently used item 2 to be evicted")
	}

	if got, ok := cache.Get(1); !ok || got != merchant1 {
		t.Fatalf("expected merchant 1 to remain in cache, got %#v, ok=%v", got, ok)
	}

	if got, ok := cache.Get(3); !ok || got != merchant3 {
		t.Fatalf("expected merchant 3 to be cached, got %#v, ok=%v", got, ok)
	}
}

func TestLRUCacheDelete(t *testing.T) {
	cache := NewLRUCache[domain.Merchant](2)
	merchant := domain.Merchant{Id: 10, Name: "Delta", IsActive: false, Balance: 10}

	cache.Put(10, merchant)
	cache.Delete(10)

	if cache.Len() != 0 {
		t.Fatalf("expected len 0 after delete, got %d", cache.Len())
	}

	if _, ok := cache.Get(10); ok {
		t.Fatal("expected deleted item to be absent from cache")
	}
}
