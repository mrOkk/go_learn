package service

import (
	"context"
)

type UrlStorage interface {
	Save(ctx context.Context, code, url string) error
	Get(ctx context.Context, code string) (string, error)
}

type CacheStorage interface {
	Save(ctx context.Context, code, url string) error
	Get(ctx context.Context, code string) (string, bool, error)
}

type RepositoryService struct {
	cache   CacheStorage
	storage UrlStorage
}

func NewRepositoryService(cache CacheStorage, storage UrlStorage) *RepositoryService {
	return &RepositoryService{
		cache:   cache,
		storage: storage,
	}
}

func (r *RepositoryService) Save(ctx context.Context, code, url string) error {
	err := r.storage.Save(ctx, code, url)
	if err != nil {
		return err
	}
	r.cache.Save(ctx, code, url)
	return nil
}

func (r *RepositoryService) Get(ctx context.Context, code string) (string, error) {
	code, ok, _ := r.cache.Get(ctx, code)
	if ok {
		return code, nil
	}
	code, err := r.storage.Get(ctx, code)
	if err != nil {
		return "", err
	}
	r.cache.Save(ctx, code, code)
	return code, nil
}
