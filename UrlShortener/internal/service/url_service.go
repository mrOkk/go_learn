package service

import (
	"UrlShortener/internal/repository"
	"context"
	"errors"
)

const retriesCount = 5

type UrlService struct {
	store UrlStorage
}

func NewUrlService(store UrlStorage) *UrlService {
	return &UrlService{
		store: store,
	}
}

func (r *UrlService) Get(ctx context.Context, code string) (string, error) {
	return r.store.Get(ctx, code)
}

func (r *UrlService) Put(ctx context.Context, url string) (string, error) {
	var err error
	for range retriesCount {
		code := generateShortCode()
		err = r.store.Save(ctx, code, url)
		if err == nil {
			return code, nil
		}
		if !errors.Is(err, &repository.ConflictError{}) {
			return "", err
		}
	}
	return "", err
}