package tests

import (
	"UrlShortener/internal/service"
	"context"
	"testing"
)

type MockStorage struct {
	m map[string]string
}

func (m *MockStorage) Save(ctx context.Context, code, url string) error {
	return nil
}

func (m *MockStorage) Get(ctx context.Context, code string) (string, error) {
	return "", nil
}

type MockCache struct {
}

func (m *MockCache) Save(ctx context.Context, code, url string) error {
	return nil
}

func (m *MockCache) Get(ctx context.Context, code string) (string, bool, error) {
	return "", true, nil
}

func GetService() *service.UrlService {
	ms := &MockStorage{}
	mc := &MockCache{}
	r := service.NewRepositoryService(mc, ms)
	return service.NewUrlService(r)
}

func TestServiceSave(t *testing.T) {
	s := GetService()
	code, err := s.Put(context.Background(), "https://fakeurl.fake")
	if err != nil {
		t.Fatal(err)
	}

	if len(code) != 6 {
		t.Fatal(code)
	}
}

func TestServiceGet(t *testing.T) {

}