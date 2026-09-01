package postgres

import (
	"App/internal/domain"
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository struct {
	pool *pgxpool.Pool
}

func NewTransactionRepository(pool *pgxpool.Pool) *TransactionRepository {
	return &TransactionRepository{
		pool: pool,
	}
}

func (r *TransactionRepository) Save(ctx context.Context, tx domain.Transaction) error {
	_, err := r.pool.Exec(ctx,
		`INSERT INTO transactions 
    			(id, merchant_id, amount, status, "timestamp")
			VALUES ($1, $2, $3, $4, $5)
			ON CONFLICT (id) DO NOTHING`,
			tx.ID, tx.MerchantID, tx.Amount, tx.Status, tx.Timestamp,
			)

	if err != nil {
		return fmt.Errorf("save transaction: %w", err)
	}

	return nil
}