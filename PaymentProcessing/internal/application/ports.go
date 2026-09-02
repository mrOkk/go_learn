package application

import (
	"App/internal/domain"
	"context"
)

type MerchantReader interface {
	GetById(ctx context.Context, id int64) (domain.Merchant, error)
}

type MerchantWriter interface {
	UpdateBalance(ctx context.Context, id int64, balance float64) error
}

type TransactionWriter interface {
	Save(ctx context.Context, transaction domain.Transaction) error
}

type TxManager interface {
	WithinTransaction(ctx context.Context, fn func(ctx context.Context) error) error
}

type CacheController[T any] interface {
	Get(key int64) (T, bool)
	Put(key int64, value T)
}