package repository

import (
	"context"
	"fmt"
)

type StubStorage struct {
	storage map[string]string
}

func NewStubStorage() *StubStorage {
	return &StubStorage{
		storage: make(map[string]string),
	}
}

type StubCache struct {
	storage map[string]string
}

func NewStubCache() *StubCache {
	return &StubCache{
		storage: make(map[string]string),
	}
}

func init() {

}

func (s *StubStorage) Save(_ context.Context, code string, url string) error {
	fmt.Printf("[StubStorage] Save(%s, %s)\n", code, url)
	s.storage[code] = url
	return nil
}

func (s *StubStorage) Get(_ context.Context, code string) (string, error) {
	url, ok := s.storage[code]
	if !ok {
		fmt.Printf("[StubStorage] Get(%s) | Url NOT found\n", code)
		return "", fmt.Errorf("url not found")
	}
	fmt.Printf("[StubStorage] Get(%s) | Url found: %s\n", code, url)
	return url, nil
}

func (s *StubCache) Save(_ context.Context, code string, url string) error {
	fmt.Printf("[StubCache] Save(%s, %s)\n", code, url)
	s.storage[code] = url
	return nil
}

func (s *StubCache) Get(_ context.Context, code string) (string, bool, error) {
	return "", false, nil
	//url, ok := s.storage[code]
	//if !ok {
	//	fmt.Printf("[StubCache] Get(%s) | Url NOT found\n", code)
	//	return "", false, fmt.Errorf("url not found")
	//}
	//fmt.Printf("[StubCache] Get(%s) | Url found: %s\n", code, url)
	//return url, true, nil
}